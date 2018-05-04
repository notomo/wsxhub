package server

import (
	"io"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
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

// Close is
func (client *Client) Close() {
	log.Info("Close connection: " + client.id)
	client.ws.Close()
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
			client.done <- true
			log.Info("listenSend done: " + client.id)
			return
		}
	}
}

func (client *Client) listenReceive() {
	for {
		select {
		case <-client.done:
			log.Info("listenReceive done: " + client.id)
			return
		default:
			log.Info("Wait in listenReceive: " + client.id)
			var message string
			err := websocket.Message.Receive(client.ws, &message)
			if err == io.EOF || err != nil {
				client.done <- true
				continue
			}
			client.server.Receive(client, message)
		}
	}
}
