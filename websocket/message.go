package websocket

type CommunicationType uint8

const (
	Join         CommunicationType = 1
	StartGame    CommunicationType = 2
	Put          CommunicationType = 3
	EnemyPut     CommunicationType = 4
	Attack       CommunicationType = 5
	AttackResult CommunicationType = 6
	EnemyAttack  CommunicationType = 7
)

type Message struct {
	Type CommunicationType
	Data Data
}

type Data interface {
}

type PutData struct {
	Index     uint8 `json:"index"`
	ShipIndex uint8 `json:"shipIndex"`
}
type PlayerIdData struct {
	PlayerId string `json:"playerId"`
}

type AttackData struct {
	Index uint8 `json:"index"`
	Hit   bool  `json:"hit"`
}

type StartGameData struct {
	PlayerGoFirst bool  `json:"playerGoFirst"`
	Width         uint8 `json:"width"`
	Height        uint8 `json:"height"`
	Ships         []any `json:"ships"`
}

func makeJoinGameMessage(playerId string) *Message {
	return &Message{
		Type: Join,
		Data: PlayerIdData{
			PlayerId: playerId,
		},
	}
}
