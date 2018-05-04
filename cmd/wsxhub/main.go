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
	}

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "Send a request and wait result",
			Action: func(context *cli.Context) error {
				if context.GlobalBool("debug") {
					log.SetLevel(log.DebugLevel)
				}
				c := client.NewClient()
				defer c.Close()
				c.Send()
				c.Receive(false)
				return nil
			},
		},
		{
			Name:  "receive",
			Usage: "Wait receiving requests",
			Action: func(context *cli.Context) error {
				if context.GlobalBool("debug") {
					log.SetLevel(log.DebugLevel)
				}
				c := client.NewClient()
				defer c.Close()
				c.Receive(true)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
