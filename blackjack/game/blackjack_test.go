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
		id:          1234,
		name:        "Dealer_1",
		cards:       []Card{},
		pType:       DealerType,
		currentGame: &blackJackGame,
	}

	s.unit.addDealer(dealer)
	player := Player{
		id:          5678,
		name:        "Player_1",
		cards:       []Card{},
		pType:       PlayerType,
		currentGame: &blackJackGame,
	}

	s.unit.addPlayer(player)
	s.unit.dealCards()
	firstBlackJackPlayer := s.unit.getFirstPlayer()

	s.gomega.Expect(firstBlackJackPlayer.isDealt()).Should(BeTrue())
}
