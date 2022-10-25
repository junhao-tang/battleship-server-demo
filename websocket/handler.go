package websocket

import "fmt"

func (h *WsHandler) handleOnDisconnect(roomId uint64, accountId string) {
}

func (h *WsHandler) handleOnConnect(roomId uint64, accountId string) {
	h.broadcast(roomId, makeJoinGameMessage(accountId)) // broadcast
}

func (h *WsHandler) handleOnMessage(roomId uint64, accountId string, b []byte) {
	msg := &Message{}
	if err := h.Serializer.Unmarshal(b, msg); err != nil {
		fmt.Println(fmt.Sprintf("unknown message from %s, payload: %s, err:%s", accountId, string(b), err))
		return
	}
	switch msg.Type {
	case Join:
		data := msg.Data.(PlayerIdData)
		h.handleJoinMessage(roomId, data.PlayerId)
		return
	}
}

func (h *WsHandler) handleJoinMessage(roomId uint64, accountId string) {
	h.broadcast(roomId, makeJoinGameMessage(accountId))
}

func (h *WsHandler) broadcast(roomId uint64, msg *Message) {
	h.broadcastExclude(roomId, msg, "")
}

func (h *WsHandler) broadcastExclude(roomId uint64, msg *Message, excludedAccountId string) {
	data, err := h.Serializer.Marshal(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("error serializing, %v", msg))
		return
	}
	for accountId, client := range h.Rooms[roomId].Clients {
		if accountId != excludedAccountId {
			client.WriteBuffer <- data
		}
	}
}

func (h *WsHandler) sendTo(roomId uint64, accountId string, msg *Message) {
	data, err := h.Serializer.Marshal(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("error serializing, %v", msg))
		return
	}
	h.Rooms[roomId].Clients[accountId].WriteBuffer <- data
}
