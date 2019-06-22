package command

import (
	"bytes"
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestPingRun(t *testing.T) {
	client := &mock.FakeWebsocketClient{
		FakeClose: func() error {
			return nil
		},
	}

	factory := &mock.FakeWebsocketClientFactory{
		FakeClient: func() (domain.WebsocketClient, error) {
			return client, nil
		},
	}

	writer := &bytes.Buffer{}
	cmd := PingCommand{
		WebsocketClientFactory: factory,
		OutputWriter:           writer,
	}

	if err := cmd.Run(); err != nil {
		t.Fatalf("should not be error: %v", err)
	}

	got := writer.String()
	want := "pong"
	if got != want {
		t.Errorf("want %v, but %v:", want, got)
	}
}
