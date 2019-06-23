package command_test

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

const outsidePort = "18881"
const insidePort = "18882"

var exitCode int
var server *testServer

func TestMain(m *testing.M) {
	server = &testServer{}
	defer func() {
		os.Exit(exitCode)
	}()
	exitCode = m.Run()
}

type testServer struct {
	cmd    *exec.Cmd
	stderr io.Reader
}

func (server *testServer) start() {
	server.cmd = exec.Command("../dist/wsxhub", "--port", insidePort, "server", "--outside", outsidePort)
	stderr, err := server.cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	server.stderr = stderr

	if err := server.cmd.Start(); err != nil {
		panic(err)
	}

	started := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			started <- true
			started <- true
			break
		}
	}()
	<-started
	<-started
}

func (server *testServer) stop() {
	if server.cmd == nil {
		return
	}

	err := server.cmd.Process.Kill()
	if err != nil {
		panic(err)
	}
	server.cmd = nil
}

func (server *testServer) waitToJoin() error {
	joined := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(server.stderr)
		for scanner.Scan() {
			joined <- true
			break
		}
	}()
	select {
	case <-joined:
		return nil
	case <-time.After(1 * time.Second):
		return errors.New("timeout for join")
	}
}
