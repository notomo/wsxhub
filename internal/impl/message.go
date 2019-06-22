package impl

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/notomo/wsxhub/internal/domain"
)

// MessageFactoryImpl :
type MessageFactoryImpl struct {
}

// FromBytes :
func (factory *MessageFactoryImpl) FromBytes(bytes []byte) (domain.Message, error) {
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(bytes, &unmarshaled); err != nil {
		return nil, err
	}

	return &MessageImpl{
		bytes:       bytes,
		unmarshaled: unmarshaled,
	}, nil
}

// FromReader :
func (factory *MessageFactoryImpl) FromReader(inputReader io.Reader) (domain.Message, error) {
	bytes, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(bytes, &unmarshaled); err != nil {
		return nil, err
	}

	return &MessageImpl{
		bytes:       bytes,
		unmarshaled: unmarshaled,
	}, nil
}

// MessageImpl :
type MessageImpl struct {
	bytes       []byte
	unmarshaled map[string]interface{}
}

// Bytes :
func (msg *MessageImpl) Bytes() []byte {
	return msg.bytes
}

// Unmarshaled :
func (msg *MessageImpl) Unmarshaled() map[string]interface{} {
	return msg.unmarshaled
}
