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
		PType: DealerType,
		Game:  &cardGame,
	}
	game.AddDealer(dealer)
	player := Player{
		Id:    2,
		Name:  "Player A",
		PType: PlayerType,
		Game:  &cardGame,
	}
	game.AddPlayer(player)

	fmt.Println("How many times would you like to play? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playCount, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < playCount; i++ {
		game.DealCards()
		playersBeforeTurn := make([]BlackJackPlayer, 0)
		for _, player := range game.GetPlayers() {
			fmt.Printf("Player: %s\n", player.String())
			player.DisplayCards()
			player.ComputeScore()
			playersBeforeTurn = append(playersBeforeTurn, player)
		}
		game.UpdatePlayers(playersBeforeTurn)

		dealer := game.GetDealer()
		fmt.Printf("Dealer: %s\n", dealer.String())
		dealer.DisplayCards()
		dealer.ComputeScore()
		game.UpdateDealer(dealer)

		var winner BlackJackPlayer
		var err error
		if winner, err = game.EarlyWinner(); err != nil {
			playersAfterTurn := make([]BlackJackPlayer, 0)
			for _, player := range game.GetPlayers() {
				player.ExecuteTurn()
				playersAfterTurn = append(playersAfterTurn, player)
			}
			game.UpdatePlayers(playersAfterTurn)

			dealer.ExecuteTurn()
			game.UpdateDealer(dealer)

			if winner, err = game.EndOfTurnWinner(); err != nil {
				fmt.Println(err)
			} else {
				if winner.PType == PlayerType {
					fmt.Printf("Winner is: %s!\n", winner.String())
				} else {
					fmt.Printf("Dealer %s won!\n", winner.String())
				}
			}
		} else {
			fmt.Printf("Winner is: %s!\n", winner.String())
		}
		game.Reset()
		fmt.Println()
	}
}
