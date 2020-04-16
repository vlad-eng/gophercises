package game

import (
	"fmt"
	. "gophercises/deck/deck"
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
}

func (g *BlackJack) GetPlayerCount() int {
	return len(g.players)
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

func (g *BlackJack) EarlyWinner() (BlackJackPlayer, error) {
	scores, maxScore := g.computeScores()
	if maxScore == BlackJackMaxScore {
		return scores[maxScore], nil
	} else if g.dealer.Score == BlackJackMaxScore {
		return scores[g.dealer.Score], nil
	} else {
		return BlackJackPlayer{}, fmt.Errorf("no winner yet")
	}
}

func (g *BlackJack) EndOfTurnWinner() BlackJackPlayer {
	scores, maxScore := g.computeScores()
	if (maxScore < g.dealer.Score && (g.dealer.Score <= BlackJackMaxScore)) ||
		maxScore > BlackJackMaxScore {

		maxScore = g.dealer.Score
	}
	return scores[maxScore]
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

func (g *BlackJack) DisplayAllCards() {
	for _, player := range g.GetPlayers() {
		fmt.Printf("Player: %s\n", player.String())
		player.DisplayCards()
	}
	fmt.Println("-------------------------")
	dealer := g.GetDealer()
	fmt.Printf("Dealer: %s\n", dealer.String())
	dealer.DisplayCards()
	fmt.Println("-------------------------")
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
