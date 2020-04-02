package deck

import (
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CardDeckTestSuite struct {
	suite.Suite
	unit   []Card
	gomega *GomegaWithT
}

func Test_CardSuitTestSuite(t *testing.T) {
	deck := New(WithJokers(1))
	testSuite := CardDeckTestSuite{unit: deck, gomega: NewGomegaWithT(t)}
	suite.Run(t, &testSuite)
}

func (s *CardDeckTestSuite) Test_NamesOfCardsCorrectlyReturned() {
	s.gomega.Expect(s.unit[0].String()).Should(Equal("AceOfSpades"))
	s.gomega.Expect(s.unit[48].String()).Should(Equal("TenOfHearts"))
	s.gomega.Expect(s.unit[52].String()).Should(Equal("Joker"))
}

func (s *CardDeckTestSuite) Test_AllCardsWithValueOfThreeAreRemoved() {
	var CardsOfTwo []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfTwo = append(CardsOfTwo, Card{Two, cType, false})
	}
	var CardsOfThree []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfThree = append(CardsOfThree, Card{Three, cType, false})
	}
	withoutCardsFunc := WithoutCards([]CardValue{Three})
	s.unit = withoutCardsFunc(s.unit)
	s.gomega.Expect(s.unit).Should(ContainElements(CardsOfTwo))
	s.gomega.Expect(s.unit).ShouldNot(ContainElements(CardsOfThree))

	withCardsFunc := WithCards(CardsOfThree)
	s.unit = withCardsFunc(s.unit)
	withSortFunc := WithSorting()
	s.unit = withSortFunc(s.unit)
}

func (s *CardDeckTestSuite) Test_DeckIsShuffled() {
	shuffledDeck := New(WithJokers(1), WithShuffling())
	s.gomega.Expect(s.unit).ShouldNot(Equal(shuffledDeck))
	s.gomega.Expect(len(s.unit)).Should(Equal(len(shuffledDeck)))
}

func (s *CardDeckTestSuite) Test_DeckIsSortedAfterShuffle() {
	shuffledDeck := New(WithJokers(1), WithShuffling())
	s.gomega.Expect(s.unit).ShouldNot(Equal(shuffledDeck))
	sortFunc := WithSorting()
	sortedDeck := sortFunc(shuffledDeck)
	s.gomega.Expect(s.unit).Should(Equal(sortedDeck))
}

func (s *CardDeckTestSuite) Test_DeckIsSortedWithCustomFunction() {
	cardTypes := []CardType{
		Hearts,
		Clubs,
		Diamonds,
		Spades,
	}
	sortedDeck := Create(Nine, Knight, cardTypes)
	values := getAllCardValues()
	removedValues := make([]CardValue, 0)
	removedValues = append(removedValues, values[0:8]...)
	removedValues = append(removedValues, values[11:]...)
	removedValues = append(removedValues, Joker)
	removalFunc := WithoutCards(removedValues)
	s.unit = removalFunc(s.unit)
	less := func(i, j int) bool {
		if s.unit[i].cType == 0 {
			return false
		}
		if s.unit[j].cType == 0 {
			return true
		}
		if s.unit[i].cType > s.unit[j].cType {
			return true
		}
		if s.unit[i].cType < s.unit[j].cType {
			return false
		}
		return s.unit[i].value < s.unit[j].value
	}
	sortFunc := WithSorting(less)
	s.unit = sortFunc(s.unit)
	s.gomega.Expect(s.unit).Should(Equal(sortedDeck))
	s.unit = New(WithJokers(1))
}

func (s *CardDeckTestSuite) Test_ConstructDeckWithThreeStandardDecks() {
	s.unit = New(WithMultipleStandardDecks(3))
	s.gomega.Expect(len(s.unit)).Should(Equal(StandardDeckLength * 3))
	s.gomega.Expect(s.unit[0:StandardDeckLength]).Should(Equal(s.unit[StandardDeckLength : StandardDeckLength*2]))
	s.gomega.Expect(s.unit[0:StandardDeckLength]).Should(Equal(s.unit[StandardDeckLength*2 : StandardDeckLength*3]))
	s.unit = New(WithJokers(1))
}
