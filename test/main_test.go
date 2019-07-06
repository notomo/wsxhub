package command_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	server.cmd = exec.Command("../dist/wsxhub", "--port", insidePort, "server", "--outside", outsidePort, "--outside-allow", "localhost:"+outsidePort)
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
		scanner.Scan()
		fmt.Println(scanner.Text())
		started <- true
		scanner.Scan()
		fmt.Println(scanner.Text())
		started <- true
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
			msg := scanner.Text()
			if strings.Contains(msg, "(inside) joined") {
				joined <- true
				break
			}
		}
	}()
	select {
	case <-joined:
		return nil
	case <-time.After(1 * time.Second):
		return errors.New("timeout for join")
	}
}

type commandClient struct {
	stdout io.ReadCloser
	stderr io.ReadCloser
	stdin  io.WriteCloser
	cmd    *exec.Cmd
	t      *testing.T
}

func newCommandClient(t *testing.T, extendedArgs ...string) *commandClient {
	bin := "../dist/wsxhub"
	baseArgs := []string{"--port", insidePort}
	args := append(baseArgs, extendedArgs...)
	cmd := exec.Command(bin, args...)
	t.Logf("command: %s %s", bin, strings.Join(args, " "))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	return &commandClient{
		stdout: stdout,
		stderr: stderr,
		stdin:  stdin,
		cmd:    cmd,
		t:      t,
	}
}

func (cmdClient *commandClient) scanStderr() chan string {
	received := make(chan string)
	go func() {
		scanner := bufio.NewScanner(cmdClient.stderr)
		for scanner.Scan() {
			msg := scanner.Text()
			cmdClient.t.Logf("scanned: %s", msg)
			received <- msg
			break
		}
	}()
	return received
}

func (cmdClient *commandClient) scanStdout() chan string {
	received := make(chan string)
	go func() {
		scanner := bufio.NewScanner(cmdClient.stdout)
		for scanner.Scan() {
			msg := scanner.Text()
			cmdClient.t.Logf("scanned: %s", msg)
			received <- msg
			break
		}
	}()

	go func() {
		scanner := bufio.NewScanner(cmdClient.stderr)
		for scanner.Scan() {
			cmdClient.t.Logf("stderr: %s", scanner.Text())
		}
	}()

	return received
}

func (cmdClient *commandClient) writeStdin(msg string) {
	if _, err := cmdClient.stdin.Write([]byte(msg)); err != nil {
		panic(err)
	}
	cmdClient.stdin.Close()
	cmdClient.t.Logf("written: %s", msg)
}
