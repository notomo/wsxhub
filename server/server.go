package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
		server.SendInside(message)
		return
	}
	for _, client := range server.outsideClients {
		client.Send(message)
	}
}

// SendInside is
func (server *Server) SendInside(message string) {
	var stringMap map[string]interface{}
	if err := json.Unmarshal([]byte(message), &stringMap); err != nil {
		panic(err)
	}
	for _, client := range server.insideClients {
		if client.Filtering(stringMap) {
			continue
		}
		client.Send(message)
	}
}

// Listen requests
func (server *Server) Listen() {
	onJoined := func(clientType ClientType) func(*websocket.Conn) {
		return func(ws *websocket.Conn) {
			client, err := NewClient(ws, server, clientType)
			if err != nil {
				ws.Close()
				log.Error(err)
				return
			}
			defer client.Close()
			server.Add(client)
			client.Listen()
		}
	}

	outsideWsMux := http.NewServeMux()
	outsideWsMux.Handle("/", websocket.Handler(onJoined(OUTSIDE)))

	insideWsMux := http.NewServeMux()
	insideWsMux.Handle("/", websocket.Handler(onJoined(INSIDE)))
	insideWsMux.Handle("/done", websocket.Handler(func(ws *websocket.Conn) {
		ws.Close()
		server.done <- true
	}))

	outsideServer := &http.Server{Addr: ":8001", Handler: outsideWsMux}
	go func() {
		log.Info("Start outside server")
		outsideServer.ListenAndServe()
	}()

	insideServer := &http.Server{Addr: ":8002", Handler: insideWsMux}
	go func() {
		log.Info("Start inside server")
		insideServer.ListenAndServe()
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
			server.SendInside(message)

		case message := <-server.insideMessage:
			server.SendOutside(message)

		case <-server.done:
			log.Info("Done")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			insideServer.Shutdown(ctx)
			outsideServer.Shutdown(ctx)
			return
		}
	}
}
