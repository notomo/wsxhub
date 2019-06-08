package domain

// Worker :
type Worker interface {
	Run(Server, func(map[string]Connection, string) error) error
	Add(Connection) error
	Delete(Connection) error
	Receive(string) error
}
