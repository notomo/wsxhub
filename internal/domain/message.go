package domain

// MessageFactory :
type MessageFactory interface {
	Message() (string, error)
}
