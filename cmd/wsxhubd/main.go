package main

import "github.com/notomo/wsxhub/server"

func main() {
	s := server.NewServer()
	s.Listen()
}
