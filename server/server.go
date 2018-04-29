package server

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Server is
type Server struct {
	outsideClients map[string]*Client
	outsideJoined  chan *Client
	outsideLeft    chan *Client
	outsideMessage chan *Message
	insideClients  map[string]*Client
	insideJoined   chan *Client
	insideLeft     chan *Client
	insideMessage  chan *Message
	done           chan bool
}

// NewServer is
func NewServer() *Server {
	outsideClients := make(map[string]*Client)
	outsideJoined := make(chan *Client)
	outsideLeft := make(chan *Client)
	outsideMessage := make(chan *Message)
	insideClients := make(map[string]*Client)
	insideJoined := make(chan *Client)
	insideLeft := make(chan *Client)
	insideMessage := make(chan *Message)
	done := make(chan bool)
	return &Server{
		outsideClients,
		outsideJoined,
		outsideLeft,
		outsideMessage,
		insideClients,
		insideJoined,
		insideLeft,
		insideMessage,
		done,
	}
}

// Add is
func (server *Server) Add(client *Client) {
	switch client.clientType {
	case INSIDE:
		server.insideJoined <- client
	case OUTSIDE:
		server.outsideJoined <- client
	}
}

// Delete is
func (server *Server) Delete(client *Client) {
	switch client.clientType {
	case INSIDE:
		server.insideLeft <- client
	case OUTSIDE:
		server.outsideLeft <- client
	}
}

// Receive is
func (server *Server) Receive(client *Client, message *Message) {
	switch client.clientType {
	case INSIDE:
		server.insideMessage <- message
	case OUTSIDE:
		server.outsideMessage <- message
	}
}

// SendOutside is
func (server *Server) SendOutside(message *Message) {
	for _, client := range server.outsideClients {
		client.Send(message)
	}
}

// SendInside is
func (server *Server) SendInside(message *Message) {
	for _, client := range server.insideClients {
		client.Send(message)
	}
}

// Listen is
func (server *Server) Listen() {
	onJoined := func(clientType ClientType) func(*websocket.Conn) {
		return func(ws *websocket.Conn) {
			client := NewClient(ws, server, clientType)
			server.Add(client)
			client.Listen()
		}
	}

	outsideWsMux := http.NewServeMux()
	outsideWsMux.Handle("/", websocket.Handler(onJoined(OUTSIDE)))

	insideWsMux := http.NewServeMux()
	insideWsMux.Handle("/", websocket.Handler(onJoined(INSIDE)))

	go func() {
		err := http.ListenAndServe(":8001", outsideWsMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := http.ListenAndServe(":8002", insideWsMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {

		case client := <-server.outsideJoined:
			server.outsideClients[client.id] = client

		case client := <-server.insideJoined:
			server.insideClients[client.id] = client

		case client := <-server.outsideLeft:
			delete(server.outsideClients, client.id)

		case client := <-server.insideLeft:
			delete(server.insideClients, client.id)

		case message := <-server.outsideMessage:
			server.SendInside(message)

		case message := <-server.insideMessage:
			server.SendOutside(message)

		case <-server.done:
			return
		}
	}
}
