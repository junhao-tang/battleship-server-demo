package websocket

import (
	"encoding/json"
	"fmt"
)

type payload struct {
	Type CommunicationType `json:"type"`
	Data string            `json:"data"`
}

type Serializer interface {
	Marshal(*Message) ([]byte, error)
	Unmarshal([]byte, *Message) error
}

type JsonSerializer struct{}

func (JsonSerializer) Marshal(msg *Message) ([]byte, error) {
	// golang string is utf8
	dataJson, _ := json.Marshal(msg.Data)
	payload := payload{
		Type: msg.Type,
		Data: string(dataJson),
	}
	return json.Marshal(payload)
}

func (JsonSerializer) Unmarshal(data []byte, msg *Message) error {
	payload := &payload{}
	err := json.Unmarshal(data, payload)
	if err != nil {
		return err
	}
	var d Data
	switch payload.Type {
	case Join:
		d = &PlayerIdData{}
		break
	default:
		return fmt.Errorf("unknown messageType, %v", payload.Type)
	}
	err = json.Unmarshal([]byte(payload.Data), d)
	if err != nil {
		return err
	}
	msg.Type = payload.Type
	msg.Data = d
	return nil
}
