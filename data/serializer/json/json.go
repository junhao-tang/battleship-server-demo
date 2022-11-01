package json

import (
	"battleship-server/data"
	"encoding/json"
	"fmt"
)

type payload struct {
	Type data.CommunicationId `json:"type"`
	Data string               `json:"data"`
}

type Serializer struct{}

func (Serializer) Marshal(msg data.Message) ([]byte, error) {
	// golang string is utf8
	dataJson, _ := json.Marshal(msg.Data)
	payload := payload{
		Type: msg.Type,
		Data: string(dataJson),
	}
	return json.Marshal(payload)
}

func (Serializer) Unmarshal(b []byte, msg *data.Message) error {
	payload := &payload{}
	err := json.Unmarshal(b, payload)
	if err != nil {
		return err
	}
	var d data.Data
	switch payload.Type {
	case data.CAttack:
		body := data.AttackData{}
		err = json.Unmarshal([]byte(payload.Data), &body)
		d = body
		break
	case data.CPut:
		body := data.PutData{}
		err = json.Unmarshal([]byte(payload.Data), &body)
		d = body
		break
	default:
		return fmt.Errorf("unknown messageType, %v", payload.Type)
	}
	if err != nil {
		return err
	}
	msg.Data = d
	msg.Type = payload.Type
	return nil
}
