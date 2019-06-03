package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// SendCommand :
type SendCommand struct {
	WebsocketClientFactory domain.WebsocketClientFactory
	OutputWriter           io.Writer
	Message                string
	Timeout                int
}

// Run :
func (cmd *SendCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Send(cmd.Message); err != nil {
		return err
	}

	message, err := client.ReceiveOnce(cmd.Timeout)
	if err != nil {
		return err
	}

	if _, err := cmd.OutputWriter.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}
