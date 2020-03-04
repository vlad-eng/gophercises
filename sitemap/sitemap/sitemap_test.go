package sitemap

import (
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	. "gophercises/link/parser"
	"strings"
	"testing"
)

type SiteParserTestSuite struct {
	suite.Suite
	unit   SiteParser
	gomega *GomegaWithT
}

func Test_SiteParserTests(t *testing.T) {
	gomega := NewGomegaWithT(t)
	testSuite := SiteParserTestSuite{unit: SiteParser{}, gomega: gomega}
	suite.Run(t, &testSuite)
}

func (s *SiteParserTestSuite) Test_ExtractsOnlyDomainLinks() {
	domain := "http://www.wikipedia.org"
	var links []Link
	var err error
	if links, err = s.unit.Parse(domain); err != nil {
		panic(err)
	}
	s.gomega.Expect(len(links)).ShouldNot(Equal(0))
}

func (s *SiteParserTestSuite) Test_ExtractLinksInChildPages() {
	domain := "http://www.wikipedia.org"
	var links []Link
	var err error
	if links, err = s.unit.Parse(domain); err != nil {
		panic(err)
	}
	s.gomega.Expect(len(links)).ShouldNot(Equal(0))
	found := false
	expectedDomain := "wikipedia"
	for _, link := range links {
		if strings.Contains(link.Href, expectedDomain) == true {
			found = true
		}
	}
	s.gomega.Expect(found).Should(Equal(true))
}
