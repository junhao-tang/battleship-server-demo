package game

import (
	"battleship-server/data"
	"battleship-server/utils/random"
	"fmt"
)

type phase uint8

const (
	pInitialize phase = iota
	pPosition
	pAttack
)

type Game struct {
	playerIds      []string
	ships          []data.Ship
	sea            map[string]*sea
	shipsAvailable map[string]map[int]bool
	phase          phase
	playerTurn     int
}

func NewGame(playerIds []string) *Game {
	s := map[string]*sea{}
	shipsAvailable := map[string]map[int]bool{}
	ships := defaultShips

	random.Shuffle(playerIds)
	for _, id := range playerIds {
		s[id] = newSea(
			gameWidth,
			gameHeight,
		)
		shipsAvailable[id] = map[int]bool{}
		for i := range ships {
			shipsAvailable[id][i] = true
		}
	}
	return &Game{
		sea:            s,
		ships:          ships[:],
		shipsAvailable: shipsAvailable,
		phase:          pPosition,
		playerIds:      playerIds,
		playerTurn:     random.RandInt(len(playerIds)),
	}
}

func (game *Game) PutShip(playerId string, index, shipIndex int) bool {
	if shipIndex < 0 || shipIndex >= len(game.ships) {
		fmt.Printf("put failed, parameter, %s, %d, %d\n", playerId, index, shipIndex)
		return false
	}
	if game.phase != pPosition {
		fmt.Printf("put failed, not phase, %s, %d, %d\n", playerId, index, shipIndex)
		return false
	}
	if !game.shipsAvailable[playerId][shipIndex] {
		fmt.Printf("put failed, not available, %s, %d, %d\n", playerId, index, shipIndex)
		return false
	}
	sea := game.sea[playerId]
	ship := game.ships[shipIndex]
	if !sea.canPut(index, ship.Width, ship.Height) {
		fmt.Printf("put failed, sea check, %s, %d, %d\n", playerId, index, shipIndex)
		return false
	}
	sea.put(index, ship.Width, ship.Height)
	delete(game.shipsAvailable[playerId], shipIndex)

	allUsed := true
	for _, m := range game.shipsAvailable {
		if len(m) != 0 {
			allUsed = false
			break
		}
	}
	if allUsed {
		game.phase = pAttack
	}
	return true
}

func (game *Game) Attack(playerId string, index int) (bool, bool) {
	if game.phase != pAttack {
		fmt.Printf("attack failed, not attack phase, %s, %d\n", playerId, index)
		return false, false
	}
	if game.playerIds[game.playerTurn] != playerId {
		fmt.Printf("attack failed, not player's turn, %s, %d\n", playerId, index)
		return false, false
	}
	// get enemy sea
	var sea *sea
	for pid, s := range game.sea {
		if pid != playerId {
			sea = s
			break
		}
	}
	if !sea.canAttack(index) {
		fmt.Printf("attack failed, sea check, %s, %d\n", playerId, index)
		return false, false
	}
	hit := sea.attack(index)
	game.playerTurn = (game.playerTurn + 1) % len(game.playerIds)
	return true, hit
}

func (game *Game) Ships() []data.Ship {
	return game.ships
}

func (game *Game) CurrentTurnPlayerId() string {
	return game.playerIds[game.playerTurn]
}

func (game *Game) GameWidth() int {
	return game.sea[game.playerIds[0]].width
}

func (game *Game) GameHeight() int {
	return game.sea[game.playerIds[0]].height
}
