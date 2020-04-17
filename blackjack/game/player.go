package game

import (
	. "gophercises/deck/deck"
)

type ParticipantType int

const (
	DealerType ParticipantType = iota
	PlayerType
)

type Player struct {
	Id        int
	Name      string
	cards     []Card
	PType     ParticipantType
	Game      *CardGame
	Score     int
	Bank      int
	BetAmount int
}

func (p *Player) String() string {
	return p.Name
}
func (p *Player) SetScore(score int) {
	p.Score = score
}

func (g *CardGame) getFirstPlayer() Player {
	return g.players[0]
}

func (g *CardGame) addPlayer(player Player) {
	g.players = append(g.players, player)
}
