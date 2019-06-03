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

// Run :
func (cmd *ReceiveCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Receive(cmd.Timeout, func(message string) error {
		if _, err := cmd.OutputWriter.Write([]byte(message)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
