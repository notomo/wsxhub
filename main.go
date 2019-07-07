package main

import (
	"fmt"
	"os"

	"github.com/notomo/wsxhub/internal/command"
	"github.com/notomo/wsxhub/internal/impl"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wsxhub"
	app.Usage = "websocket client from stdio"
	app.Version = "0.0.6"
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
					OutputWriter:   os.Stdout,
					Timeout:        context.Int("timeout"),
					MessageFactory: &impl.MessageFactoryImpl{},
					InputReader:    os.Stdin,
				}
				return cmd.Run()
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
			Name:  "notify",
			Usage: "Send a request, but don't wait response",
			Action: func(context *cli.Context) error {
				cmd := command.NotifyCommand{
					WebsocketClientFactory: &impl.WebsocketClientFactoryImpl{
						Port: context.GlobalString("port"),
					},
					MessageFactory: &impl.MessageFactoryImpl{},
					InputReader:    os.Stdin,
				}
				return cmd.Run()
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
				return cmd.Run()
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
				return cmd.Run()
			},
		},
		{
			Name:  "server",
			Usage: "Start server",
			Action: func(context *cli.Context) error {
				outsideWorker := impl.NewWorker("outside")
				insideWorker := impl.NewWorker("inside")
				filterClauseFactory := &impl.FilterClauseFactoryImpl{}
				messageFactory := &impl.MessageFactoryImpl{}
				port := context.GlobalString("port")
				cmd := command.ServerCommand{
					OutsideServerFactory: &impl.ServerFactoryImpl{
						Port:                context.String("outside"),
						Worker:              outsideWorker,
						TargetWorker:        insideWorker,
						FilterClauseFactory: filterClauseFactory,
						MessageFactory:      messageFactory,
						HostPattern:         context.String("outside-allow"),
					},
					InsideServerFactory: &impl.ServerFactoryImpl{
						Port:                port,
						Worker:              insideWorker,
						TargetWorker:        outsideWorker,
						FilterClauseFactory: filterClauseFactory,
						MessageFactory:      messageFactory,
						HostPattern:         "localhost:" + port,
					},
				}
				return cmd.Run()
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "outside",
					Usage: "port for outside",
					Value: "8001",
				},
				cli.StringFlag{
					Name:  "outside-allow",
					Usage: "allowed request host pattern",
					Value: "localhost:8001",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
