package main

import "github.com/notomo/wsxhub/client"

func main() {
	c := client.NewClient()
	go c.Receive()
	c.Send()
}
