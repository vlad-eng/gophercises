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
		Bank:  1000,
	}
	game.AddDealer(dealer)
	player := Player{
		Id:    2,
		Name:  "Player A",
		PType: PlayerType,
		Game:  &cardGame,
		Bank:  1000,
	}
	game.AddPlayer(player)

	fmt.Println("How many times would you like to play? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playCount, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < playCount; i++ {
		game.PlaceBets()
		game.DealCards()
		playersBeforeTurn := make([]BlackJackPlayer, 0)
		for _, player := range game.GetPlayers() {
			fmt.Printf("Player: %s\n", player.String())
			player.DisplayCards(true)
			player.ComputeScore()
			playersBeforeTurn = append(playersBeforeTurn, player)
		}
		game.UpdatePlayers(playersBeforeTurn)

		dealer := game.GetDealer()
		fmt.Printf("Dealer: %s\n", dealer.String())
		dealer.DisplayCards(false)
		dealer.ComputeScore()
		game.UpdateDealer(dealer)

		var winner BlackJackPlayer
		var nonWinner BlackJackPlayer
		var err error
		if winner, nonWinner, err = game.EarlyOutcome(); err != nil {
			playersAfterTurn := make([]BlackJackPlayer, 0)
			for _, player := range game.GetPlayers() {
				player.ExecuteTurn()
				playersAfterTurn = append(playersAfterTurn, player)
			}
			game.UpdatePlayers(playersAfterTurn)

			dealer.ExecuteTurn()
			game.UpdateDealer(dealer)

			if winner, nonWinner, err = game.EndOfTurnOutcome(); err != nil {
				fmt.Println(err)
			}
		}

		if err == nil {
			winner.WinBankUpdate()
			nonWinner.LossBankUpdate()
			fmt.Printf("Winner is: %s!\n", winner.String())
			if winner.PType == PlayerType {
				game.UpdatePlayers([]BlackJackPlayer{winner})
				game.UpdateDealer(nonWinner)
			} else {
				game.UpdateDealer(winner)
				game.UpdatePlayers([]BlackJackPlayer{nonWinner})
			}
			winner.DisplayAmount()
			nonWinner.DisplayAmount()
		}
		game.Reset()
		fmt.Println()
	}
}
