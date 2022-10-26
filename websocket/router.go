package websocket

import (
	"battleship-server/data/serializer/json"

	"github.com/gin-gonic/gin"
)

func NewRouter(addr string) *gin.Engine {
	hdl := newWsHandler(json.Serializer{})
	r := gin.Default()
	r.GET("", hdl.EstablishWsConnection)
	r.Run(addr)
	return r
}
