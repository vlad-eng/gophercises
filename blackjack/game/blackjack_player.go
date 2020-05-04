package game

import (
	"bufio"
	"fmt"
	. "gophercises/deck/deck"
	"os"
	"strconv"
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
	if p.PType == PlayerType {
		isDoubleDown := p.doubleDown()
		toHit := true
		var hitCard Card
		if isDoubleDown {
			if hitCard, err = p.dealCard(true); err != nil {
				return err
			}
			p.cards = Add(p.cards, hitCard)
			p.UpdateScore(hitCard)
			p.DisplayCards(true)
		} else {
			for toHit {
				if toHit, err = p.toHit(); err != nil {
					return err
				}
				p.DisplayCards(true)
			}
		}
	} else {
		var hitCard Card
		for p.Score < 17 {
			if hitCard, err = p.dealCard(true); err != nil {
				return err
			}
			p.cards = Add(p.cards, hitCard)
			p.UpdateScore(hitCard)
		}

		for i := range p.cards {
			p.cards[i].SetVisible(true)
		}

		p.DisplayCards(true)
	}
	return nil
}

func (p *BlackJackPlayer) doubleDown() bool {
	fmt.Println("Double down? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := scanner.Text()
	isDoubleDown := false
	if strings.Compare(answer, "y") == 0 {
		p.BetAmount *= 2
		isDoubleDown = true
	}
	return isDoubleDown
}

func (p *BlackJackPlayer) placeBet() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(p.Name + ": How much would you like to bet?")
	scanner.Scan()
	p.BetAmount, _ = strconv.Atoi(scanner.Text())
	p.BetAmount = 10
}

func (p *BlackJackPlayer) toHit() (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println()
	fmt.Print(p.Name + ": Do you want to hit? ")
	scanner.Scan()
	fmt.Println()
	var hitCard Card
	var err error
	answer := scanner.Text()
	toHit := strings.Compare(answer, "y")
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

func (p *BlackJackPlayer) DisplayCards(showSecondCard bool) {
	fmt.Printf(p.Name + ": ")
	for _, card := range p.cards[:len(p.cards)-1] {
		fmt.Printf("%s; ", card.String())
	}
	if showSecondCard {
		fmt.Printf("%s \n", p.cards[len(p.cards)-1].String())
	} else {
		fmt.Println("Face down card")
	}
}

func (p *BlackJackPlayer) String() string {
	player := Player(*p)
	return player.String()
}

func (p *BlackJackPlayer) ComputeScore() {
	game := BlackJack(*p.Game)
	for _, card := range p.cards {
		cardScore := game.GetCardScore(card)
		p.Score += cardScore
	}
	hasAce, _ := p.hasOneAce()
	if hasAce {
		p.Score += 10
	}
}

func (p *BlackJackPlayer) UpdateScore(card Card) {
	game := BlackJack(*p.Game)
	p.Score += game.GetCardScore(card)
	hasDealtAce, hasHitAce := p.hasOneAce()
	if hasDealtAce && p.Score > BlackJackMaxScore {
		p.Score -= 10
	}
	if hasHitAce && p.Score <= BlackJackMaxScore-10 {
		p.Score += 10
	}
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

//Returns a boolean tuple, if it has one ace for dealt cards
//and if it has one ace for hit cards
func (p *BlackJackPlayer) hasOneAce() (bool, bool) {
	dealtAcesCount := 0
	hitAcesCount := 0
	for i, card := range p.cards {
		value := card.GetValue()
		if i < 2 && value == 1 {
			dealtAcesCount++
		} else if value == 1 {
			hitAcesCount++
		}
	}

	hasDealtAce := false
	if dealtAcesCount == 1 {
		hasDealtAce = true
	}

	hasHitAce := false
	if hitAcesCount == 1 {
		hasHitAce = true
	}

	return hasDealtAce, hasHitAce
}

func (p *BlackJackPlayer) WinBankUpdate() {
	p.Bank += p.BetAmount
	p.BetAmount = 0
}

func (p *BlackJackPlayer) LossBankUpdate() {
	p.Bank -= p.BetAmount
	p.BetAmount = 0
}

func (p *BlackJackPlayer) DisplayAmount() {
	fmt.Printf("%s current amount: %d\n", p.Name, p.Bank)
}
