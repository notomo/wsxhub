package main

import (
	"os"

	"github.com/notomo/wsxhub/client"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wsxhub"
	app.Usage = "websocket client from stdio"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Show debug messages",
		},
		cli.StringFlag{
			Name:  "regex, r",
			Usage: "Filter received json value by regular expression",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "Filter received json key",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "filter, f",
			Usage: "Filter received json",
			Value: "{}",
		},
		cli.IntFlag{
			Name:  "timeout, t",
			Usage: "Timeout seconds for receiving",
			Value: 0,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "Send a request and wait result",
			Action: func(context *cli.Context) error {
				if context.GlobalBool("debug") {
					log.SetLevel(log.DebugLevel)
				}
				c, err := client.NewClientWithID(context.GlobalString("key"))
				if err != nil {
					return err
				}
				defer c.Close()
				sendErr := c.Send(context.String("json"))
				if sendErr != nil {
					return cli.NewExitError("input json parse error", 1)
				}

				receiveErr := c.Receive(false, context.GlobalInt("timeout"))
				if receiveErr != nil {
					log.Error(receiveErr)
					return receiveErr
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json, j",
					Usage: "Sent json",
				},
			},
		},
		{
			Name:  "receive",
			Usage: "Wait receiving requests",
			Action: func(context *cli.Context) error {
				if context.GlobalBool("debug") {
					log.SetLevel(log.DebugLevel)
				}
				c, err := client.NewClient(context.GlobalString("filter"), context.GlobalString("key"), context.GlobalString("regex"))
				if err != nil {
					return err
				}
				defer c.Close()
				receiveErr := c.Receive(true, context.GlobalInt("timeout"))
				if receiveErr != nil {
					log.Error(receiveErr)
					return receiveErr
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
