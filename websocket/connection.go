package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"battleship-server/data/serializer"
	"battleship-server/game"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	clientBufferSize = 1024
	// Time allowed to write a message to the peer.
	writeWaitMax = 10 * time.Second
	// Maximum message size allowed from peer.
	maxMessagesize = 512
)

type wsHandler struct {
	upgrader   *websocket.Upgrader
	serializer serializer.Serializer
	rooms      map[uint64]*room
}

type client struct {
	accountId   string
	roomId      uint64
	conn        *websocket.Conn
	readBuffer  chan []byte // buffer to store msg from socket
	writeBuffer chan []byte // buffer to write to socket
}

type room struct {
	clients map[string]*client
	game    *game.Game
}

func newWsHandler(serializer serializer.Serializer) *wsHandler {
	upgrader := &websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // remove at prod
	return &wsHandler{
		upgrader:   upgrader,
		serializer: serializer,
		rooms:      map[uint64]*room{},
	}
}

func (h *wsHandler) EstablishWsConnection(c *gin.Context) {
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
	// TODO: need lock
	if room, exists := h.rooms[roomId]; exists {
		if client := room.clients[accountId]; client != nil {
			client.conn.Close()
		}
	}
	// room might have closed, check again, replace with atomic check next time
	r, exists := h.rooms[roomId]
	if !exists {
		fmt.Println("creating room", roomId)
		r = &room{
			clients: map[string]*client{},
		}
		h.rooms[roomId] = r
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &client{
		accountId:   accountId,
		roomId:      roomId,
		conn:        conn,
		writeBuffer: make(chan []byte, clientBufferSize),
		readBuffer:  make(chan []byte, clientBufferSize),
	}
	r.clients[accountId] = client

	h.handleOnConnect(roomId, accountId)
	go h.listenClientConnection(client)
	go h.listenClientWriteBuffer(client)
	go h.listenClientReadBuffer(client)
}

func (h *wsHandler) listenClientConnection(client *client) {
	// reading message from conn, then write to READ-buffer
	defer func() {
		client.conn.Close()
		delete(h.rooms[client.roomId].clients, client.accountId)
		close(client.readBuffer)
		close(client.writeBuffer)
		h.handleOnDisconnect(client.roomId, client.accountId)
		fmt.Println("ended " + client.accountId)
		if len(h.rooms[client.roomId].clients) == 0 {
			delete(h.rooms, client.roomId)
			fmt.Println("closing room", client.roomId)
		}
	}()
	client.conn.SetReadLimit(maxMessagesize)
	for {
		_, b, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(fmt.Sprintf("closed unexpectedly: %+v", err))
			}
			break
		}
		client.readBuffer <- b
	}
}

func (h *wsHandler) listenClientWriteBuffer(client *client) {
	// reading message from client WRITE-buffer & send message
	for {
		b, more := <-client.writeBuffer
		if !more {
			break
		}
		fmt.Println(fmt.Sprintf("out, payload: %+v", string(b)))
		client.conn.SetWriteDeadline(time.Now().Add(writeWaitMax))
		if err := client.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			fmt.Println(fmt.Sprintf("failed to send message from %s, payload: %s, err: %s", client.accountId, string(b), err))
		}
	}
}

func (h *wsHandler) listenClientReadBuffer(client *client) {
	// reading message from client READ-buffer
	// and process accordingly
	for {
		b, more := <-client.readBuffer
		if !more {
			return
		}
		h.handleOnMessage(client.roomId, client.accountId, b)
	}
}
