package data

type CommunicationId uint8

const (
	CJoin         CommunicationId = 1
	CStartGame    CommunicationId = 2
	CPut          CommunicationId = 3
	CEnemyPut     CommunicationId = 4
	CAttack       CommunicationId = 5
	CAttackResult CommunicationId = 6
	CEnemyAttack  CommunicationId = 7
)

type Message struct {
	Type CommunicationId
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
	StartingPlayerId string `json:"startingPlayerId"`
	Width            uint8  `json:"width"`
	Height           uint8  `json:"height"`
	Ships            []Ship `json:"ships"`
}

func NewJoinGameMessage(playerId string) Message {
	return Message{
		Type: CJoin,
		Data: PlayerIdData{
			PlayerId: playerId,
		},
	}
}

func NewStartGameMessage(startingPlayerId string, width, height int, ships []Ship) Message {
	return Message{
		Type: CStartGame,
		Data: StartGameData{
			StartingPlayerId: startingPlayerId,
			Width:            uint8(width),
			Height:           uint8(height),
			Ships:            ships,
		},
	}
}

func NewPutShipMessage(shipIndex, index int) Message {
	return Message{
		Type: CPut,
		Data: PutData{
			uint8(index),
			uint8(shipIndex),
		},
	}
}
func NewEnemyPutShipMessage(shipIndex int) Message {
	return Message{
		Type: CEnemyPut,
		Data: PutData{
			ShipIndex: uint8(shipIndex),
		},
	}
}

func NewAttackResultMessage(index int, hit bool) Message {
	return Message{
		Type: CAttackResult,
		Data: AttackData{
			uint8(index),
			hit,
		},
	}
}
func NewEnemyAttackMessage(index int) Message {
	return Message{
		Type: CEnemyAttack,
		Data: AttackData{
			Index: uint8(index),
		},
	}
}
