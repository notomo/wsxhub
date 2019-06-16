package main

import (
	"os"

	"github.com/notomo/wsxhub/internal/command"
	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/impl"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wsxhub"
	app.Usage = "websocket client from stdio"
	app.Version = "0.0.5"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Usage: "Set port",
			Value: "8002",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "Send a request and wait result",
			Action: func(context *cli.Context) error {
				cmd := command.SendCommand{
					WebsocketClientFactory: &impl.WebsocketClientFactoryImpl{
						Port:         context.GlobalString("port"),
						FilterSource: context.String("filter"),
					},
					OutputWriter: os.Stdout,
					Timeout:      context.Int("timeout"),
					MessageFactory: &impl.MessageFactoryImpl{
						InputReader: os.Stdin,
					},
				}
				if err := cmd.Run(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "filter",
					Usage: "Filter received json",
				},
				cli.IntFlag{
					Name:  "timeout",
					Usage: "Timeout seconds for receiving",
					Value: 0,
				},
			},
		},
		{
			Name:  "receive",
			Usage: "Wait receiving requests",
			Action: func(context *cli.Context) error {
				factory := &impl.WebsocketClientFactoryImpl{
					Port:         context.GlobalString("port"),
					FilterSource: context.String("filter"),
					Debounce:     context.Int("debounce"),
				}
				cmd := command.ReceiveCommand{
					WebsocketClientFactory: factory,
					OutputWriter:           os.Stdout,
					Timeout:                context.Int("timeout"),
				}
				if err := cmd.Run(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "debounce",
					Usage: "Debounce interval(ms)",
					Value: 0,
				},
				cli.StringFlag{
					Name:  "filter",
					Usage: "Filter received json",
				},
				cli.IntFlag{
					Name:  "timeout",
					Usage: "Timeout seconds for receiving",
					Value: 0,
				},
			},
		},
		{
			Name:  "ping",
			Usage: "Test request to wsxhubd",
			Action: func(context *cli.Context) error {
				factory := &impl.WebsocketClientFactoryImpl{
					Port: context.GlobalString("port"),
				}
				cmd := command.PingCommand{
					WebsocketClientFactory: factory,
					OutputWriter:           os.Stdout,
				}
				if err := cmd.Run(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		{
			Name:  "server",
			Usage: "Start server",
			Action: func(context *cli.Context) error {
				outsideWorker := &impl.WorkerImpl{
					Name:         "outside",
					Joined:       make(chan domain.Connection),
					Received:     make(chan string),
					Left:         make(chan domain.Connection),
					Done:         make(chan bool),
					Conns:        make(map[string]domain.Connection),
					OutputWriter: os.Stdout,
				}
				insideWorker := &impl.WorkerImpl{
					Name:         "inside",
					Joined:       make(chan domain.Connection),
					Received:     make(chan string),
					Left:         make(chan domain.Connection),
					Done:         make(chan bool),
					Conns:        make(map[string]domain.Connection),
					OutputWriter: os.Stdout,
				}
				filterClauseFactory := &impl.FilterClauseFactoryImpl{}
				cmd := command.ServerCommand{
					OutputWriter: os.Stdout,
					OutsideServerFactory: &impl.ServerFactoryImpl{
						Port:                context.String("outside"),
						Worker:              outsideWorker,
						TargetWorker:        insideWorker,
						FilterClauseFactory: filterClauseFactory,
						OutputWriter:        os.Stdout,
					},
					InsideServerFactory: &impl.ServerFactoryImpl{
						Port:                context.GlobalString("port"),
						Worker:              insideWorker,
						TargetWorker:        outsideWorker,
						FilterClauseFactory: filterClauseFactory,
						OutputWriter:        os.Stdout,
					},
				}
				if err := cmd.Run(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "outside",
					Usage: "port for outside",
					Value: "8001",
				},
			},
		},
	}

	app.Run(os.Args)
}
