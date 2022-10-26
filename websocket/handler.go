package websocket

import (
	"fmt"

	"battleship-server/data"
	"battleship-server/game"
)

const maxPlayerCount = 2

func (h *wsHandler) handleOnDisconnect(roomId uint64, accountId string) {
	// TODO signal required
}

func (h *wsHandler) handleOnConnect(roomId uint64, accountId string) {
	h.broadcast(roomId, data.NewJoinGameMessage(accountId))

	room := h.rooms[roomId]
	if len(room.clients) == maxPlayerCount {
		playerIds := make([]string, 0, maxPlayerCount)
		for id := range room.clients {
			playerIds = append(playerIds, id)
		}
		game := game.NewGame(playerIds)
		room.game = game
		h.broadcast(roomId, data.NewStartGameMessage(
			game.CurrentTurnPlayerId(),
			game.GameWidth(),
			game.GameHeight(),
			game.Ships(),
		))
	}
}

func (h *wsHandler) handleOnMessage(roomId uint64, accountId string, b []byte) {
	msg := &data.Message{}
	if err := h.serializer.Unmarshal(b, msg); err != nil {
		fmt.Println(fmt.Sprintf("unknown message from %s, payload: %s, err:%s", accountId, string(b), err))
		return
	}
	fmt.Printf("handling: %v", msg)
	switch msg.Type {
	case data.CPut:
		data := msg.Data.(data.PutData)
		h.handlePutMessage(roomId, accountId, int(data.ShipIndex), int(data.Index))
		return
	case data.CAttack:
		data := msg.Data.(data.AttackData)
		h.handleAttackMessage(roomId, accountId, int(data.Index))
		return
	}
}

func (h *wsHandler) handlePutMessage(roomId uint64, accountId string, shipIndex, index int) {
	game := h.rooms[roomId].game
	if game == nil {
		return
	}
	if game.PutShip(accountId, index, shipIndex) {
		h.sendTo(roomId, accountId, data.NewPutShipMessage(shipIndex, index))
		h.broadcastExclude(roomId, accountId, data.NewEnemyPutShipMessage(shipIndex))
	}
}

func (h *wsHandler) handleAttackMessage(roomId uint64, accountId string, index int) {
	game := h.rooms[roomId].game
	if game == nil {
		return
	}
	success, result := game.Attack(accountId, index)
	if success {
		h.sendTo(roomId, accountId, data.NewAttackResultMessage(index, result))
		h.broadcastExclude(roomId, accountId, data.NewEnemyAttackMessage(index))
	}
}

func (h *wsHandler) broadcast(roomId uint64, msg data.Message) {
	h.broadcastExclude(roomId, "", msg)
}

func (h *wsHandler) broadcastExclude(roomId uint64, excludedAccountId string, msg data.Message) {
	data, err := h.serializer.Marshal(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("error serializing, %v", msg))
		return
	}
	for accountId, client := range h.rooms[roomId].clients {
		if accountId != excludedAccountId {
			client.writeBuffer <- data
		}
	}
}

func (h *wsHandler) sendTo(roomId uint64, accountId string, msg data.Message) {
	data, err := h.serializer.Marshal(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("error serializing, %v", msg))
		return
	}
	h.rooms[roomId].clients[accountId].writeBuffer <- data
}
