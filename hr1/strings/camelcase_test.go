package strings

import (
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CamelCaseTestSuite struct {
	suite.Suite
	unit   StringTokenizer
	gomega *GomegaWithT
}

func Test_CamelCaseTestSuite(t *testing.T) {
	gomega := NewGomegaWithT(t)
	testSuite := CamelCaseTestSuite{unit: StringTokenizer{}, gomega: gomega}
	suite.Run(t, &testSuite)
}

func (s *CamelCaseTestSuite) Test_EmptyStringHasNoWords() {
	camelcase := ""
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(0))
}

func (s *CamelCaseTestSuite) Test_OneSimpleWordReturnsOne() {
	camelcase := "camel"
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(1))
}

func (s *CamelCaseTestSuite) Test_TwoSimpleWordsReturnTwo() {
	camelcase := "camelCase"
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(2))
}

func (s *CamelCaseTestSuite) Test_TwoSimpleNonCamelCaseWordsLeadToPanic() {
	s.gomega.Expect(panicCase).To(Panic())
}

func panicCase() {
	st := StringTokenizer{}
	st.Tokenize("CamelCase")
}

func (s *CamelCaseTestSuite) Test_FiveSimpleWordsReturnFive() {
	camelcase := "  saveChangesInTheEditor    "
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(5))
}

func (s *CamelCaseTestSuite) Test_PhraseWithThreeCamelCaseWordsReturnsSumOfSimpleWordCount() {
	camelcase := "camelCase    wasIs a      wideSpreadConvention"
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(8))
}

func (s *CamelCaseTestSuite) Test_PhraseWithNonCamelCaseWordReturnsSumOfSimpleWordsFromCamelCaseWords() {
	camelcase := "camelCase    wasIs a      wideSpreadConvention ,,,	RightAmI?"
	st := StringTokenizer{}
	s.gomega.Expect(st.Tokenize(camelcase)).Should(Equal(8))
}
