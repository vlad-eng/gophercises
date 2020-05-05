package game

import (
	"bufio"
	"fmt"
	. "gophercises/deck/deck"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
	if p.PType == Human || p.PType == AI {
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

func (p *BlackJackPlayer) userInput(question string) string {
	fmt.Println()
	fmt.Println(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (p *BlackJackPlayer) aiInput(question string) string {
	if strings.Contains(question, "Double down? ") {
		if p.Score >= 10 && p.Score <= 11 {
			return p.probabilisticAnswer(0.9)
		} else if p.Score == 9 {
			return p.probabilisticAnswer(0.4)
		} else {
			return "n"
		}
	} else if strings.Contains(question, ": How much would you like to bet?") {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		return string(r.Intn(MaxPossibleBetAmount))
	} else if strings.Contains(question, ": Do you want to hit? ") {
		if p.Score <= 11 {
			return "y"
		} else if p.Score >= 12 && p.Score <= 16 {
			dealerShownValue := p.Game.dealer.cards[0].GetValue()
			if dealerShownValue <= 6 {
				return "n"
			} else {
				return "y"
			}
		} else {
			return "n"
		}
	}
	panic("Unknown question")
}

func (p *BlackJackPlayer) probabilisticAnswer(probabilityForYes float32) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randValue := r.Float32()
	if randValue <= probabilityForYes {
		return "y"
	} else {
		return "n"
	}
}

func (p *BlackJackPlayer) doubleDown() bool {
	var answer string
	if p.PType == Human {
		answer = p.userInput("Double down? ")
	} else {
		answer = p.aiInput("Double down? ")
	}
	isDoubleDown := false
	if strings.Compare(answer, "y") == 0 {
		p.BetAmount *= 2
		isDoubleDown = true
	}
	return isDoubleDown
}

func (p *BlackJackPlayer) placeBet() {
	var answer string
	if p.PType == Human {
		answer = p.userInput(p.Name + ": How much would you like to bet?")
	} else {
		answer = p.aiInput(p.Name + ": How much would you like to bet?")
	}

	p.BetAmount, _ = strconv.Atoi(answer)
	p.BetAmount = 10
}

func (p *BlackJackPlayer) toHit() (bool, error) {
	var answer string
	if p.PType == Human {
		answer = p.userInput(p.Name + ": Do you want to hit? ")
	} else {
		answer = p.aiInput(p.Name + ": Do you want to hit? ")
	}

	var hitCard Card
	var err error
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
	if p.PType == Human || p.PType == AI {
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
