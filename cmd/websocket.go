package main

import "battleship-server/websocket"

func main() {
	websocket.Router(":80")
}
