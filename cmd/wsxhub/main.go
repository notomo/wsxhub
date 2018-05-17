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
			Name:  "key, k",
			Usage: "Filter received json key",
			Value: "{}",
		},
		cli.StringFlag{
			Name:  "filter, f",
			Usage: "Filter received json",
			Value: "{}",
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
				c := client.NewClientWithID(context.GlobalString("key"))
				defer c.Close()
				c.Send(context.String("json"))
				c.Receive(false)
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
				c := client.NewClient(context.GlobalString("filter"), context.GlobalString("key"))
				defer c.Close()
				c.Receive(true)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
