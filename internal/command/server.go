package command

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// ServerCommand :
type ServerCommand struct {
	OutputWriter         io.Writer
	OutsideServerFactory domain.ServerFactory
	InsideServerFactory  domain.ServerFactory
}

// Run :
func (cmd *ServerCommand) Run() error {
	outsideServer, err := cmd.OutsideServerFactory.Server(
		domain.NewRoute(
			"/",
			func(conn domain.Connection) error {
				defer conn.Close()
				return conn.Listen()
			},
		),
	)
	if err != nil {
		return err
	}

	insideServer, err := cmd.InsideServerFactory.Server(
		domain.NewRoute(
			"/",
			func(conn domain.Connection) error {
				defer conn.Close()
				return conn.Listen()
			},
		),
	)
	if err != nil {
		return err
	}

	go insideServer.Start(
		func(conns map[string]domain.Connection, msg string) error {
			for _, conn := range conns {
				if err := conn.Send(msg); err != nil {
					return err
				}
			}
			return nil
		},
	)

	return outsideServer.Start(
		func(conns map[string]domain.Connection, msg string) error {
			for _, conn := range conns {
				if err := conn.Send(msg); err != nil {
					return err
				}
			}
			return nil
		},
	)
}
