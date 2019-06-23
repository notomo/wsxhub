package mock

import (
	"io"

	"github.com/notomo/wsxhub/internal/domain"
)

// FakeMessageFactory :
type FakeMessageFactory struct {
	domain.MessageFactory
	FakeFromReader func(io.Reader) (domain.Message, error)
}

// FromReader :
func (factory *FakeMessageFactory) FromReader(inputReader io.Reader) (domain.Message, error) {
	return factory.FakeFromReader(inputReader)
}

// FakeMessage :
type FakeMessage struct {
	domain.Message
	FakeBytes func() []byte
}

// Bytes :
func (factory *FakeMessage) Bytes() []byte {
	return factory.FakeBytes()
}