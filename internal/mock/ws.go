package mock

import "github.com/notomo/wsxhub/internal/domain"

// FakeWebsocketClientFactory :
type FakeWebsocketClientFactory struct {
	domain.WebsocketClientFactory
	FakeClient func() (domain.WebsocketClient, error)
}

// Client :
func (factory *FakeWebsocketClientFactory) Client() (domain.WebsocketClient, error) {
	return factory.FakeClient()
}

// FakeWebsocketClient :
type FakeWebsocketClient struct {
	domain.WebsocketClient
	FakeSend        func([]byte) error
	FakeClose       func() error
	FakeReceive     func(int, func([]byte) error) error
	FakeReceiveOnce func(int) ([]byte, error)
}

// Send :
func (factory *FakeWebsocketClient) Send(bytes []byte) error {
	return factory.FakeSend(bytes)
}

// Close :
func (factory *FakeWebsocketClient) Close() error {
	return factory.FakeClose()
}

// Receive :
func (factory *FakeWebsocketClient) Receive(timeout int, callback func([]byte) error) error {
	return factory.FakeReceive(timeout, callback)
}

// ReceiveOnce :
func (factory *FakeWebsocketClient) ReceiveOnce(timeout int) ([]byte, error) {
	return factory.FakeReceiveOnce(timeout)
}
