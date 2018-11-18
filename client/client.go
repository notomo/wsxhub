package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// Client is Websocket client
type Client struct {
	ws        *websocket.Conn
	done      chan bool
	message   chan string
	requestID string
}

// NewClient returns a client
func NewClient(filterString string, keyFilterString string, regexFilterString string) (*Client, error) {
	return newClient(filterString, keyFilterString, regexFilterString, "")
}

// NewClientWithID returns a client with the request id
func NewClientWithID(keyFilterString string) (*Client, error) {
	requestID := xid.New().String()
	filterString := fmt.Sprintf("{\"id\":\"%s\"}", requestID)
	return newClient(filterString, keyFilterString, "", requestID)
}

func newClient(filterString string, keyFilterString string, regexFilterString string, requestID string) (*Client, error) {
	params := url.Values{"filter": {filterString}, "key": {keyFilterString}, "regex": {regexFilterString}}
	url := fmt.Sprintf("ws://localhost:8002/?%s", params.Encode())
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	done := make(chan bool)
	message := make(chan string)
	return &Client{ws, done, message, requestID}, nil
}

// Send is
func (client *Client) Send(jsons ...string) {
	if len(jsons) == 0 {
		go client.listenSend()
		client.readStdin()
		<-client.done
	} else {
		client.send(jsons[0])
	}
}

func (client *Client) send(message string) {
	var decodedMessage interface{}
	if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
		panic(err)
	}
	decodedMessage.(map[string]interface{})["id"] = client.requestID
	bytes, err := json.Marshal(decodedMessage)
	if err != nil {
		panic(err)
	}
	var sendMessage = string(bytes)
	websocket.Message.Send(client.ws, sendMessage)
}

func (client *Client) listenSend() {
	select {
	case message := <-client.message:
		client.send(message)
		client.done <- true
	case <-client.done:
		return
	}
}

func (client *Client) readStdin() {
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		buf = append(buf, line...)
		if !isPrefix {
			break
		}
	}
	var message = string(buf)
	client.message <- message
}

// Close is
func (client *Client) Close() {
	client.ws.Close()
}

// Receive is
func (client *Client) Receive(loop bool, timeout int) {
	go client.writeStdout(loop)
	client.listenReceive(loop, timeout)
	<-client.done
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
				client.message <- "{\"status\": \"ng\"}"
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
				panic(err)
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
