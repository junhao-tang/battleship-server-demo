package game

import "battleship-server/data"

const (
	gameWidth  = 8
	gameHeight = 8
)

var defaultShips = [...]data.Ship{
	{Width: 1, Height: 4},
	{Width: 4, Height: 1},
	{Width: 2, Height: 2},
	{Width: 3, Height: 3},
}
