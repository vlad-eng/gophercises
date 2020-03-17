package sitemap

import (
	"encoding/xml"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
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
	testSuite := SiteParserTestSuite{gomega: gomega}
	suite.Run(t, &testSuite)
}

func (s *SiteParserTestSuite) Test_ExtractsOnlyDomainLinks() {
	domain := "http://www.wikipedia.org"
	s.unit = InitSiteParser(domain)
	links := s.unit.Parse()
	s.gomega.Expect(len(links)).ShouldNot(Equal(0))
}

func (s *SiteParserTestSuite) Test_ExtractLinksInChildPages() {
	domain := "http://www.thesmallthingsblog.com"
	s.unit = InitSiteParser(domain)
	links := s.unit.Parse()
	found := false
	pageUrl := "http://www.thesmallthingsblog.com/shop/the-master-list/"
	for _, link := range links {
		if strings.Contains(link, pageUrl) {
			found = true
		}
	}
	s.gomega.Expect(found).Should(Equal(true))
}

func (s *SiteParserTestSuite) Test_FormatProducesExpectedXml() {
	domain := "http://www.thesmallthingsblog.com"
	s.unit = InitSiteParser(domain)
	links := s.unit.Parse()

	expectedSiteMap := UrlSet{
		Url: []Url{
			{Loc: "http://www.thesmallthingsblog.com"},
			{Loc: "http://www.thesmallthingsblog.com/category/beauty/hair/"},
			{Loc: "http://www.thesmallthingsblog.com/category/hair-tutorial/"},
			{Loc: "http://www.thesmallthingsblog.com/category/hair-products/"},
		},
	}

	marshalledExpectedSiteMap, _ := xml.MarshalIndent(expectedSiteMap, "", "   ")
	expectedResult := xml.Header + string(marshalledExpectedSiteMap)

	actualResult := XmlEncode(links[:4])
	s.gomega.Expect(actualResult).Should(Equal(expectedResult))
}
