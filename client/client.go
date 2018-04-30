package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// Client is Websocket client
type Client struct {
	ws      *websocket.Conn
	done    chan bool
	message chan string
}

// NewClient is
func NewClient() *Client {
	ws, err := websocket.Dial("ws://localhost:8002", "", "http://localhost/")
	if err != nil {
		panic(err)
	}
	done := make(chan bool)
	message := make(chan string)
	return &Client{ws, done, message}
}

// Send is
func (client *Client) Send() {
	go client.listenSend()
	client.readStdin()
}

func (client *Client) listenSend() {
	for {
		select {
		case message := <-client.message:
			var decodedMessage interface{}
			if err := json.Unmarshal([]byte(message), &decodedMessage); err != nil {
				panic(err)
			}
			if decodedMessage.(map[string]interface{})["requestId"] == nil {
				decodedMessage.(map[string]interface{})["requestId"] = xid.New()
				bytes, err := json.Marshal(decodedMessage)
				if err != nil {
					panic(err)
				}
				websocket.Message.Send(client.ws, string(bytes))
			} else {
				websocket.Message.Send(client.ws, message)
			}
		case <-client.done:
			return
		}
	}
}

func (client *Client) readStdin() {
	for {
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
		client.message <- string(buf)
	}
}

// Receive is
func (client *Client) Receive() {
	go client.writeStdout()
	client.listenReceive()
}

func (client *Client) listenReceive() {
	for {
		select {
		case <-client.done:
			return
		default:
			var message string
			err := websocket.Message.Receive(client.ws, &message)
			if err == io.EOF {
				client.done <- true
			}
			client.message <- message
		}
	}
}

func (client *Client) writeStdout() {
	for {
		select {
		case message := <-client.message:
			writer := bufio.NewWriter(os.Stdout)
			fmt.Fprintln(writer, message)
			if err := writer.Flush(); err != nil {
				panic(err)
			}
		case <-client.done:
			return
		}
	}
}
