package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// Client is Websocket client
type Client struct {
	ws        *websocket.Conn
	done      chan bool
	message   chan string
	writeErr  chan error
	requestID string
}

// NewClient returns a client
func NewClient(port string, filterString string, keyFilterString string, regexFilterString string, debounceInterval int) (*Client, error) {
	return newClient(port, filterString, keyFilterString, regexFilterString, debounceInterval, "")
}

// NewClientWithID returns a client with the request id
func NewClientWithID(port string, keyFilterString string) (*Client, error) {
	requestID := xid.New().String()
	filterString := fmt.Sprintf("{\"id\":\"%s\"}", requestID)
	return newClient(port, filterString, keyFilterString, "", 0, requestID)
}

func newClient(port string, filterString string, keyFilterString string, regexFilterString string, debounceInterval int, requestID string) (*Client, error) {
	params := url.Values{"filter": {filterString}, "key": {keyFilterString}, "regex": {regexFilterString}, "debounceInterval": {strconv.Itoa(debounceInterval)}}
	url := fmt.Sprintf("ws://localhost:%s/?%s", port, params.Encode())
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	done := make(chan bool)
	message := make(chan string)
	writeErr := make(chan error)
	return &Client{ws, done, message, writeErr, requestID}, nil
}

// Send a message to wsxhubd
func (client *Client) Send(message string) error {
	var decodedMessage interface{}
	if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
		return err
	}
	decodedMessage.(map[string]interface{})["id"] = client.requestID
	bytes, err := json.Marshal(decodedMessage)
	if err != nil {
		return err
	}
	var sendMessage = string(bytes)
	websocket.Message.Send(client.ws, sendMessage)
	return nil
}

// Close the connection
func (client *Client) Close() {
	client.ws.Close()
}

// Receive messages
func (client *Client) Receive(loop bool, timeout int) error {
	go client.writeStdout(loop)
	client.listenReceive(loop, timeout)
	select {
	case writeErr := <-client.writeErr:
		return writeErr
	case <-client.done:
		return nil
	}
}

func (client *Client) listenReceive(loop bool, timeout int) {
	for {
		select {
		case <-client.done:
			return
		default:
			var message string
			if timeout > 0 {
				client.ws.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			}
			err := websocket.Message.Receive(client.ws, &message)
			if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
				client.message <- "{\"error\": {\"data\":{\"name\": \"TimeoutError\"}, \"message\": \"Timeout error\"}}"
				return
			}
			if err == io.EOF || err != nil {
				client.done <- true
				continue
			}
			client.message <- message
			if !loop {
				return
			}
		}
	}
}

func (client *Client) writeStdout(loop bool) {
	for {
		select {
		case message := <-client.message:
			writer := bufio.NewWriter(os.Stdout)
			fmt.Fprintln(writer, message)
			if err := writer.Flush(); err != nil {
				client.writeErr <- err
			}
			if !loop {
				client.done <- true
				return
			}
		case <-client.done:
			return
		}
	}
}
