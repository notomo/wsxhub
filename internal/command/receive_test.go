package command

import (
	"bytes"
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestReceiveRun(t *testing.T) {
	message := "received"
	client := &mock.FakeWebsocketClient{
		FakeClose: func() error {
			return nil
		},
		FakeReceive: func(timeout int, callback func([]byte) error) error {
			return callback([]byte(message))
		},
	}

	factory := &mock.FakeWebsocketClientFactory{
		FakeClient: func() (domain.WebsocketClient, error) {
			return client, nil
		},
	}

	writer := &bytes.Buffer{}
	cmd := ReceiveCommand{
		WebsocketClientFactory: factory,
		OutputWriter:           writer,
	}

	if err := cmd.Run(); err != nil {
		t.Fatalf("should not be error: %v", err)
	}

	got := writer.String()
	want := message + "\n"
	if got != want {
		t.Errorf("want %v, but %v:", want, got)
	}
}
