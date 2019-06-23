package command

import (
	"github.com/notomo/wsxhub/internal/domain"
)

// ServerCommand :
type ServerCommand struct {
	OutsideServerFactory domain.ServerFactory
	InsideServerFactory  domain.ServerFactory
}

// Run : starts a wsxhub server
// Inside server responds to wsxhub clients.
// Outside server responds to the other clients.
func (cmd *ServerCommand) Run() error {
	outsideServer, err := cmd.OutsideServerFactory.Server(
		domain.NewRoute(
			"/",
			func(conn domain.Connection) error {
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
				return conn.Listen()
			},
		),
	)
	if err != nil {
		return err
	}

	go insideServer.Start()

	return outsideServer.Start()
}
