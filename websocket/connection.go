package websocket

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	CLIENT_BUFFER_SIZE = 1024
	// Time allowed to write a message to the peer.
	WRITE_WAIT_MAX = 10 * time.Second
	// Maximum message size allowed from peer.
	MAX_MESSAGE_SIZE = 512
)

type WsHandler struct {
	Upgrader   *websocket.Upgrader
	Serializer Serializer
	Rooms      map[uint64]*Room
}

type Client struct {
	AccountId   string
	RoomId      uint64
	Conn        *websocket.Conn
	ReadBuffer  chan []byte // buffer to store msg from socket
	WriteBuffer chan []byte // buffer to write to socket
}

type Room struct {
	Clients map[string]*Client
}

func (h *WsHandler) establishWsConnection(c *gin.Context) {
	accessToken := c.Query("access_token")
	if accessToken == "" {
		c.AbortWithStatusJSON(400, gin.H{"msg": "invalid_param"})
		return
	}
	roomId, err := strconv.ParseUint(c.Query("room_id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "invalid_param"})
		return
	}
	accountId := accessToken

	if room, exists := h.Rooms[roomId]; exists {
		if client := room.Clients[accountId]; client != nil {
			client.Conn.Close()
		}
	}
	// room might have closed, check again, replace with atomic check next time
	room, exists := h.Rooms[roomId]
	if !exists {
		fmt.Println("creating room", roomId)
		room = &Room{
			Clients: map[string]*Client{},
		}
		h.Rooms[roomId] = room
	}
	conn, err := h.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		AccountId:   accountId,
		RoomId:      roomId,
		Conn:        conn,
		WriteBuffer: make(chan []byte, CLIENT_BUFFER_SIZE),
		ReadBuffer:  make(chan []byte, CLIENT_BUFFER_SIZE),
	}
	room.Clients[accountId] = client

	h.handleOnConnect(roomId, accountId)
	go h.listenClientConnection(client)
	go h.listenClientWriteBuffer(client)
	go h.listenClientReadBuffer(client)
}

func (h *WsHandler) listenClientConnection(client *Client) {
	// reading message from conn, then write to READ-buffer
	defer func() {
		client.Conn.Close()
		delete(h.Rooms[client.RoomId].Clients, client.AccountId)
		close(client.ReadBuffer)
		close(client.WriteBuffer)
		h.handleOnDisconnect(client.RoomId, client.AccountId)
		fmt.Println("ended " + client.AccountId)
		if len(h.Rooms[client.RoomId].Clients) == 0 {
			delete(h.Rooms, client.RoomId)
			fmt.Println("closing room", client.RoomId)
		}
	}()
	client.Conn.SetReadLimit(MAX_MESSAGE_SIZE)
	for {
		_, b, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(fmt.Sprintf("closed unexpectedly: %+v", err))
			}
			break
		}
		client.ReadBuffer <- b
	}
}

func (h *WsHandler) listenClientWriteBuffer(client *Client) {
	// reading message from client WRITE-buffer & send message
	for {
		b, more := <-client.WriteBuffer
		if !more {
			break
		}
		fmt.Println(fmt.Sprintf("out, payload: %+v", string(b)))
		client.Conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT_MAX))
		if err := client.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
			fmt.Println(fmt.Sprintf("failed to send message from %s, payload: %s, err: %s", client.AccountId, string(b), err))
		}
	}
}

func (h *WsHandler) listenClientReadBuffer(client *Client) {
	// reading message from client READ-buffer
	// and process accordingly
	for {
		b, more := <-client.ReadBuffer
		if !more {
			return
		}
		h.handleOnMessage(client.RoomId, client.AccountId, b)
	}
}
