package game

import (
	"fmt"
	. "gophercises/deck/deck"
)

const DealtCardsCount = 2

type Player struct {
	id          int
	name        string
	cards       []Card
	currentGame *CardGame
}

type BlackJackPlayer Player

type CardGame struct {
	dealer  *Player
	players []Player
	deck    []Card
}

type BlackJack CardGame

func (p *BlackJackPlayer) isDealt() bool {
	return p.cards != nil && len(p.cards) >= 2
}

func (g *BlackJack) dealCards() error {
	var err error
	for i, player := range g.players {
		p := BlackJackPlayer(player)

		if err = p.dealCards(); err != nil {
			return err
		}
		g.players[i] = Player(p)
	}

	//TODO: add dealing for dealer at the end of dealing stage
	//dealer := BlackJackPlayer(*g.dealer)
	//if err = dealer.dealCards(); err != nil {
	//	return err
	//}

	return nil
}

func (p *BlackJackPlayer) dealCards() error {
	var dealtCard Card
	var err error
	for i := 0; i < DealtCardsCount; i++ {
		if dealtCard, err = p.dealCard(); err != nil {
			return err
		}
		p.cards = Add(p.cards, dealtCard)
	}
	return nil
}

func (p *BlackJackPlayer) dealCard() (Card, error) {
	deck := BlackJack(*p.currentGame).deck
	var dealtCard Card
	var err error
	if p.currentGame.deck, dealtCard, err = RemoveLast(deck); err != nil {
		return Card{}, err
	}
	return dealtCard, nil

	return Card{}, fmt.Errorf("unknown game type")
}

func NewGame() *BlackJack {
	game := BlackJack{
		deck: New(WithShuffling()),
	}
	return &game
}

func (g *BlackJack) getFirstPlayer() BlackJackPlayer {
	return BlackJackPlayer(g.players[0])
}

func (g *BlackJack) addPlayer(player Player) {
	cardGame := CardGame(*g)
	player.currentGame = &cardGame
	g.players = append(g.players, player)
}

func (g *BlackJack) addDealer(player Player) {
	g.dealer = &player
}

func (g *BlackJack) hasDealer() bool {
	return g.dealer != nil
}
