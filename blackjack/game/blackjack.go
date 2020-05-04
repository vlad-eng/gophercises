package game

import (
	"bufio"
	"fmt"
	. "gophercises/deck/deck"
	"os"
	"strconv"
)

const DealtCardsCount = 2
const FaceCardValue = 10
const BlackJackMaxScore = 21

type CardGame struct {
	dealer  Player
	players []Player
	deck    []Card
}

//TODO: Refactor this to use the embedding pattern for CardGame & BlackJack, maybe also for Player & BlackJackPlayer
type BlackJack CardGame

func NewGame(gameOptions ...func([]Card) []Card) *BlackJack {
	game := BlackJack{
		deck: New(gameOptions...),
	}
	return &game
}

func (g *BlackJack) Play() {
	cardGame := g.GetCardGame()
	dealer := Player{
		Id:    1,
		Name:  "Mr. X",
		PType: DealerType,
		Game:  &cardGame,
		Bank:  1000,
	}
	g.AddDealer(dealer)
	player := Player{
		Id:    2,
		Name:  "Player A",
		PType: PlayerType,
		Game:  &cardGame,
		Bank:  1000,
	}
	g.AddPlayer(player)

	fmt.Println("How many times would you like to play? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playCount, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < playCount; i++ {
		g.PlaceBets()
		g.DealCards()
		playersBeforeTurn := make([]BlackJackPlayer, 0)
		for _, player := range g.GetPlayers() {
			fmt.Printf("Player: %s\n", player.String())
			player.DisplayCards(true)
			player.ComputeScore()
			playersBeforeTurn = append(playersBeforeTurn, player)
		}
		g.UpdatePlayers(playersBeforeTurn)

		dealer := g.GetDealer()
		fmt.Printf("Dealer: %s\n", dealer.String())
		dealer.DisplayCards(false)
		dealer.ComputeScore()
		g.UpdateDealer(dealer)

		var winner BlackJackPlayer
		var nonWinner BlackJackPlayer
		var err error
		if winner, nonWinner, err = g.EarlyOutcome(); err != nil {
			playersAfterTurn := make([]BlackJackPlayer, 0)
			for _, player := range g.GetPlayers() {
				player.ExecuteTurn()
				playersAfterTurn = append(playersAfterTurn, player)
			}
			g.UpdatePlayers(playersAfterTurn)

			dealer.ExecuteTurn()
			g.UpdateDealer(dealer)

			if winner, nonWinner, err = g.EndOfTurnOutcome(); err != nil {
				fmt.Println(err)
			}
		}

		if err == nil {
			winner.WinBankUpdate()
			nonWinner.LossBankUpdate()
			fmt.Printf("Winner is: %s!\n", winner.String())
			if winner.PType == PlayerType {
				g.UpdatePlayers([]BlackJackPlayer{winner})
				g.UpdateDealer(nonWinner)
			} else {
				g.UpdateDealer(winner)
				g.UpdatePlayers([]BlackJackPlayer{nonWinner})
			}
			winner.DisplayAmount()
			nonWinner.DisplayAmount()
		}
		g.Reset()
		fmt.Println()
	}
}

func (g *BlackJack) GetCardGame() CardGame {
	return CardGame(*g)
}

func (g *BlackJack) GetDealer() BlackJackPlayer {
	return BlackJackPlayer(g.dealer)
}

func (g *BlackJack) GetPlayers() []BlackJackPlayer {
	bPlayers := make([]BlackJackPlayer, 0)
	for _, player := range g.players {
		bPlayer := BlackJackPlayer(player)
		bPlayers = append(bPlayers, bPlayer)
	}
	return bPlayers
}

func (g *BlackJack) GetPlayer(playerIndex int) *BlackJackPlayer {
	p := BlackJackPlayer(g.players[playerIndex])
	return &p
}

func (g *BlackJack) UpdatePlayers(blackJackPlayers []BlackJackPlayer) {
	for i := range g.players {
		g.players[i] = Player(blackJackPlayers[i])
	}
}

func (g *BlackJack) UpdateDealer(blackJackDealer BlackJackPlayer) {
	g.dealer = Player(blackJackDealer)
	for _, player := range g.players {
		//only works for one player against dealer
		g.dealer.BetAmount = player.BetAmount
	}
}

func (g *BlackJack) GetPlayerCount() int {
	return len(g.players)
}

func (g *BlackJack) PlaceBets() {
	for i, player := range g.players {
		p := BlackJackPlayer(player)
		p.placeBet()
		//only works for one player against dealer
		g.dealer.BetAmount += p.BetAmount
		g.players[i] = Player(p)
	}
}

func (g *BlackJack) DealCards() error {
	var err error
	for i, player := range g.players {
		p := BlackJackPlayer(player)
		if err = p.dealCards(); err != nil {
			return err
		}
		g.players[i] = Player(p)
	}

	dealerPlayer := BlackJackPlayer(g.dealer)
	if dealerPlayer.dealCards(); err != nil {
		return err
	}
	g.dealer = Player(dealerPlayer)

	return nil
}

func (g *BlackJack) EarlyOutcome() (BlackJackPlayer, BlackJackPlayer, error) {
	scores, maxScore := g.computeScores()
	if maxScore == BlackJackMaxScore {
		return scores[maxScore], scores[g.dealer.Score], nil
	} else if g.dealer.Score == BlackJackMaxScore {
		return scores[g.dealer.Score], scores[maxScore], nil
	} else {
		return BlackJackPlayer{}, BlackJackPlayer{}, fmt.Errorf("no winner yet")
	}
}

func (g *BlackJack) EndOfTurnOutcome() (BlackJackPlayer, BlackJackPlayer, error) {
	scores, maxScore := g.computeScores()
	if maxScore == g.dealer.Score && maxScore <= BlackJackMaxScore {
		return BlackJackPlayer{}, BlackJackPlayer{}, fmt.Errorf("push: no winner")
	}
	if (maxScore < g.dealer.Score && (g.dealer.Score <= BlackJackMaxScore)) ||
		maxScore > BlackJackMaxScore {

		return BlackJackPlayer(g.dealer), scores[maxScore], nil
	}
	return scores[maxScore], BlackJackPlayer(g.dealer), nil
}

func (g *BlackJack) computeScores() (map[int]BlackJackPlayer, int) {
	var maxScore int
	scores := make(map[int]BlackJackPlayer)
	for _, player := range g.players {
		p := BlackJackPlayer(player)
		scores[p.Score] = p
		if maxScore < p.Score {
			maxScore = p.Score
		}
	}
	scores[g.dealer.Score] = BlackJackPlayer(g.dealer)
	return scores, maxScore
}

func (g *BlackJack) AddPlayer(player Player) {
	game := CardGame(*g)
	game.addPlayer(player)
	g.players = game.players
}

func (g *BlackJack) GetFirstPlayer() BlackJackPlayer {
	game := CardGame(*g)
	return BlackJackPlayer(game.getFirstPlayer())
}

func (g *BlackJack) AddDealer(player Player) {
	g.dealer = player
}

func (g *BlackJack) HasDealer() bool {
	return g.dealer.Id != 0 && g.dealer.Name != ""
}

func (g *BlackJack) GetCardScore(card Card) int {
	cardScore := card.GetValue()
	if cardScore >= Knight {
		return FaceCardValue
	} else {
		return int(cardScore)
	}
}

func (g *BlackJack) Reset() {
	g.dealer.cards = nil
	g.dealer.Score = 0
	for i := range g.players {
		g.players[i].cards = nil
		g.players[i].Score = 0
	}
}
