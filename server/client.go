package server

import (
	"io"
	"net/url"

	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// ClientType is a targeting client type
type ClientType int

const (
	// INSIDE targets the wsxhub client
	INSIDE ClientType = iota
	// OUTSIDE targets the other client
	OUTSIDE
)

// Client is Websocket client
type Client struct {
	id          string
	ws          *websocket.Conn
	server      *Server
	done        chan bool
	message     chan string
	clientType  ClientType
	filter      *StringMapFilter
	keyFilter   *KeyFilter
	regexFilter *RegexFilter
}

// NewClient creates a client
func NewClient(ws *websocket.Conn, server *Server, clientType ClientType) (*Client, error) {
	id := xid.New()
	done := make(chan bool)
	message := make(chan string)

	filterString := ws.Request().FormValue("filter")
	decodedFilterString, decodedFilterErr := url.QueryUnescape(filterString)
	if decodedFilterErr != nil {
		return nil, decodedFilterErr
	}
	filter, filterErr := NewStringMapFilterFromString(decodedFilterString)
	if filterErr != nil {
		return nil, filterErr
	}

	keyFilterString := ws.Request().FormValue("key")
	var keyFilter *KeyFilter
	decodedKeyFilterString, decodedKeyFilterErr := url.QueryUnescape(keyFilterString)
	if decodedKeyFilterErr != nil {
		return nil, decodedKeyFilterErr
	}
	keyFilter, keyFilterErr := NewKeyFilterFromString(decodedKeyFilterString)
	if keyFilterErr != nil {
		return nil, keyFilterErr
	}

	regexFilterString := ws.Request().FormValue("regex")
	decodedRegexFilterString, decodedRegexFilterErr := url.QueryUnescape(regexFilterString)
	if decodedRegexFilterErr != nil {
		return nil, decodedRegexFilterErr
	}
	regexFilter, regexFilterErr := NewRegexFilterFromString(decodedRegexFilterString)
	if regexFilterErr != nil {
		return nil, regexFilterErr
	}

	return &Client{id.String(), ws, server, done, message, clientType, filter, keyFilter, regexFilter}, nil
}

// Send a message
func (client *Client) Send(message string) {
	client.message <- message
}

// Filtering returns true if stringMap is not match filters
func (client *Client) Filtering(stringMap map[string]interface{}) bool {
	if !client.keyFilter.Match(stringMap) {
		return true
	}
	if !client.filter.isSubsetOf(stringMap) {
		return true
	}
	return !client.regexFilter.Match(stringMap)
}

// Close the connection
func (client *Client) Close() {
	client.ws.Close()
}

// Listen sending and receiving
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
			return
		}
	}
}

func (client *Client) listenReceive() {
	for {
		select {
		case <-client.done:
			return
		default:
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
