package main

import "battleship-server/websocket"

func main() {
	websocket.NewRouter(":80")
}
