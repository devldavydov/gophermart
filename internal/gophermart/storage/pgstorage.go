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
	_databaseRequestTimeout = 10 * time.Second
	_userUniqueConstraint   = `duplicate key value violates unique constraint "users_username_key"`
	_orderKeyConstraint     = `duplicate key value violates unique constraint "orders_pkey"`
	_balanceSumConstraint   = `new row for relation "balance" violates check constraint "balance_current_check"`
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

func (pg *PgStorage) CreateUser(login, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	tx, err := pg.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userId int
	err = tx.QueryRowContext(ctx, _sqlCreateUser, login, password).Scan(&userId)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Message == _userUniqueConstraint {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}

	_, err = tx.ExecContext(ctx, _sqlCreateUserBalance, userId)
	if err != nil {
		return 0, err
	}

	return userId, tx.Commit()
}

func (pg *PgStorage) FindUser(login string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var userId int
	var userPassword string
	err := pg.db.QueryRowContext(ctx, _sqlFindUser, login).Scan(&userId, &userPassword)
	switch {
	case err == sql.ErrNoRows:
		return 0, "", ErrUserNotFound
	case err != nil:
		return 0, "", err
	}

	return userId, userPassword, nil
}

func (pg *PgStorage) AddOrder(userId int, orderNum string) error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	foundUserOrder, err := pg.findUserOrder(userId, orderNum)
	if err != nil {
		return err
	}

	if foundUserOrder {
		return ErrOrderAlreadyExistsFromUser
	}

	_, err = pg.db.ExecContext(ctx, _sqlAddOrder, orderNum, userId)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Message == _orderKeyConstraint {
			return ErrOrderAlreadyExists
		}
		return err
	}

	return nil
}

func (pg *PgStorage) ListOrders(userId int) ([]OrderItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	rows, err := pg.db.QueryContext(ctx, _sqlListOrders, userId)
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
			item.Accrual = &r.accrual.Int32
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

func (pg *PgStorage) GetBalance(userId int) (*Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var balance Balance
	err := pg.db.QueryRowContext(ctx, _sqlGetUserBalance, userId).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (pg *PgStorage) ListWithdrawals(userId int) ([]WithdrawalItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	rows, err := pg.db.QueryContext(ctx, _sqlListWithdrawals, userId)
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

func (pg *PgStorage) BalanceWithdraw(userId int, orderNum string, sum float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update balance
	_, err = tx.ExecContext(ctx, _sqlUpdateUserBalance, userId, sum)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Message == _balanceSumConstraint {
			return ErrNotEnoughBalance
		}
		return err
	}

	// Insert withdrawal record
	_, err = tx.ExecContext(ctx, _sqlAddWithdrawal, userId, orderNum, sum)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pg *PgStorage) findUserOrder(userId int, orderNum string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var dummy int
	err := pg.db.QueryRowContext(ctx, _sqlFindUserOrder, orderNum, userId).Scan(&dummy)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}

	return true, err
}

func (pg *PgStorage) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	for _, createTbl := range []string{_sqlCreateTableUser, _sqlCreateTableOrders, _sqlCreateTableBalance, _sqlCreateTableWithdrawals} {
		_, err := pg.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}
