package blackbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

const outsidePort = "18881"
const insidePort = "18882"
const id = "1"

func TestPingFailure(t *testing.T) {
	err := exec.Command("../dist/wsxhub", "--port", insidePort, "ping").Wait()
	if err == nil {
		t.Fatal("`wsxhub ping` must fail if `wsxhubd` is not executed.")
	}
}

func TestSend(t *testing.T) {
	// Server
	err := exec.Command("../dist/wsxhub", "--port", insidePort, "server", "--outside", outsidePort).Start()
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
	cmd := exec.Command("../dist/wsxhub", "--port", insidePort, "send", "--timeout", "5")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
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

	msg := fmt.Sprintf(`{"id":"%s"}`, id)
	if _, err := stdin.Write([]byte(msg)); err != nil {
		t.Fatal(err)
	}
	stdin.Close()

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

func TestPing(t *testing.T) {
	message, err := exec.Command("../dist/wsxhub", "--port", insidePort, "ping").Output()
	if err != nil {
		t.Fatal(err)
	}

	want := "pong"
	got := strings.TrimSpace(string(message))
	if got != want {
		t.Fatalf("want %v, but %v", want, got)
	}
}
