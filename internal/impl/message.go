package impl

import (
	"io"
	"io/ioutil"
)

// MessageFactoryImpl :
type MessageFactoryImpl struct {
	InputReader io.Reader
}

// Message :
func (factory *MessageFactoryImpl) Message() (string, error) {
	bytes, err := ioutil.ReadAll(factory.InputReader)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
