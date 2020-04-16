package game

import (
	"bufio"
	"fmt"
	. "gophercises/deck/deck"
	"os"
	"strings"
)

type BlackJackPlayer Player

func (p *BlackJackPlayer) isDealt() bool {
	return p.cards != nil && len(p.cards) >= 2
}

func (p *BlackJackPlayer) dealCards() error {
	var dealtCard Card
	var err error
	cardsVisibility := p.getCardVisibility()

	for i := 0; i < DealtCardsCount; i++ {
		if dealtCard, err = p.dealCard(cardsVisibility[i]); err != nil {
			return err
		}
		p.cards = Add(p.cards, dealtCard)
	}

	return nil
}

func (p *BlackJackPlayer) dealCard(isVisible bool) (Card, error) {
	deck := BlackJack(*p.Game).deck
	var dealtCard Card
	var err error
	if p.Game.deck, dealtCard, err = RemoveLast(deck); err != nil {
		return Card{}, err
	}
	dealtCard.SetVisible(isVisible)
	return dealtCard, nil
}

func (p *BlackJackPlayer) ExecuteTurn() (err error) {
	toHit := true
	if p.PType == PlayerType {
		for toHit {
			if toHit, err = p.toHit(); err != nil {
				return err
			}
			p.DisplayCards()
		}
	} else {
		//TODO: dealer's turn
		//first version below:
		p.cards[1].SetVisible(true)
		p.DisplayCards()
	}
	return nil
}

func (p *BlackJackPlayer) toHit() (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println()
	fmt.Print(p.Name + ": Hit or Stand: ")
	scanner.Scan()
	fmt.Println()
	var hitCard Card
	var err error
	answer := scanner.Text()
	//answer := "H"
	toHit := strings.Compare(answer, "H")
	if toHit == 0 {
		if hitCard, err = p.dealCard(true); err != nil {
			return false, err
		}
		p.cards = Add(p.cards, hitCard)
		p.UpdateScore(hitCard)
		return true, nil
	}
	return false, nil
}

func (p *BlackJackPlayer) DisplayCards() {
	fmt.Printf(p.Name + ": ")
	for _, card := range p.cards[:len(p.cards)-1] {
		fmt.Printf("%s; ", card.String())
	}
	fmt.Printf("%s \n", p.cards[len(p.cards)-1].String())
}

func (p *BlackJackPlayer) String() string {
	player := Player(*p)
	return player.String()
}

func (p *BlackJackPlayer) ComputeScore() {
	var score int
	game := BlackJack(*p.Game)
	for _, card := range p.cards {
		cardScore := game.GetCardScore(card)
		score += cardScore
	}
	p.Score = score
}

func (p *BlackJackPlayer) UpdateScore(card Card) {
	game := BlackJack(*p.Game)
	p.Score += game.GetCardScore(card)
}

func (p *BlackJackPlayer) getCardVisibility() []bool {
	var cardsVisibility []bool
	if p.PType == PlayerType {
		cardsVisibility = []bool{true, true}
	} else {
		cardsVisibility = []bool{true, false}
	}
	return cardsVisibility
}
