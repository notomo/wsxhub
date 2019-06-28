package mock

import "github.com/notomo/wsxhub/internal/domain"

// FakeWorker :
type FakeWorker struct {
	domain.Worker
	FakeDelete           func(domain.Connection) error
	FakeAdd              func(domain.Connection) error
	FakeReceive          func(domain.Message) error
	FakeNotifySendResult func(error)
}

// Delete :
func (factory *FakeWorker) Delete(connection domain.Connection) error {
	return factory.FakeDelete(connection)
}

// Add :
func (factory *FakeWorker) Add(connection domain.Connection) error {
	return factory.FakeAdd(connection)
}

// Receive :
func (factory *FakeWorker) Receive(message domain.Message) error {
	return factory.FakeReceive(message)
}

// NotifySendResult :
func (factory *FakeWorker) NotifySendResult(err error) {
	factory.FakeNotifySendResult(err)
}
