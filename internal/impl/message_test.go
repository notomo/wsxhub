package impl

import (
	"bytes"
	"testing"
)

func TestFromReader(t *testing.T) {
	rawMessage := `{"id":"1"}`
	stdin := bytes.NewBufferString(rawMessage)

	factory := MessageFactoryImpl{}
	message, err := factory.FromReader(stdin)
	if err != nil {
		t.Fatalf("should not be error: %v", err)
	}

	{
		got := string(message.Bytes())
		want := rawMessage
		if got != want {
			t.Errorf("want %v, but %v:", want, got)
		}
	}
	{
		got := message.Unmarshaled()["id"]
		want := "1"
		if got != want {
			t.Errorf("want %v, but %v:", want, got)
		}
	}
}
