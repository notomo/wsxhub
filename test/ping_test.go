package command_test

import (
	"testing"
	"time"
)

func TestPingFailure(t *testing.T) {
	cmdClient := newCommandClient(t, "ping")
	received := cmdClient.scanStderr()

	if err := cmdClient.cmd.Run(); err == nil {
		t.Fatal("`wsxhub ping` must fail if `wsxhub server` is not executed.")
	}
	exitCode := cmdClient.cmd.ProcessState.ExitCode()
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

	cmdClient := newCommandClient(t, "ping")
	received := cmdClient.scanStdout()
	if err := cmdClient.cmd.Run(); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-received:
		want := "pong"
		if got != want {
			t.Errorf("want %v, but %v", want, got)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("output not found")
	}
}
