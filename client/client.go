package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Client is Websocket client
type Client struct {
	ws           *websocket.Conn
	done         chan bool
	message      chan string
	filterString string
}

// NewClient is
func NewClient(filterString string) *Client {
	ws, err := websocket.Dial("ws://localhost:8002", "", "http://localhost/")
	if err != nil {
		panic(err)
	}
	log.Debug("Connect")
	done := make(chan bool)
	message := make(chan string)
	return &Client{ws, done, message, filterString}
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
	var requestID = decodedMessage.(map[string]interface{})["requestId"]
	if decodedMessage.(map[string]interface{})["requestId"] == nil {
		requestID = xid.New()
		decodedMessage.(map[string]interface{})["requestId"] = requestID
	}
	var filter = map[string]interface{}{"requestId": requestID}
	decodedMessage.(map[string]interface{})["filter"] = filter
	bytes, err := json.Marshal(decodedMessage)
	if err != nil {
		panic(err)
	}
	var sendMessage = string(bytes)
	log.Debug("Try to send in listenSend: " + sendMessage)
	websocket.Message.Send(client.ws, sendMessage)
}

func (client *Client) listenSend() {
	select {
	case message := <-client.message:
		client.send(message)
		client.done <- true
	case <-client.done:
		log.Debug("Done listenSend")
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
	log.Debug("Read on readStdin: " + message)
	client.message <- message
}

// Close is
func (client *Client) Close() {
	log.Debug("Close")
	client.ws.Close()
}

// Receive is
func (client *Client) Receive(loop bool) {
	go client.writeStdout(loop)
	client.listenReceive(loop)
	<-client.done
}

func (client *Client) listenReceive(loop bool) {
	for {
		select {
		case <-client.done:
			log.Debug("Done listenReceive")
			return
		default:
			log.Debug("Wait on listenReceive")
			var message string
			err := websocket.Message.Receive(client.ws, &message)
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
			log.Debug("Try to write in writeStdout: " + message)
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
			log.Debug("Done writeStdout")
			return
		}
	}
}
