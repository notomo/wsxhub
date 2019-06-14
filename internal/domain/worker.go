package domain

// Worker :
type Worker interface {
	Run() error
	Add(Connection) error
	Delete(Connection) error
	Receive(string) error
}
