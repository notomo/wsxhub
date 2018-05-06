package server

import (
	"encoding/json"
	"fmt"
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
	filter     *StringMapFilter
}

// NewClient is
func NewClient(ws *websocket.Conn, server *Server, clientType ClientType) *Client {
	id := xid.New()
	done := make(chan bool)
	message := make(chan string)
	filter := NewStringMapFilter(map[string]interface{}{})
	return &Client{id.String(), ws, server, done, message, clientType, filter}
}

// Send is
func (client *Client) Send(message string) {
	client.message <- message
}

// Filtering is
func (client *Client) Filtering(stringMap map[string]interface{}) bool {
	return !client.filter.isSubsetOf(stringMap)
}

func (client *Client) setFilter(message string) {
	log.Info("Set filter for receive message: " + client.id)
	var decodedMessage interface{}
	if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
		panic(err)
	}
	value, ok := decodedMessage.(map[string]interface{})["filter"]
	if !ok {
		return
	}
	client.filter = NewStringMapFilter(value.(map[string]interface{}))
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
			log.Info(fmt.Sprintf("Sent: %s : %s", client.id, message))
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
			log.Info("Done listenReceive: " + client.id)
			return
		default:
			log.Info("Wait in listenReceive: " + client.id)
			var message string
			err := websocket.Message.Receive(client.ws, &message)
			if err == io.EOF || err != nil {
				client.done <- true
				continue
			}
			log.Info("Received in listenReceive: " + client.id + " : " + message)
			if client.clientType == INSIDE {
				client.setFilter(message)
			}
			client.server.Receive(client, message)
		}
	}
}
