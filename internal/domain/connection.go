package domain

// Connection :
type Connection interface {
	Close() error
	Listen() error
	ID() string
	Send(string) error
}
