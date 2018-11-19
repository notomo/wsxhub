package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Server represents the websocket server
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

// NewServer creates a server
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

// Add a client to this server
func (server *Server) Add(client *Client) {
	switch client.clientType {
	case INSIDE:
		server.insideJoined <- client
	case OUTSIDE:
		server.outsideJoined <- client
	}
}

// Delete a client from this server
func (server *Server) Delete(client *Client) {
	switch client.clientType {
	case INSIDE:
		server.insideLeft <- client
	case OUTSIDE:
		server.outsideLeft <- client
	}
}

// Receive a message from client
func (server *Server) Receive(client *Client, message string) {
	switch client.clientType {
	case INSIDE:
		server.insideMessage <- message
	case OUTSIDE:
		server.outsideMessage <- message
	}
}

// SendOutside sends a message to the outside client
func (server *Server) SendOutside(message string) error {
	if len(server.outsideClients) == 0 {
		return server.SendInside(message)
	}
	for _, client := range server.outsideClients {
		client.Send(message)
	}
	return nil
}

// SendInside sends a message to the inside client
func (server *Server) SendInside(message string) error {
	var stringMap map[string]interface{}
	if err := json.Unmarshal([]byte(message), &stringMap); err != nil {
		return err
	}
	for _, client := range server.insideClients {
		if client.Filtering(stringMap) {
			continue
		}
		client.Send(message)
	}

	return nil
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
			err := server.SendInside(message)
			if err != nil {
				log.Error(err)
			}

		case message := <-server.insideMessage:
			err := server.SendOutside(message)
			if err != nil {
				log.Error(err)
			}

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
