package command

import (
	"bytes"
	"io"
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestSendRun(t *testing.T) {
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
			return nil
		},
		FakeReceiveOnce: func(timeout int) ([]byte, error) {
			return []byte(want), nil
		},
	}

	factory := &mock.FakeWebsocketClientFactory{
		FakeClient: func() (domain.WebsocketClient, error) {
			return client, nil
		},
	}

	writer := &bytes.Buffer{}
	cmd := SendCommand{
		WebsocketClientFactory: factory,
		MessageFactory:         msgFactory,
		OutputWriter:           writer,
	}

	if err := cmd.Run(); err != nil {
		t.Fatalf("should not be error: %v", err)
	}

	got := writer.String()
	if got != want {
		t.Errorf("want %v, but %v:", want, got)
	}
}
