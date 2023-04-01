package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrOrderAlreadyExists         = errors.New("order already exists")
	ErrOrderAlreadyExistsFromUser = errors.New("order already exists from user")
	ErrNoOrders                   = errors.New("no orders")
	ErrNoWithdrawals              = errors.New("no withdrawals")
	ErrNotEnoughBalance           = errors.New("not enough balance")
)

const (
	_databaseInitTimeout = 10 * time.Second

	_constraintCheckViolation  pq.ErrorCode = "23514"
	_constraintUniqueViolation pq.ErrorCode = "23505"
	_constraintBalanceCheck                 = "balance_current_check"
	_constraintUsernameCheck                = "users_username_key"
	_constraintOrderPkeyCheck               = "orders_pkey"
)

type PgStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPgStorage(pgConnString string, logger *logrus.Logger) (*PgStorage, error) {
	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		return nil, err
	}

	pgstorage := &PgStorage{db: db, logger: logger}

	if err = pgstorage.init(); err != nil {
		return nil, err
	}

	return pgstorage, nil
}

var _ Storage = (*PgStorage)(nil)

func (pg *PgStorage) Close() {
	if pg.db == nil {
		return
	}

	err := pg.db.Close()
	if err != nil {
		pg.logger.Errorf("Database conn close err: %v", err)
	}
}

func (pg *PgStorage) CreateUser(ctx context.Context, login, password string) (int, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRowContext(ctx, _sqlCreateUser, login, password).Scan(&userID)
	if err != nil {
		var pqErr *pq.Error
		if !errors.As(err, &pqErr) {
			return 0, err
		}

		if pqErr.Code == _constraintUniqueViolation && pqErr.Constraint == _constraintUsernameCheck {
			return 0, ErrUserAlreadyExists
		}

		return 0, err
	}

	_, err = tx.ExecContext(ctx, _sqlCreateUserBalance, userID)
	if err != nil {
		return 0, err
	}

	return userID, tx.Commit()
}

func (pg *PgStorage) FindUser(ctx context.Context, login string) (int, string, error) {
	var userID int
	var userPassword string
	err := pg.db.QueryRowContext(ctx, _sqlFindUser, login).Scan(&userID, &userPassword)
	switch {
	case err == sql.ErrNoRows:
		return 0, "", ErrUserNotFound
	case err != nil:
		return 0, "", err
	}

	return userID, userPassword, nil
}

func (pg *PgStorage) AddOrder(ctx context.Context, userID int, orderNum string) error {
	foundUserOrder, err := pg.findUserOrder(ctx, userID, orderNum)
	if err != nil {
		return err
	}

	if foundUserOrder {
		return ErrOrderAlreadyExistsFromUser
	}

	_, err = pg.db.ExecContext(ctx, _sqlAddOrder, orderNum, userID)
	if err != nil {
		var pqErr *pq.Error
		if !errors.As(err, &pqErr) {
			return err
		}

		if pqErr.Code == _constraintUniqueViolation && pqErr.Constraint == _constraintOrderPkeyCheck {
			return ErrOrderAlreadyExists
		}

		return err
	}

	return nil
}

func (pg *PgStorage) ListOrders(ctx context.Context, userID int) ([]OrderItem, error) {
	rows, err := pg.db.QueryContext(ctx, _sqlListOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem

	for rows.Next() {
		var r orderItemRow
		err = rows.Scan(&r.number, &r.status, &r.accrual, &r.uploadedAt)
		if err != nil {
			return nil, err
		}

		item := OrderItem{
			Number:     r.number,
			Status:     r.status,
			UploadedAt: r.uploadedAt,
		}
		if r.accrual.Valid {
			item.Accrual = &r.accrual.Float64
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrNoOrders
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (pg *PgStorage) GetOrdersToProcess(ctx context.Context) ([]OrderItem, error) {
	rows, err := pg.db.QueryContext(ctx, _sqlListOrdersForProcessing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem

	for rows.Next() {
		var r orderItemRow
		err = rows.Scan(&r.number, &r.userID)
		if err != nil {
			return nil, err
		}

		items = append(items, OrderItem{Number: r.number, UserID: r.userID})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (pg *PgStorage) ProcessOrder(ctx context.Context, orderNum string) error {
	_, err := pg.db.ExecContext(ctx, _sqlUpdateOrderForProcessing, orderNum)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) FinishOrder(ctx context.Context, orderNum string, userID int, success bool, accrual float64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if !success {
		_, err = tx.ExecContext(ctx, _sqlUpdateOrderForInvalid, orderNum)
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	_, err = tx.ExecContext(ctx, _sqlUpdateOrderForProcessed, orderNum, accrual)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, _sqlUpdateUserBalance, userID, accrual, 0)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pg *PgStorage) GetBalance(ctx context.Context, userID int) (*Balance, error) {
	var balance Balance
	err := pg.db.QueryRowContext(ctx, _sqlGetUserBalance, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (pg *PgStorage) ListWithdrawals(ctx context.Context, userID int) ([]WithdrawalItem, error) {
	rows, err := pg.db.QueryContext(ctx, _sqlListWithdrawals, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []WithdrawalItem

	for rows.Next() {
		var item WithdrawalItem
		err = rows.Scan(&item.Order, &item.Sum, &item.ProcessedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrNoWithdrawals
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (pg *PgStorage) BalanceWithdraw(ctx context.Context, userID int, orderNum string, sum float64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update balance
	_, err = tx.ExecContext(ctx, _sqlUpdateUserBalance, userID, -sum, sum)
	if err != nil {
		var pqErr *pq.Error
		if !errors.As(err, &pqErr) {
			return err
		}

		if pqErr.Code == _constraintCheckViolation && pqErr.Constraint == _constraintBalanceCheck {
			return ErrNotEnoughBalance
		}

		return err
	}

	// Insert withdrawal record
	_, err = tx.ExecContext(ctx, _sqlAddWithdrawal, userID, orderNum, sum)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pg *PgStorage) findUserOrder(ctx context.Context, userID int, orderNum string) (bool, error) {
	var dummy int
	err := pg.db.QueryRowContext(ctx, _sqlFindUserOrder, orderNum, userID).Scan(&dummy)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}

	return true, err
}

func (pg *PgStorage) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseInitTimeout)
	defer cancel()

	for _, createTbl := range []string{_sqlCreateTableUser, _sqlCreateTableOrders, _sqlCreateTableBalance, _sqlCreateTableWithdrawals} {
		_, err := pg.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}
