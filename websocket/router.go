package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Router(addr string) *gin.Engine {
	upgrader := &websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // remove at prod
	hdl := &WsHandler{
		Upgrader:   upgrader,
		Serializer: &JsonSerializer{},
		Rooms:      map[uint64]*Room{},
	}
	r := gin.Default()
	r.GET("", hdl.establishWsConnection)
	r.Run(addr)
	return r
}
