package main

import (
	. "gophercises/blackjack/game"
	. "gophercises/deck/deck"
)

func main() {
	game := NewGame(WithShuffling())
	game.Play()
}
