package command_test

import (
	"bufio"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestReceive(t *testing.T) {
	server.start()
	defer server.stop()

	cmd := exec.Command("../dist/wsxhub", "--port", insidePort, "receive")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := server.waitToJoin(); err != nil {
		t.Fatal(err)
	}

	received := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			received <- scanner.Text()
			break
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			t.Logf(scanner.Text())
		}
	}()

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

	want := message
	select {
	case got := <-received:
		if got != want {
			t.Errorf("want %v, but %v", want, got)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
