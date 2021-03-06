package impl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/notomo/wsxhub/internal"
	"github.com/notomo/wsxhub/internal/domain"
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
	ws, resp, wsErr := websocket.DefaultDialer.Dial(u, nil)
	if wsErr != nil {
		if resp == nil {
			return nil, wsErr
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, wsErr
		}

		msg := fmt.Sprintf("%s: %s", wsErr, body)
		return nil, errors.New(msg)
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
func (client *WebsocketClientImpl) Send(message []byte) error {
	return client.ws.WriteMessage(websocket.TextMessage, message)
}

// ReceiveOnce :
func (client *WebsocketClientImpl) ReceiveOnce(timeout int) ([]byte, error) {
	if timeout > 0 {
		client.ws.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}

	_, message, err := client.ws.ReadMessage()
	if err != nil {
		if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
			return nil, internal.ErrTimeout
		} else if _, ok := err.(*websocket.CloseError); ok {
			return nil, internal.ErrEOF
		}
		return nil, err
	}

	return message, nil
}

// Receive :
func (client *WebsocketClientImpl) Receive(timeout int, callback func([]byte) error) error {
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
