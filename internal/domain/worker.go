package domain

// Worker : handles events from connections
type Worker interface {
	Run() error
	Add(Connection) error
	Delete(Connection) error
	Receive(Message) error
	NotifySendResult(error)
	Finish()
}
