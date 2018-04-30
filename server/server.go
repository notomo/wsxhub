package server

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type requestClient struct {
	client    *Client
	RequestID string
}

// Server is
type Server struct {
	outsideClients       map[string]*Client
	outsideJoined        chan *Client
	outsideLeft          chan *Client
	outsideMessage       chan string
	insideClients        map[string]*Client
	insideWaitingClients map[string]*Client
	insideStartedWaiting chan *requestClient
	insideEndedWaiting   chan *requestClient
	insideJoined         chan *Client
	insideLeft           chan *Client
	insideMessage        chan string
	done                 chan bool
}

// NewServer is
func NewServer() *Server {
	outsideClients := make(map[string]*Client)
	outsideJoined := make(chan *Client)
	outsideLeft := make(chan *Client)
	outsideMessage := make(chan string)
	insideClients := make(map[string]*Client)
	insideWaitingClients := make(map[string]*Client)
	insideStartedWaiting := make(chan *requestClient)
	insideEndedWaiting := make(chan *requestClient)
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
		insideWaitingClients,
		insideStartedWaiting,
		insideEndedWaiting,
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
		var decodedMessage requestMessage
		if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
			panic(err)
		}
		server.insideStartedWaiting <- &requestClient{client, decodedMessage.RequestID}
		server.insideMessage <- message
	case OUTSIDE:
		server.outsideMessage <- message
	}
}

// SendOutside is
func (server *Server) SendOutside(message string) {
	for _, client := range server.outsideClients {
		client.Send(message)
	}
}

type requestMessage struct {
	RequestID string `json:"requestId"`
}

// SendInside is
func (server *Server) SendInside(message string) {
	var decodedMessage requestMessage
	if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
		panic(err)
	}
	if decodedMessage.RequestID != "" {
		if client, ok := server.insideWaitingClients[decodedMessage.RequestID]; ok {
			client.Send(message)
			server.insideEndedWaiting <- &requestClient{client, decodedMessage.RequestID}
		}
	} else {
		for _, client := range server.insideClients {
			client.Send(message)
		}
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

		case requestClient := <-server.insideStartedWaiting:
			server.insideWaitingClients[requestClient.RequestID] = requestClient.client

		case requestClient := <-server.insideEndedWaiting:
			delete(server.insideWaitingClients, requestClient.RequestID)

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
