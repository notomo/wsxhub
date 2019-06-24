package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// PingCommand :
type PingCommand struct {
	WebsocketClientFactory domain.WebsocketClientFactory
	OutputWriter           io.Writer
}

// Run : confirms connection and then outputs "pong"
func (cmd *PingCommand) Run() error {
	client, err := cmd.WebsocketClientFactory.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	if _, err := cmd.OutputWriter.Write([]byte("pong")); err != nil {
		return err
	}

	return nil
}
