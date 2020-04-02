package game

import (
	. "gophercises/deck/deck"
)

const DealtCardsCount = 2

type CardGame struct {
	dealer  Player
	players []Player
	deck    []Card
}

type BlackJack CardGame

func NewGame() *BlackJack {
	game := BlackJack{
		deck: New(WithShuffling()),
	}
	return &game
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

	dealerPlayer := BlackJackPlayer(g.dealer)
	if dealerPlayer.dealCards(); err != nil {
		return err
	}

	return nil
}

func (g *BlackJack) addPlayer(player Player) {
	game := CardGame(*g)
	game.addPlayer(player)
	g.players = game.players
}

func (g *BlackJack) getFirstPlayer() BlackJackPlayer {
	game := CardGame(*g)
	return BlackJackPlayer(game.getFirstPlayer())
}

func (g *BlackJack) addDealer(player Player) {
	g.dealer = player
}

func (g *BlackJack) hasDealer() bool {
	return g.dealer.id != 0 && g.dealer.name != ""
}
