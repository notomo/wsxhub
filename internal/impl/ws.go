package impl

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/notomo/wsxhub/internal"
	"github.com/notomo/wsxhub/internal/domain"
	"golang.org/x/net/websocket"
)

// WebsocketClientFactoryImpl :
type WebsocketClientFactoryImpl struct {
	Port         string
	FilterSource string
	Debounce     int
}

// Client :
func (factory *WebsocketClientFactoryImpl) Client() (domain.WebsocketClient, error) {
	params := url.Values{"filter": {factory.FilterSource}, "debounce": {strconv.Itoa(factory.Debounce)}}
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
func (client *WebsocketClientImpl) Send(message string) error {
	return websocket.Message.Send(client.ws, message)
}

// ReceiveOnce :
func (client *WebsocketClientImpl) ReceiveOnce(timeout int) (string, error) {
	if timeout > 0 {
		client.ws.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}

	var message string
	if err := websocket.Message.Receive(client.ws, &message); err != nil {
		if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
			return "", internal.ErrTimeout
		} else if err == io.EOF {
			return "", internal.ErrEOF
		}
		return "", err
	}

	return message, nil
}

// Receive :
func (client *WebsocketClientImpl) Receive(timeout int, callback func(string) error) error {
	for {
		message, err := client.ReceiveOnce(timeout)
		if err == internal.ErrEOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := callback(message); err != nil {
			return err
		}
	}
}

// Close :
func (client *WebsocketClientImpl) Close() error {
	return client.ws.Close()
}
