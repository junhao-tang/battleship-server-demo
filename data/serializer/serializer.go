package serializer

import "battleship-server/data"

type Serializer interface {
	Marshal(data.Message) ([]byte, error)
	Unmarshal([]byte, *data.Message) error
}
