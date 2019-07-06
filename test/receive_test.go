package command_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestReceive(t *testing.T) {
	server.start()
	defer server.stop()

	cmdClient := newCommandClient(t, "receive")
	if err := cmdClient.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := server.waitToJoin(); err != nil {
		t.Fatal(err)
	}

	received := cmdClient.scanStdout()

	u := fmt.Sprintf("ws://localhost:%s", outsidePort)
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	message := "{}"
	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-received:
		want := message
		if got != want {
			t.Errorf("want %v, but %v", want, got)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
