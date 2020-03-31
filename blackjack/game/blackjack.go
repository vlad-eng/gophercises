package game

import . "gophercises/deck/deck"

type Player struct {
	id    int
	name  string
	cards []Card
}

type BlackJackPlayer Player

type CardGame struct {
	dealer  *Player
	players []Player
}

type BlackJack CardGame

func (p *BlackJackPlayer) isDealt() bool {
	return false
}

func NewGame() *BlackJack {
	game := BlackJack{}
	return &game
}

func (g *BlackJack) getFirstPlayer() BlackJackPlayer {
	return BlackJackPlayer(g.players[0])
}

func (g *BlackJack) addPlayer(player Player) {
	g.players = append(g.players, player)
}

func (g *BlackJack) addDealer(player Player) {
	g.dealer = &player
}

func (g *BlackJack) hasDealer() bool {
	return g.dealer != nil
}

func (g *BlackJack) dealCards() {
	//TODO
}
