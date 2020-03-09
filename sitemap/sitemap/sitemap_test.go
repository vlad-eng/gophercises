package sitemap

import (
	"encoding/xml"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	. "gophercises/link/parser"
	"strings"
	"testing"
	"time"
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
	links := make(chan Link, 1000)
	doneChannel := make(chan bool, 1)
	go s.unit.Parse(domain, links, doneChannel)

	for {
		select {
		case <-doneChannel:
			s.gomega.Expect(len(links)).ShouldNot(Equal(0))
			return
		case <-time.After(1000 * time.Millisecond):
			s.gomega.Expect(len(links)).ShouldNot(Equal(0))
			return
		}
	}
}

func (s *SiteParserTestSuite) Test_ExtractLinksInChildPages() {
	domain := "http://www.thesmallthingsblog.com"
	visitedLinks := make(chan Link, 1000)
	doneChannel := make(chan bool, 1)
	go s.unit.Parse(domain, visitedLinks, doneChannel)
	found := false
	expectedDomain := "http://www.thesmallthingsblog.com/shop/the-master-list/"
	links := make([]Link, 1000)

	for {
		select {
		case link := <-visitedLinks:
			links = append(links, link)
			if strings.Contains(link.Href, expectedDomain) {
				found = true
			}
		case <-doneChannel:
			s.gomega.Expect(found).Should(Equal(true))
			return
		case <-time.After(1000 * time.Millisecond):
			s.gomega.Expect(found).Should(Equal(true))
			return
		}
	}
}

func (s *SiteParserTestSuite) Test_FormatProducesExpectedXml() {
	domain := "http://www.thesmallthingsblog.com"
	visitedLinks := make(chan Link, 1000)
	isFinished := make(chan bool, 1)
	var links []Link
	go s.unit.Parse(domain, visitedLinks, isFinished)

	expectedSiteMap := UrlSet{
		Url: []Url{
			{Loc: "http://www.thesmallthingsblog.com/category/beauty/hair/"},
			{Loc: "http://www.thesmallthingsblog.com/category/hair-tutorial/"},
			{Loc: "http://www.thesmallthingsblog.com/category/hair-products/"},
			{Loc: "http://www.thesmallthingsblog.com/category/ask-kate/"},
		},
	}

	marshalledExpectedSiteMap, _ := xml.MarshalIndent(expectedSiteMap, "", "   ")
	expectedResult := xml.Header + string(marshalledExpectedSiteMap)

	for {
		select {
		case link := <-visitedLinks:
			links = append(links, link)
		case <-isFinished:
			actualResult := s.unit.Format(links[:4])
			s.gomega.Expect(actualResult).Should(Equal(expectedResult))
			return
		case <-time.After(300 * time.Millisecond):
			actualResult := s.unit.Format(links[:4])
			s.gomega.Expect(actualResult).Should(Equal(expectedResult))
			return
		}
	}
}
