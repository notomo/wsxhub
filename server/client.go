package server

import (
	"io"
	"net/url"

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

// NewClient is
func NewClient(ws *websocket.Conn, server *Server, clientType ClientType) *Client {
	id := xid.New()
	done := make(chan bool)
	message := make(chan string)
	filterString := ws.Request().FormValue("filter")
	var filter *StringMapFilter
	if decoded, err := url.QueryUnescape(filterString); err == nil {
		filter = NewStringMapFilterFromString(decoded)
	} else {
		panic(err)
	}
	keyFilterString := ws.Request().FormValue("key")
	var keyFilter *KeyFilter
	if decoded, err := url.QueryUnescape(keyFilterString); err == nil {
		keyFilter = NewKeyFilterFromString(decoded)
	} else {
		panic(err)
	}
	regexFilterString := ws.Request().FormValue("regex")
	var regexFilter *RegexFilter
	if decoded, err := url.QueryUnescape(regexFilterString); err == nil {
		regexFilter = NewRegexFilterFromString(decoded)
	} else {
		panic(err)
	}
	return &Client{id.String(), ws, server, done, message, clientType, filter, keyFilter, regexFilter}
}

// Send is
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

// Close is
func (client *Client) Close() {
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
