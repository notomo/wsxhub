package impl

import (
	"fmt"
	"net/url"

	"github.com/notomo/wsxhub/internal/domain"
	"golang.org/x/net/websocket"
)

// WebsocketClientFactoryImpl :
type WebsocketClientFactoryImpl struct {
	Port string
}

// Client :
func (factory *WebsocketClientFactoryImpl) Client() (domain.WebsocketClient, error) {
	params := url.Values{}
	u := fmt.Sprintf("ws://localhost:%s/?%s", factory.Port, params.Encode())
	ws, err := websocket.Dial(u, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	return &WebsocketClientImpl{
		ws: ws,
	}, nil
}

// WebsocketClientImpl :
type WebsocketClientImpl struct {
	ws *websocket.Conn
}

// Send :
func (client *WebsocketClientImpl) Send() error {
	return nil
}

// ReceiveOnce :
func (client *WebsocketClientImpl) ReceiveOnce() (string, error) {
	return "", nil
}

// Receive :
func (client *WebsocketClientImpl) Receive(callback func(string) error) error {
	return nil
}

// Close :
func (client *WebsocketClientImpl) Close() error {
	return client.ws.Close()
}
