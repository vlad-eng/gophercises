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
	deck := New(withJokers(1), withoutCards(CardsOfThree))
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
		CardsOfTwo = append(CardsOfTwo, Card{Two, cType})
	}
	var CardsOfThree []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfThree = append(CardsOfThree, Card{Three, cType})
	}
	s.gomega.Expect(s.unit).Should(ContainElements(CardsOfTwo))
	s.gomega.Expect(s.unit).ShouldNot(ContainElements(CardsOfThree))
}

func (s *CardDeckTestSuite) Test_DeckIsShuffled() {
	var CardsOfThree []Card
	for cType := Spades; cType <= Hearts; cType++ {
		CardsOfThree = append(CardsOfThree, Card{Three, cType})
	}
	unshuffledDeck := New(withJokers(1), withoutCards(CardsOfThree))
	s.gomega.Expect(s.unit).ShouldNot(Equal(unshuffledDeck))
}
