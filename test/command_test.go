package blackbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

const outsidePort = "18881"
const insidePort = "18882"
const id = "1"

func TestSend(t *testing.T) {
	// Server
	err := exec.Command("../dist/wsxhubd", "-outside", outsidePort, "-inside", insidePort).Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 50)

	// Outside client
	url := fmt.Sprintf("ws://localhost:" + outsidePort)
	ws, err := websocket.Dial(url, "", "http://localhost:"+outsidePort)
	if err != nil {
		t.Fatal(err)
	}
	go receiveAndReply(ws)

	// Inside client
	cmd := exec.Command("../dist/wsxhub", "-p", insidePort, "--timeout", "5", "send", "--json", "{}", "--id", id)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			message := scanner.Text()

			var received map[string]interface{}
			if err := json.Unmarshal([]byte(message), &received); err != nil {
				panic(err)
			}

			want := id
			got := received["id"]
			if got != want {
				t.Fatalf("want %v, but %v", want, got)
			}
		}
	}()
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func receiveAndReply(ws *websocket.Conn) {
	var msg string

	err := websocket.Message.Receive(ws, &msg)
	if err != nil {
		panic(err)
	}

	err = websocket.Message.Send(ws, msg)
	if err != nil {
		panic(err)
	}
}
