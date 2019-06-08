package impl

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
	"golang.org/x/net/websocket"
)

// ConnectionImpl :
type ConnectionImpl struct {
	ws           *websocket.Conn
	worker       domain.Worker
	targetWorker domain.Worker
	id           string
	done         chan bool
}

// ID :
func (conn *ConnectionImpl) ID() string {
	return conn.id
}

// Close :
func (conn *ConnectionImpl) Close() error {
	if err := conn.worker.Delete(conn); err != nil {
		return err
	}
	return conn.ws.Close()
}

// Listen :
func (conn *ConnectionImpl) Listen() error {
	if err := conn.worker.Add(conn); err != nil {
		return err
	}
	for {
		select {
		case <-conn.done:
			return nil
		default:
			var message string
			err := websocket.Message.Receive(conn.ws, &message)
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			if err := conn.targetWorker.Receive(message); err != nil {
				return err
			}
		}
	}
}

// Send :
func (conn *ConnectionImpl) Send(message string) error {
	return websocket.Message.Send(conn.ws, message)
}
