package main

import (
	"bufio"
	"fmt"
	. "gophercises/blackjack/game"
	. "gophercises/deck/deck"
	"os"
	"strconv"
)

func main() {
	game := NewGame(WithShuffling())
	cardGame := game.GetCardGame()
	dealer := Player{
		Id:    1,
		Name:  "Mr. X",
		PType: Dealer,
		Game:  &cardGame,
		Bank:  1000,
	}
	game.AddDealer(dealer)
	player := Player{
		Id:    2,
		Name:  "Player A",
		PType: AI,
		Game:  &cardGame,
		Bank:  1000,
	}
	game.AddPlayer(player)

	fmt.Println("How many times would you like to play? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playCount, _ := strconv.Atoi(scanner.Text())

	game.Play(playCount)
}
