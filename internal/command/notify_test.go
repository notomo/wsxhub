package command

import (
	"io"
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestNotifyRun(t *testing.T) {
	want := "send"

	msg := &mock.FakeMessage{
		FakeBytes: func() []byte {
			return []byte(want)
		},
	}

	msgFactory := &mock.FakeMessageFactory{
		FakeFromReader: func(inputReader io.Reader) (domain.Message, error) {
			return msg, nil
		},
	}

	client := &mock.FakeWebsocketClient{
		FakeClose: func() error {
			return nil
		},
		FakeSend: func(bytes []byte) error {
			if got := string(bytes); got != want {
				t.Errorf("want %v, but %v:", want, got)
			}
			return nil
		},
	}

	factory := &mock.FakeWebsocketClientFactory{
		FakeClient: func() (domain.WebsocketClient, error) {
			return client, nil
		},
	}

	cmd := NotifyCommand{
		WebsocketClientFactory: factory,
		MessageFactory:         msgFactory,
	}

	if err := cmd.Run(); err != nil {
		t.Fatalf("should not be error: %v", err)
	}
}
