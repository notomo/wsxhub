package mock

import (
	"github.com/notomo/wsxhub/internal/domain"
)

// FakeConnection :
type FakeConnection struct {
	domain.Connection
	FakeID   func() string
	FakeSend func(domain.Message) (bool, error)
}

// ID :
func (conn *FakeConnection) ID() string {
	return conn.FakeID()
}

// Send :
func (conn *FakeConnection) Send(message domain.Message) (bool, error) {
	return conn.FakeSend(message)
}
