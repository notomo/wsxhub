package server

// Message is
type Message struct {
	Body string `json:"body"`
}

// String is
func (message *Message) String() string {
	return message.Body
}
