package game

import . "gophercises/deck/deck"

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
	deck := BlackJack(*p.currentGame).deck
	var dealtCard Card
	var err error
	if p.currentGame.deck, dealtCard, err = RemoveLast(deck); err != nil {
		return Card{}, err
	}
	dealtCard.SetVisible(isVisible)
	return dealtCard, nil
}

func (p *BlackJackPlayer) getCardVisibility() []bool {
	var cardsVisibility []bool
	if p.pType == PlayerType {
		cardsVisibility = []bool{true, true}
	} else {
		cardsVisibility = []bool{true, false}
	}
	return cardsVisibility
}
