package domain

// WebsocketClientFactory :
type WebsocketClientFactory interface {
	Client() (WebsocketClient, error)
}

// WebsocketClient :
type WebsocketClient interface {
	Send([]byte) error
	ReceiveOnce(int) ([]byte, error)
	Receive(int, func([]byte) error) error
	Close() error
}
