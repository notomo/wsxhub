package impl

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/notomo/wsxhub/internal/domain"
)

// MessageFactoryImpl :
type MessageFactoryImpl struct {
}

// FromBytes :
func (factory *MessageFactoryImpl) FromBytes(bytes []byte) (domain.Message, error) {
	var unknown interface{}
	if err := json.Unmarshal(bytes, &unknown); err != nil {
		return nil, err
	}

	if m, ok := unknown.(map[string]interface{}); ok {
		return &MessageImpl{
			bytes:       bytes,
			unmarshaled: []map[string]interface{}{m},
		}, nil
	}

	unknowns, ok := unknown.([]interface{})
	if !ok {
		return nil, fmt.Errorf("message must be map or map[]")
	}

	maps := []map[string]interface{}{}
	for _, u := range unknowns {
		m, ok := u.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("message must be map or map[]")
		}
		maps = append(maps, m)
	}

	return &MessageImpl{
		bytes:       bytes,
		unmarshaled: maps,
	}, nil
}

// FromReader :
func (factory *MessageFactoryImpl) FromReader(inputReader io.Reader) (domain.Message, error) {
	bytes, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}
	return factory.FromBytes(bytes)
}

// MessageImpl :
type MessageImpl struct {
	bytes       []byte
	unmarshaled []map[string]interface{}
}

// Bytes :
func (msg *MessageImpl) Bytes() []byte {
	return msg.bytes
}

// Unmarshaled :
func (msg *MessageImpl) Unmarshaled() []map[string]interface{} {
	return msg.unmarshaled
}
