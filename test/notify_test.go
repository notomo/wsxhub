package command_test

import (
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
)

func TestNotify(t *testing.T) {
	server.start()
	defer server.stop()

	cmdClient := newCommandClient(t, "notify")
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

	if err := cmdClient.cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	got := string(message)
	want := msg
	if got != want {
		t.Errorf("want %v, but %v", want, got)
	}
}
