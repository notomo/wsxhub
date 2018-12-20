package main

import (
	"os"

	"github.com/notomo/wsxhub/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wsxhubd"
	app.Usage = "websocket server"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "outside, o",
			Usage: "port for outside",
			Value: "8001",
		},
		cli.StringFlag{
			Name:  "inside, i",
			Usage: "port for inside",
			Value: "8002",
		},
	}

	app.Action = func(context *cli.Context) error {
		s := server.NewServer(context.GlobalString("outside"), context.GlobalString("inside"))
		s.Listen()
		return nil
	}

	app.Run(os.Args)
}
