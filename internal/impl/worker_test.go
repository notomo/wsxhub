package impl

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestAddDelete(t *testing.T) {
	writer := &bytes.Buffer{}
	log.SetOutput(writer)

	id := "1"
	conn := &mock.FakeConnection{
		FakeID: func() string {
			return id
		},
	}

	worker := NewWorker("test")

	go func() {
		worker.Add(conn)
		worker.Finish()
	}()
	if err := worker.Run(); err != nil {
		t.Errorf("should not be error: %v", err)
	}

	if got, want := len(worker.Conns), 1; got != want {
		t.Errorf("want %v, but %v:", want, got)
	}

	go func() {
		worker.Delete(conn)
		worker.Finish()
	}()
	if err := worker.Run(); err != nil {
		t.Errorf("should not be error: %v", err)
	}

	if got, want := len(worker.Conns), 0; got != want {
		t.Errorf("want %v, but %v:", want, got)
	}
}

func TestReceive(t *testing.T) {
	writer := &bytes.Buffer{}
	log.SetOutput(writer)

	message := &mock.FakeMessage{}
	id := "1"

	t.Run("ok", func(t *testing.T) {
		conn := &mock.FakeConnection{
			FakeID: func() string {
				return id
			},
			FakeSend: func(msg domain.Message) (bool, error) {
				if message != msg {
					t.Errorf("should be the same message, but actual: %v, %v", message, msg)
				}
				return true, nil
			},
		}

		worker := NewWorker("test")

		go func() {
			worker.Add(conn)
			worker.Receive(message)
			worker.Finish()
		}()
		if err := worker.Run(); err != nil {
			t.Errorf("should not be error: %v", err)
		}
	})

	t.Run("fail to send", func(t *testing.T) {
		conn := &mock.FakeConnection{
			FakeID: func() string {
				return id
			},
			FakeSend: func(msg domain.Message) (bool, error) {
				if message != msg {
					t.Errorf("should be the same message, but actual: %v, %v", message, msg)
				}
				return false, fmt.Errorf("err")
			},
		}

		worker := NewWorker("test")

		go func() {
			worker.Add(conn)
			worker.Receive(message)
			worker.Finish()
		}()
		if err := worker.Run(); err != nil {
			t.Errorf("should not be error: %v", err)
		}
	})
}

func TestNotifySendResult(t *testing.T) {

	errMsg := "err"

	t.Run("sent", func(t *testing.T) {
		writer := &bytes.Buffer{}
		log.SetOutput(writer)

		worker := NewWorker("test")

		go func() {
			worker.NotifySendResult(nil)
			worker.Finish()
		}()
		if err := worker.Run(); err != nil {
			t.Errorf("should not be error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		writer := &bytes.Buffer{}
		log.SetOutput(writer)

		worker := NewWorker("test")

		go func() {
			worker.NotifySendResult(fmt.Errorf(errMsg))
			worker.Finish()
		}()
		if err := worker.Run(); err != nil {
			t.Errorf("should not be error: %v", err)
		}

		if got := writer.String(); !strings.Contains(got, errMsg) {
			t.Errorf("should contain errMsg, but actual: %s", got)
		}
	})
}
