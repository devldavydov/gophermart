package storage

import "context"

type Storage interface {
	CreateUser(ctx context.Context, login, password string) (int, error)
	FindUser(ctx context.Context, login string) (int, string, error)

	AddOrder(ctx context.Context, userID int, orderNum string) error
	ListOrders(ctx context.Context, userID int) ([]OrderItem, error)
	GetOrdersToProcess(ctx context.Context) ([]OrderItem, error)
	ProcessOrder(ctx context.Context, orderNum string) error
	FinishOrder(ctx context.Context, orderNum string, userID int, success bool, accrual float64) error

	GetBalance(ctx context.Context, userID int) (*Balance, error)
	ListWithdrawals(ctx context.Context, userID int) ([]WithdrawalItem, error)
	BalanceWithdraw(ctx context.Context, userID int, orderNum string, sum float64) error

	Close()
}
