package impl

import (
	"fmt"
	"testing"
	"time"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestID(t *testing.T) {
	id := "1"
	connection := &ConnectionImpl{
		id: id,
	}

	if got := connection.ID(); got != id {
		t.Errorf("want %v, but %v:", id, got)
	}
}

func TestClose(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		client := &mock.FakeWebsocketClient{
			FakeClose: func() error {
				return nil
			},
		}
		worker := &mock.FakeWorker{
			FakeDelete: func(connection domain.Connection) error {
				return nil
			},
		}
		connection := &ConnectionImpl{
			websocketClient: client,
			worker:          worker,
		}

		err := connection.Close()

		if err != nil {
			t.Errorf("should not be error, but actual: %v", err)
		}
	})

	t.Run("fail to delete connection", func(t *testing.T) {
		worker := &mock.FakeWorker{
			FakeDelete: func(connection domain.Connection) error {
				return fmt.Errorf("err")
			},
		}
		connection := &ConnectionImpl{
			worker: worker,
		}

		err := connection.Close()

		if err == nil {
			t.Errorf("should be error, but actual nil")
		}
	})
}

func TestListen(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		bytes := []byte("message")
		client := &mock.FakeWebsocketClient{
			FakeReceive: func(timeout int, callback func([]byte) error) error {
				return callback(bytes)
			},
		}

		worker := &mock.FakeWorker{
			FakeAdd: func(connection domain.Connection) error {
				return nil
			},
		}
		message := &mock.FakeMessage{}
		targetWorker := &mock.FakeWorker{
			FakeReceive: func(m domain.Message) error {
				if message != m {
					t.Errorf("should be the same message, but actual: %v, %v", message, m)
				}
				return nil
			},
		}
		messageFactory := &mock.FakeMessageFactory{
			FakeFromBytes: func(b []byte) (domain.Message, error) {
				if string(bytes) != string(b) {
					t.Errorf("should be the same bytes, but actual: %v, %v", bytes, b)
				}
				return message, nil
			},
		}

		connection := &ConnectionImpl{
			websocketClient: client,
			worker:          worker,
			targetWorker:    targetWorker,
			messageFactory:  messageFactory,
		}

		err := connection.Listen()

		if err != nil {
			t.Errorf("should not be error, but actual: %v", err)
		}
	})

	t.Run("fail to create message", func(t *testing.T) {
		client := &mock.FakeWebsocketClient{
			FakeReceive: func(timeout int, callback func([]byte) error) error {
				return callback([]byte(""))
			},
		}

		worker := &mock.FakeWorker{
			FakeAdd: func(connection domain.Connection) error {
				return nil
			},
		}

		messageFactory := &mock.FakeMessageFactory{
			FakeFromBytes: func(_ []byte) (domain.Message, error) {
				return nil, fmt.Errorf("err")
			},
		}

		connection := &ConnectionImpl{
			websocketClient: client,
			worker:          worker,
			messageFactory:  messageFactory,
		}

		err := connection.Listen()

		if err == nil {
			t.Errorf("should be error, but actual nil")
		}
	})

	t.Run("fail to add connection", func(t *testing.T) {
		worker := &mock.FakeWorker{
			FakeAdd: func(connection domain.Connection) error {
				return fmt.Errorf("err")
			},
		}
		connection := &ConnectionImpl{
			worker: worker,
		}

		err := connection.Listen()

		if err == nil {
			t.Errorf("should be error, but actual nil")
		}
	})
}

func TestSend(t *testing.T) {
	t.Run("filtered", func(t *testing.T) {
		message := &mock.FakeMessage{}

		filterClause := &mock.FakeFilterClause{
			FakeMatch: func(m domain.Message) (bool, error) {
				if message != m {
					t.Errorf("should be the same message, but actual: %v, %v", message, m)
				}
				return false, nil
			},
		}

		connection := &ConnectionImpl{
			filterClause: filterClause,
		}

		sent, err := connection.Send(message)
		if sent == true {
			t.Errorf("should not send")
		}
		if err != nil {
			t.Errorf("should not be error, but actual: %v", err)
		}
	})

	t.Run("filter error", func(t *testing.T) {
		message := &mock.FakeMessage{}

		filterClause := &mock.FakeFilterClause{
			FakeMatch: func(m domain.Message) (bool, error) {
				return false, fmt.Errorf("err")
			},
		}

		connection := &ConnectionImpl{
			filterClause: filterClause,
		}

		if _, err := connection.Send(message); err == nil {
			t.Errorf("should be error")
		}
	})

	t.Run("send", func(t *testing.T) {
		client := &mock.FakeWebsocketClient{
			FakeSend: func(b []byte) error {
				return nil
			},
		}

		bytes := []byte("message")
		message := &mock.FakeMessage{
			FakeBytes: func() []byte {
				return bytes
			},
		}

		filterClause := &mock.FakeFilterClause{
			FakeMatch: func(_ domain.Message) (bool, error) {
				return true, nil
			},
		}

		connection := &ConnectionImpl{
			websocketClient: client,
			filterClause:    filterClause,
		}

		sent, err := connection.Send(message)
		if sent == false {
			t.Errorf("should send")
		}
		if err != nil {
			t.Errorf("should not be error, but actual: %v", err)
		}
	})

	t.Run("timer stop and start", func(t *testing.T) {
		client := &mock.FakeWebsocketClient{
			FakeSend: func(b []byte) error {
				return nil
			},
		}

		bytes := []byte("message")
		message := &mock.FakeMessage{
			FakeBytes: func() []byte {
				return bytes
			},
		}

		filterClause := &mock.FakeFilterClause{
			FakeMatch: func(_ domain.Message) (bool, error) {
				return true, nil
			},
		}

		notified := make(chan bool)
		worker := &mock.FakeWorker{
			FakeNotifySendResult: func(err error) {
				notified <- true
			},
		}

		connection := &ConnectionImpl{
			websocketClient: client,
			filterClause:    filterClause,
			worker:          worker,
			debounce:        100,
			debounceTimer:   time.NewTimer(time.Duration(100) * time.Millisecond),
		}

		sent, err := connection.Send(message)

		select {
		case <-notified:
			if sent == true {
				t.Errorf("should not send")
			}
			if err != nil {
				t.Errorf("should not be error, but actual: %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("should be notified")
		}
	})
}
