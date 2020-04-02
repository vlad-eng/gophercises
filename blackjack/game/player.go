package game

import . "gophercises/deck/deck"

type ParticipantType int

const (
	DealerType ParticipantType = iota
	PlayerType
)

type Player struct {
	id          int
	name        string
	cards       []Card
	pType       ParticipantType
	currentGame *CardGame
}

func (g *CardGame) getFirstPlayer() Player {
	return g.players[0]
}

func (g *CardGame) addPlayer(player Player) {
	g.players = append(g.players, player)
}
