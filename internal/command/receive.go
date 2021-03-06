package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// ReceiveCommand :
type ReceiveCommand struct {
	WebsocketClientFactory domain.WebsocketClientFactory
	OutputWriter           io.Writer
	Timeout                int
}

// Run : outputs the received messages
func (cmd *ReceiveCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Receive(cmd.Timeout, func(message []byte) error {
		_, err := cmd.OutputWriter.Write(append(message, '\n'))
		return err
	}); err != nil {
		return err
	}

	return nil
}
