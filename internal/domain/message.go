package domain

// MessageFactory :
type MessageFactory interface {
	Message() (Message, error)
	FromBytes([]byte) (Message, error)
}

// Message :
type Message interface {
	Bytes() []byte
	Unmarshaled() map[string]interface{}
}
