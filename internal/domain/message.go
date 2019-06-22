package domain

import "io"

// MessageFactory :
type MessageFactory interface {
	FromReader(io.Reader) (Message, error)
	FromBytes([]byte) (Message, error)
}

// Message :
type Message interface {
	Bytes() []byte
	Unmarshaled() map[string]interface{}
}
