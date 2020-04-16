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
		Id:     1,
		Name:   "Mr. X",
		PType:  DealerType,
		Game:   &cardGame,
		Amount: 1000,
	}
	game.AddDealer(dealer)
	player := Player{
		Id:     2,
		Name:   "Player A",
		PType:  PlayerType,
		Game:   &cardGame,
		Amount: 1000,
	}
	game.AddPlayer(player)

	fmt.Println("How many times would you like to play? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playCount, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < playCount; i++ {
		fmt.Println("How much would you like to bet?")
		scanner.Scan()
		betAmount, _ := strconv.Atoi(scanner.Text())
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
		var loser BlackJackPlayer
		var err error
		if winner, loser, err = game.EarlyOutcome(); err != nil {
			playersAfterTurn := make([]BlackJackPlayer, 0)
			for _, player := range game.GetPlayers() {
				player.ExecuteTurn()
				playersAfterTurn = append(playersAfterTurn, player)
			}
			game.UpdatePlayers(playersAfterTurn)

			dealer.ExecuteTurn()
			game.UpdateDealer(dealer)

			if winner, loser, err = game.EndOfTurnOutcome(); err != nil {
				fmt.Println(err)
			}
		}

		if err == nil {
			winner.UpdateAmount(betAmount)
			loser.UpdateAmount(-1 * betAmount)
			fmt.Printf("Winner is: %s!\n", winner.String())
			if winner.PType == PlayerType {
				game.UpdatePlayers([]BlackJackPlayer{winner})
				game.UpdateDealer(loser)
			} else {
				game.UpdateDealer(winner)
				game.UpdatePlayers([]BlackJackPlayer{loser})
			}
		}
		winner.DisplayAmount()
		loser.DisplayAmount()
		game.Reset()
		fmt.Println()
	}
}
