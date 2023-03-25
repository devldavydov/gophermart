package storage

type Storage interface {
	CreateUser(login, password string) (int, error)
	FindUser(login string) (int, string, error)

	AddOrder(userID int, orderNum string) error
	ListOrders(userID int) ([]OrderItem, error)
	GetOrdersToProcess() ([]OrderItem, error)
	ProcessOrder(orderNum string) error
	FinishOrder(orderNum string, userID int, success bool, accrual float64) error

	GetBalance(userID int) (*Balance, error)
	ListWithdrawals(userID int) ([]WithdrawalItem, error)
	BalanceWithdraw(userID int, orderNum string, sum float64) error

	Close()
}
