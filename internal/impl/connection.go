package impl

import (
	"github.com/notomo/wsxhub/internal/domain"
)

// ConnectionImpl :
type ConnectionImpl struct {
	websocketClient domain.WebsocketClient
	worker          domain.Worker
	targetWorker    domain.Worker
	id              string
	filterClause    domain.FilterClause
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
	return conn.websocketClient.Receive(0, func(message string) error {
		return conn.targetWorker.Receive(message)
	})
}

// Send :
func (conn *ConnectionImpl) Send(message string) error {
	return conn.websocketClient.Send(message)
}

// IsTarget :
func (conn *ConnectionImpl) IsTarget(targetMap map[string]interface{}) bool {
	return conn.filterClause.Match(targetMap)
}
