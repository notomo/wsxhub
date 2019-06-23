package command_test

import (
	"bufio"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSend(t *testing.T) {
	server.start()
	defer server.stop()

	cmd := exec.Command("../dist/wsxhub", "--port", insidePort, "send")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("ws://localhost:%s", outsidePort)
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	id := "1"
	msg := fmt.Sprintf(`{"id":"%s"}`, id)
	if _, err := stdin.Write([]byte(msg)); err != nil {
		t.Fatal(err)
	}
	stdin.Close()

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatal(err)
	}

	sent := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sent <- scanner.Text()
			break
		}
	}()

	if err := cmd.Wait(); err != nil {
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