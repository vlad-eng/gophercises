package adventure

import (
	"github.com/golang/go/src/pkg/net/http/httptest"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"log"
	. "net/http"
	"testing"
)

type ChooseAdventureTestSuite struct {
	suite.Suite
	gomega *GomegaWithT
}

func TestChooseAdventure(t *testing.T) {
	gomega := NewGomegaWithT(t)
	testSuite := ChooseAdventureTestSuite{gomega: gomega}
	suite.Run(t, &testSuite)
}

func (s *ChooseAdventureTestSuite) Test_AdventureStep() {
	expectedNextPage := Page{
		Title: "The Great Debate",
		Story: []string{"After a bit everyone settles down the two people on stage begin having a debate. You don't recall too many specifics, but for some reason you have a feeling you are supposed to pick sides."},
		Options: []Option{
			{
				Text: "Clearly that man in the fox outfit was the winner.",
				Arc:  "sean-kelly",
			},
			{
				Text: "I don't think those fake abs would help much in a feat of strength, but our caped friend clearly won this bout. Let's go congratulate him.",
				Arc:  "mark-bates",
			},
			{
				Text: "Slip out the back before anyone asks us to pick a side.",
				Arc:  "home",
			},
		},
	}

	var testStory *Story
	var err error
	if testStory, err = Decode("gopher.json"); err != nil {
		panic(err)
	}
	var nextPage *Page
	if nextPage, err = testStory.NextPage("new-york", "debate"); err != nil {
		panic(err)
	}
	s.gomega.Expect(*nextPage).Should(Equal(expectedNextPage))
}

func (s *ChooseAdventureTestSuite) Test_PageHandlerRetrievesCorrectPage() {
	var testStory *Story
	var err error
	if testStory, err = Decode("gopher.json"); err != nil {
		panic(err)
	}
	storyServer := StoryServer{PreviousRequest: "/new-york", Story: testStory}
	writer := httptest.NewRecorder()
	request := createRequest("/debate")
	storyServer.ServeHTTP(writer, request)

	s.gomega.Expect(writer.Code).Should(Equal(StatusOK))
}

func (s *ChooseAdventureTestSuite) Test_InternalServerErrorWhenIncorrectRequest() {
	var testStory *Story
	var err error
	if testStory, err = Decode("gopher.json"); err != nil {
		panic(err)
	}
	storyServer := StoryServer{PreviousRequest: "/new-york", Story: testStory}
	writer := httptest.NewRecorder()
	request := createRequest("mark-kelly")
	storyServer.ServeHTTP(writer, request)
	s.gomega.Expect(writer.Code).Should(Equal(StatusInternalServerError))
}

func createRequest(path string) *Request {
	var req *Request
	var err error
	req, err = NewRequest("GET", "localhost:8088", nil)
	req.URL.Path = path
	if err != nil {
		log.Fatal("couldn't create request due to: ", err)
	}

	return req
}
