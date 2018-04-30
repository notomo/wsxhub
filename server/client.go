package server

import (
	"io"

	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// ClientType is
type ClientType int

const (
	// INSIDE is
	INSIDE ClientType = iota
	// OUTSIDE is
	OUTSIDE
)

// Client is Websocket client
type Client struct {
	id         string
	ws         *websocket.Conn
	server     *Server
	done       chan bool
	message    chan string
	clientType ClientType
}

// NewClient is
func NewClient(ws *websocket.Conn, server *Server, clientType ClientType) *Client {
	id := xid.New()
	done := make(chan bool)
	message := make(chan string)
	return &Client{id.String(), ws, server, done, message, clientType}
}

// Send is
func (client *Client) Send(message string) {
	select {
	case client.message <- message:
	default:
		client.server.Delete(client)
	}
}

// Listen is
func (client *Client) Listen() {
	go client.listenSend()
	client.listenReceive()
}

func (client *Client) listenSend() {
	for {
		select {
		case message := <-client.message:
			websocket.Message.Send(client.ws, message)
		case <-client.done:
			client.server.Delete(client)
			return
		}
	}
}

func (client *Client) listenReceive() {
	for {
		select {
		case <-client.done:
			client.server.Delete(client)
			return
		default:
			var message string
			err := websocket.Message.Receive(client.ws, &message)
			client.server.Receive(client, message)
			if err == io.EOF {
				client.done <- true
			}
		}
	}
}
