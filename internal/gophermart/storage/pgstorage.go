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
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

const (
	_databaseRequestTimeout = 10 * time.Second
	_userUniqueConstraint   = `duplicate key value violates unique constraint "users_username_key"`
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

	var userId int
	err := pg.db.QueryRowContext(ctx, _sqlCreateUser, login, password).Scan(&userId)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Message == _userUniqueConstraint {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}

	return userId, nil
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

func (pg *PgStorage) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	_, err := pg.db.ExecContext(ctx, _sqlCreateTableUser)
	if err != nil {
		return err
	}

	return nil
}
