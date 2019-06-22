package domain

// Connection :
type Connection interface {
	ID() string
	Listen() error
	Send(Message) (bool, error)
	Close() error
}
