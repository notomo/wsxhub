package domain

// WebsocketClientFactory :
type WebsocketClientFactory interface {
	Client() (WebsocketClient, error)
}

// WebsocketClient :
type WebsocketClient interface {
	Send() error
	ReceiveOnce() (string, error)
	Receive(func(string) error) error
	Close() error
}
