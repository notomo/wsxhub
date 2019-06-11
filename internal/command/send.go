package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// SendCommand :
type SendCommand struct {
	WebsocketClientFactory domain.WebsocketClientFactory
	OutputWriter           io.Writer
	MessageFactory         domain.MessageFactory
	Timeout                int
}

// Run :
func (cmd *SendCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	message, err := cmd.MessageFactory.Message()
	if err != nil {
		return err
	}

	if err := client.Send(message); err != nil {
		return err
	}

	received, err := client.ReceiveOnce(cmd.Timeout)
	if err != nil {
		return err
	}

	if _, err := cmd.OutputWriter.Write([]byte(received)); err != nil {
		return err
	}

	return nil
}
