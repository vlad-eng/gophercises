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
	fmt.Println("first test")
	game := BlackJack{}
	dealer := Player{
		id:    1234,
		name:  "Dealer_1",
		cards: make([]Card, 0),
	}

	game.addDealer(dealer)
	player := Player{
		cards: make([]Card, 0),
	}
	game.addPlayer(player)
	game.dealCards()
	firstBlackJackPlayer := game.getFirstPlayer()
	s.gomega.Expect(firstBlackJackPlayer.isDealt()).Should(BeTrue())
}
