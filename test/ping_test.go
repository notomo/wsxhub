package command_test

import (
	"bufio"
	"os/exec"
	"testing"
	"time"
)

func TestPingFailure(t *testing.T) {
	cmd := exec.Command("../dist/wsxhub", "--port", insidePort, "ping")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	received := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			received <- scanner.Text()
			break
		}
	}()

	if err := cmd.Run(); err == nil {
		t.Fatal("`wsxhub ping` must fail if `wsxhub server` is not executed.")
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 1 {
		t.Fatalf("The exit code must be 1, but actual: %d", exitCode)
	}

	select {
	case <-received:
	case <-time.After(1 * time.Second):
		t.Fatalf("stderr output not found")
	}
}

func TestPing(t *testing.T) {
	server.start()
	defer server.stop()

	message, err := exec.Command("../dist/wsxhub", "--port", insidePort, "ping").Output()
	if err != nil {
		t.Fatal(err)
	}

	want := "pong"
	got := string(message)
	if got != want {
		t.Errorf("want %v, but %v", want, got)
	}
}
