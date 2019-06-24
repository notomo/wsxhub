package impl

import (
	"time"

	"github.com/notomo/wsxhub/internal/domain"
)

// ConnectionImpl :
type ConnectionImpl struct {
	websocketClient domain.WebsocketClient
	worker          domain.Worker
	targetWorker    domain.Worker
	id              string
	filterClause    domain.FilterClause
	debounce        int
	debounceTimer   *time.Timer
	messageFactory  domain.MessageFactory
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
	return conn.websocketClient.Close()
}

// Listen :
func (conn *ConnectionImpl) Listen() error {
	if err := conn.worker.Add(conn); err != nil {
		return err
	}
	return conn.websocketClient.Receive(0, func(bytes []byte) error {
		message, err := conn.messageFactory.FromBytes(bytes)
		if err != nil {
			return err
		}

		return conn.targetWorker.Receive(message)
	})
}

// Send :
func (conn *ConnectionImpl) Send(message domain.Message) (bool, error) {
	if !conn.filterClause.Match(message.Unmarshaled()) {
		return false, nil
	}

	if conn.debounceTimer != nil {
		conn.debounceTimer.Stop()
	}
	if conn.debounce > 0 {
		conn.debounceTimer = time.AfterFunc(time.Duration(conn.debounce)*time.Millisecond, func() {
			err := conn.websocketClient.Send(message.Bytes())
			conn.worker.NotifySendResult(err)
		})
		return false, nil
	}
	return true, conn.websocketClient.Send(message.Bytes())
}
