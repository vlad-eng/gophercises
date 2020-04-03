package game

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	. "gophercises/deck/deck"
	"testing"
)

type BlackjackTestSuite struct {
	suite.Suite
	unit   *BlackJack
	gomega *GomegaWithT
}

func Test_BlackjackTestSuite(t *testing.T) {
	testSuite := BlackjackTestSuite{unit: NewGame(), gomega: NewGomegaWithT(t)}
	suite.Run(t, &testSuite)
	fmt.Println("setup")
}

func (s *BlackjackTestSuite) Test_FirstPlayerIsDealt() {
	blackJackGame := CardGame(*s.unit)
	dealer := Player{
		Id:    1234,
		Name:  "Dealer_1",
		cards: []Card{},
		PType: DealerType,
		Game:  &blackJackGame,
	}

	s.unit.AddDealer(dealer)
	player := Player{
		Id:    5678,
		Name:  "Player_1",
		cards: []Card{},
		PType: PlayerType,
		Game:  &blackJackGame,
	}

	s.unit.AddPlayer(player)
	s.unit.DealCards()
	firstBlackJackPlayer := s.unit.GetFirstPlayer()

	s.gomega.Expect(firstBlackJackPlayer.isDealt()).Should(BeTrue())
}
