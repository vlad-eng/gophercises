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
	var CardsOfThree []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfThree = append(CardsOfThree, Card{Three, cType})
	}
	deck := New(withJokers(1))
	testSuite := CardDeckTestSuite{unit: deck, gomega: NewGomegaWithT(t)}
	suite.Run(t, &testSuite)
}

func (s *CardDeckTestSuite) Test_NamesOfCardsCorrectlyReturned() {
	cardNames := make([]string, 0)
	for i := 0; i < len(s.unit); i++ {
		cardNames = append(cardNames, s.unit[i].String())
	}
	s.gomega.Expect(cardNames).Should(ContainElement("AceOfSpades"))
	s.gomega.Expect(cardNames).Should(ContainElement("TenOfHearts"))
	s.gomega.Expect(cardNames).Should(ContainElement("Joker"))

	//s.gomega.Expect(s.unit[48].String()).Should(ContainElement("TenOfHearts"))
	//s.gomega.Expect(s.unit[52].String()).Should(ContainElement("Joker"))
}

func (s *CardDeckTestSuite) Test_AllCardsWithValueOfThreeAreRemoved() {
	var CardsOfTwo []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfTwo = append(CardsOfTwo, Card{Two, cType})
	}
	var CardsOfThree []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfThree = append(CardsOfThree, Card{Three, cType})
	}
	withoutCardsFunc := withoutCards(CardsOfThree)
	s.unit = withoutCardsFunc(s.unit)
	s.gomega.Expect(s.unit).Should(ContainElements(CardsOfTwo))
	s.gomega.Expect(s.unit).ShouldNot(ContainElements(CardsOfThree))

	withCardsFunc := withCards(CardsOfThree)
	s.unit = withCardsFunc(s.unit)
}

func (s *CardDeckTestSuite) Test_DeckIsShuffled() {
	shuffledDeck := New(withJokers(1), withShuffling())
	s.gomega.Expect(s.unit).ShouldNot(Equal(shuffledDeck))
	s.gomega.Expect(len(s.unit)).Should(Equal(len(shuffledDeck)))
}
