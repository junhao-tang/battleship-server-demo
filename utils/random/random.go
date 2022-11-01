package random

import (
	"math/rand"
	"time"
)

var seed = rand.NewSource(time.Now().UnixNano())
var random = rand.New(seed)

func RandInt(max int) int {
	return random.Intn(max)
}

func Shuffle(ls []string) {
	random.Shuffle(len(ls), func(i, j int) { ls[i], ls[j] = ls[j], ls[i] })
}
