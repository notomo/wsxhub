package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// NotifyCommand :
type NotifyCommand struct {
	WebsocketClientFactory domain.WebsocketClientFactory
	MessageFactory         domain.MessageFactory
	InputReader            io.Reader
}

// Run : notifies a message to wsxhub server, but doesn't wait a response
func (cmd *NotifyCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	message, err := cmd.MessageFactory.FromReader(cmd.InputReader)
	if err != nil {
		return err
	}

	if err := client.Send(message.Bytes()); err != nil {
		return err
	}

	return nil
}
