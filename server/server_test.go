package server

import (
	"fmt"
	"net/url"
	"testing"

	"golang.org/x/net/websocket"
)

func insideConnect(filter string, key string) *websocket.Conn {
	params := url.Values{"filter": {filter}, "key": {key}}
	url := fmt.Sprintf("ws://localhost:8002/?%s", params.Encode())
	ws, _ := websocket.Dial(url, "", "http://localhost/")
	return ws
}

func outsideConnect() *websocket.Conn {
	url := fmt.Sprintf("ws://localhost:8001/")
	ws, _ := websocket.Dial(url, "", "http://localhost/")
	return ws
}

func TestSendButOutsideIsEmpty(t *testing.T) {
	go func() {
		s := NewServer()
		s.Listen()
	}()

	inside := insideConnect("", "")

	var want = "{\"body\":1}"
	websocket.Message.Send(inside, want)

	var got string
	websocket.Message.Receive(inside, &got)

	if got != want {
		t.Fatalf("want %q, but %q:", want, got)
	}
	inside.Close()
}

type Message struct {
	Sent     int `json:"sent"`
	Required int `json:"required"`
}

type FilteredMessage struct {
	Sent int `json:"sent"`
}

func TestSendAndReceive(t *testing.T) {
	inside := insideConnect("{\"sent\":1}", "{\"required\":true}")
	outside := outsideConnect()

	var want = 1
	var insideSent = Message{want, 1}
	websocket.JSON.Send(inside, insideSent)

	var outsideGot Message
	websocket.JSON.Receive(outside, &outsideGot)

	if outsideGot.Sent != want {
		t.Fatalf("want %q, but %q:", want, outsideGot.Sent)
	}

	websocket.JSON.Send(outside, Message{2, 1})
	websocket.JSON.Send(outside, FilteredMessage{1})
	websocket.JSON.Send(outside, outsideGot)

	var insideGot Message
	websocket.JSON.Receive(inside, &insideGot)
	if insideGot.Sent != want {
		t.Fatalf("want %q, but %q:", want, insideGot.Sent)
	}
	outside.Close()
	inside.Close()
}

func TestShutdown(t *testing.T) {
	url := fmt.Sprintf("ws://localhost:8002/done")
	ws, _ := websocket.Dial(url, "", "http://localhost/")
	var got Message
	err := websocket.JSON.Receive(ws, got)
	if err == nil {
		t.Fatalf("should be error")
	}
	ws.Close()
}
