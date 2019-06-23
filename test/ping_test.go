package command_test

import (
	"os/exec"
	"testing"
)

func TestPingFailure(t *testing.T) {
	err := exec.Command("../dist/wsxhub", "--port", insidePort, "ping").Wait()
	if err == nil {
		t.Fatal("`wsxhub ping` must fail if `wsxhub server` is not executed.")
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
