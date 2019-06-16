package domain

// Connection :
type Connection interface {
	ID() string
	Listen() error
	Send(string) error
	Close() error
	IsTarget(map[string]interface{}) bool
}
