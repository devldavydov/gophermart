package storage

type Storage interface {
	CreateUser(login, password string) (int, error)
	FindUser(login string) (int, string, error)
	AddOrder(userId int, orderNum string) error
	ListOrders(userId int) ([]OrderItem, error)
	Close()
}
