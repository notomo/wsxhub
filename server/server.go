package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Server is
type Server struct {
	outsideClients map[string]*Client
	outsideJoined  chan *Client
	outsideLeft    chan *Client
	outsideMessage chan string
	insideClients  map[string]*Client
	insideJoined   chan *Client
	insideLeft     chan *Client
	insideMessage  chan string
	done           chan bool
}

// NewServer is
func NewServer() *Server {
	outsideClients := make(map[string]*Client)
	outsideJoined := make(chan *Client)
	outsideLeft := make(chan *Client)
	outsideMessage := make(chan string)
	insideClients := make(map[string]*Client)
	insideJoined := make(chan *Client)
	insideLeft := make(chan *Client)
	insideMessage := make(chan string)
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
func (server *Server) Receive(client *Client, message string) {
	switch client.clientType {
	case INSIDE:
		server.insideMessage <- message
	case OUTSIDE:
		server.outsideMessage <- message
	}
}

// SendOutside is
func (server *Server) SendOutside(message string) {
	if len(server.outsideClients) == 0 {
		log.Info("Sent to outside but there is no clients: " + message)
		server.SendInside(message)
		return
	}
	for _, client := range server.outsideClients {
		client.Send(message)
	}
	log.Info("Sent to outside: " + message)
}

// SendInside is
func (server *Server) SendInside(message string) {
	for _, client := range server.insideClients {
		client.Send(message)
	}
	log.Info("Sent to inside: " + message)
}

// Listen is
func (server *Server) Listen() {
	onJoined := func(clientType ClientType) func(*websocket.Conn) {
		return func(ws *websocket.Conn) {
			client := NewClient(ws, server, clientType)
			defer client.Close()
			server.Add(client)
			client.Listen()
		}
	}

	outsideWsMux := http.NewServeMux()
	outsideWsMux.Handle("/", websocket.Handler(onJoined(OUTSIDE)))

	insideWsMux := http.NewServeMux()
	insideWsMux.Handle("/", websocket.Handler(onJoined(INSIDE)))

	go func() {
		log.Info("Start outside server")
		err := http.ListenAndServe(":8001", outsideWsMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Info("Start inside server")
		err := http.ListenAndServe(":8002", insideWsMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {

		case client := <-server.outsideJoined:
			log.Info("Joined outside client: " + client.id)
			server.outsideClients[client.id] = client

		case client := <-server.insideJoined:
			log.Info("Joined inside client: " + client.id)
			server.insideClients[client.id] = client

		case client := <-server.outsideLeft:
			log.Info("Left outside client: " + client.id)
			delete(server.outsideClients, client.id)
			client.Close()

		case client := <-server.insideLeft:
			log.Info("Left inside client: " + client.id)
			delete(server.insideClients, client.id)
			client.Close()

		case message := <-server.outsideMessage:
			log.Info("Receive from outside: " + message)
			server.SendInside(message)

		case message := <-server.insideMessage:
			log.Info("Receive from inside: " + message)
			server.SendOutside(message)

		case <-server.done:
			log.Info("Done")
			return
		}
	}
}
