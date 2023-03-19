package storage

type Storage interface {
	CreateUser(string, string) (int, error)
	FindUser(string) (int, string, error)
	Close()
}
