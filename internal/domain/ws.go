package domain

// WebsocketClientFactory :
type WebsocketClientFactory interface {
	Client() (WebsocketClient, error)
}

// WebsocketClient :
type WebsocketClient interface {
	Send(string) error
	ReceiveOnce(int) (string, error)
	Receive(int, func(string) error) error
	Close() error
}
