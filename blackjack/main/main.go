package main

import (
	"fmt"
	. "gophercises/blackjack/game"
	. "gophercises/deck/deck"
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
	//scanner := bufio.NewScanner(os.Stdin)
	//scanner.Scan()
	//playCount, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < 2; i++ {
		game.DealCards()
		players := game.GetPlayers()
		for _, player := range players {
			fmt.Printf("Player: %s\n", player.String())
			player.DisplayCards()
		}
		dealer := game.GetDealer()
		fmt.Printf("Dealer: %s\n", dealer.String())
		dealer.DisplayCards()

		playersAfterTurn := make([]BlackJackPlayer, 0)
		for _, player := range players {
			player.ExecuteTurn()
			playersAfterTurn = append(playersAfterTurn, player)
		}
		game.UpdatePlayers(playersAfterTurn)

		dealer.ExecuteTurn()
		game.UpdateDealer(dealer)

		var winner BlackJackPlayer
		var err error
		if winner, err = game.DecideWinner(); err != nil {
			fmt.Printf("%s", err.Error())
		}

		if winner.PType == PlayerType {
			fmt.Printf("Winner is: %s!\n", winner.String())
		} else {
			fmt.Printf("Dealer %s won!\n", winner.String())
		}

		game.Reset()
		fmt.Println()
	}

}
