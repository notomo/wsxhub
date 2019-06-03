package main

import (
	"os"

	"github.com/notomo/wsxhub/internal/command"
	"github.com/notomo/wsxhub/internal/impl"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wsxhub"
	app.Usage = "websocket client from stdio"
	app.Version = "0.0.5"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Show debug messages",
		},
		cli.StringFlag{
			Name:  "regex",
			Usage: "Filter received json value by regular expression",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "Filter received json key",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "filter",
			Usage: "Filter received json",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "port",
			Usage: "Set port",
			Value: "8002",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "Timeout seconds for receiving",
			Value: 0,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "Send a request and wait result",
			Action: func(context *cli.Context) error {
				factory := &impl.WebsocketClientFactoryImpl{
					Port: context.GlobalString("port"),
				}
				cmd := command.SendCommand{
					WebsocketClientFactory: factory,
					OutputWriter:           os.Stdout,
					Timeout:                context.GlobalInt("timeout"),
					Message:                context.String("json"),
				}
				if err := cmd.Run(); err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "Sent json",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "Set id",
					Value: "",
				},
			},
		},
		{
			Name:  "receive",
			Usage: "Wait receiving requests",
			Action: func(context *cli.Context) error {
				factory := &impl.WebsocketClientFactoryImpl{
					Port: context.GlobalString("port"),
				}
				cmd := command.ReceiveCommand{
					WebsocketClientFactory: factory,
					OutputWriter:           os.Stdout,
					Timeout:                context.GlobalInt("timeout"),
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
	}

	app.Run(os.Args)
}
