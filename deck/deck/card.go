package deck

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type CardValue int
type CardType int

const (
	Joker CardValue = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Knight
	Queen
	King
)

func (v *CardValue) String() string {
	values := []string{
		"Joker",
		"Ace",
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
		"Knight",
		"Queen",
		"King",
	}

	return values[*v]
}

const (
	NoType CardType = iota
	Spades
	Diamonds
	Clubs
	Hearts
)

func (t *CardType) String() string {
	types := []string{
		"",
		"Spades",
		"Diamonds",
		"Clubs",
		"Hearts",
	}

	return types[*t]
}

type Card struct {
	value CardValue
	cType CardType
}

func (c *Card) String() string {
	name := fmt.Sprintf("%sOf%s", c.value.String(), c.cType.String())
	name = strings.TrimSuffix(name, "Of")
	return name
}

func New(options ...func([]Card) []Card) []Card {
	deck := make([]Card, 0)
	for cType := Spades; cType <= Hearts; cType++ {
		for value := Ace; value <= King; value++ {
			card := Card{value, cType}
			deck = append(deck, card)
		}
	}

	for _, option := range options {
		deck = option(deck)
	}
	return deck
}

func withJokers(numOfJokers int) func([]Card) []Card {
	addJokers := func(d []Card) []Card {
		jokerCard := Card{value: Joker}
		for j := 0; j < numOfJokers; j++ {
			d = append(d, jokerCard)
		}
		return d
	}
	return addJokers
}

func withoutCards(cards []Card) func([]Card) []Card {
	removeCards := func(d []Card) []Card {
		for i, card := range cards {
			val := (i * 13) + (int(card.value) - 1)
			d = append(d[:val], d[val+1:]...)
		}
		return d
	}
	return removeCards
}

func withCards(cards []Card) func([]Card) []Card {
	addedCards := func(d []Card) []Card {
		for _, card := range cards {
			d = append(d, card)
		}
		return d
	}
	return addedCards
}

func withShuffling() func([]Card) []Card {
	shuffleFunc := func(deck []Card) []Card {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for n := len(deck); n > 1; n-- {
			randIndex := r.Intn(n)
			deck[randIndex], deck[n-1] = deck[n-1], deck[randIndex]
		}
		return deck
	}
	return shuffleFunc
}
