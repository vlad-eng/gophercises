package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const StandardDeckLength = 52

type ByTypeAndValue []Card

func (t ByTypeAndValue) Len() int      { return len(t) }
func (t ByTypeAndValue) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ByTypeAndValue) Less(i, j int) bool {
	if t[i].cType == 0 {
		return false
	}
	if t[j].cType == 0 {
		return true
	}
	if t[i].cType < t[j].cType {
		return true
	}
	if t[i].cType > t[j].cType {
		return false
	}
	return t[i].value < t[j].value
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
			card := Card{value, cType, false}
			deck = append(deck, card)
		}
	}

	for _, option := range options {
		deck = option(deck)
	}
	return deck
}

func WithJokers(numOfJokers int) func([]Card) []Card {
	addJokers := func(d []Card) []Card {
		jokerCard := Card{value: Joker}
		for j := 0; j < numOfJokers; j++ {
			d = append(d, jokerCard)
		}
		return d
	}
	return addJokers
}

func WithoutCards(cardValues []CardValue) func([]Card) []Card {
	removeCards := func(d []Card) []Card {
		deck := make([]Card, 0)
		for _, checkedCard := range d {
			if Contains(cardValues, checkedCard.value) == false {
				deck = append(deck, checkedCard)
			}
		}
		return deck
	}
	return removeCards
}

func WithCards(cards []Card) func([]Card) []Card {
	addedCards := func(d []Card) []Card {
		for _, card := range cards {
			d = append(d, card)
		}
		return d
	}
	return addedCards
}

func WithShuffling() func([]Card) []Card {
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

func WithSorting(comparisonFuncs ...func(int, int) bool) func([]Card) []Card {
	var sortFunc func([]Card) []Card

	if len(comparisonFuncs) == 0 {
		sortFunc = func(deck []Card) []Card {
			sort.Sort(ByTypeAndValue(deck))
			return deck
		}
		return sortFunc

	}

	sortFunc = func(deck []Card) []Card {
		comparatorFunc := comparisonFuncs[0]
		sort.Slice(deck, comparatorFunc)
		return deck
	}

	return sortFunc
}

func WithMultipleStandardDecks(deckCount int) func([]Card) []Card {
	multipleDecksFunc := func(standardDeck []Card) []Card {
		multipleDecks := make([]Card, 0)
		for i := 1; i <= deckCount; i++ {
			multipleDecks = append(multipleDecks, standardDeck...)
		}
		return multipleDecks
	}
	return multipleDecksFunc
}

func Create(firstValue, lastValue CardValue, types []CardType) []Card {
	deck := make([]Card, 0)
	for _, cType := range types {
		for value := firstValue; value <= lastValue; value++ {
			deck = append(deck, Card{cType: cType, value: value})
		}
	}
	return deck
}

func Add(deck []Card, card Card) []Card {
	deck = append(deck, card)
	return deck
}

func RemoveLast(deck []Card) ([]Card, Card, error) {
	if len(deck) > 0 {
		deck = deck[:len(deck)-1]
		removedCard := deck[len(deck)-1]
		return deck, removedCard, nil
	}

	return deck, Card{}, fmt.Errorf("no cards left in deck")
}

func Contains(s []CardValue, e CardValue) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
