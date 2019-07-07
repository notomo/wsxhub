package command_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSend(t *testing.T) {
	cmdClient := newCommandClient(t, "send")

	cmdClient.startServer()
	defer cmdClient.stopServer()

	if err := cmdClient.cmd.Start(); err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("ws://localhost:%s", outsidePort)
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	msg := `{"id":"1"}`
	cmdClient.writeStdin(msg)

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatal(err)
	}

	sent := cmdClient.scanStdout()

	if err := cmdClient.cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-sent:
		want := string(message)
		if got != want {
			t.Errorf("want %v, but %v", want, got)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestBatchSend(t *testing.T) {
	filter := `{"filters": [{"map": {"id": "1"}}, {"map": {"id": "2"}}]}`
	cmdClient := newCommandClient(t, "send", "--filter", filter)

	cmdClient.startServer()
	defer cmdClient.stopServer()

	if err := cmdClient.cmd.Start(); err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("ws://localhost:%s", outsidePort)
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	msg := `[{"id":"1"},{"id":"2"}]`
	cmdClient.writeStdin(msg)

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatal(err)
	}

	sent := cmdClient.scanStdout()

	if err := cmdClient.cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	want := string(message)
	select {
	case got := <-sent:
		if got != want {
			t.Errorf("want %v, but %v", want, got)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
