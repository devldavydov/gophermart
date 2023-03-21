package storage

type Storage interface {
	CreateUser(login, password string) (int, error)
	FindUser(login string) (int, string, error)
	AddOrder(userId int, orderNum string) error
	ListOrders(userId int) ([]OrderItem, error)
	GetBalance(userId int) (*Balance, error)
	ListWithdrawals(userId int) ([]WithdrawalItem, error)
	Close()
}
