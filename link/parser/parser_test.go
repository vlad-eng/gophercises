package parser

import (
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type LinkParserTestSuite struct {
	suite.Suite
	unit   PageParser
	gomega *GomegaWithT
}

func Test_LinkParserTests(t *testing.T) {
	gomega := NewGomegaWithT(t)
	testSuite := LinkParserTestSuite{unit: PageParser{}, gomega: gomega}
	suite.Run(t, &testSuite)
}

func (s *LinkParserTestSuite) Test_TheOnlyLinkCorrectlyExtracted() {
	var links []Link
	var err error
	htmlPath, _ := filepath.Abs("../ex1.html")
	var htmlBytes []byte
	if htmlBytes, err = ioutil.ReadFile(htmlPath); err != nil {
		panic("couldn't read input file" + err.Error())
	}
	if links, err = s.unit.Parse(string(htmlBytes)); err != nil {
		panic(err)
	}

	s.gomega.Expect(len(links)).Should(Equal(1))
	s.gomega.Expect(links[0].Href).Should(Equal("/other-page"))
	s.gomega.Expect(links[0].Text).Should(Equal("A link to another page"))
}

func (s *LinkParserTestSuite) Test_AllLinksAreExtracted() {
	var links []Link
	var err error
	htmlPath, _ := filepath.Abs("../ex2.html")
	var htmlBytes []byte
	if htmlBytes, err = ioutil.ReadFile(htmlPath); err != nil {
		panic("couldn't read input file" + err.Error())
	}
	if links, err = s.unit.Parse(string(htmlBytes)); err != nil {
		panic(err)
	}

	s.gomega.Expect(len(links)).Should(Equal(2))
	s.gomega.Expect(links[0].Href).Should(Equal("https://www.twitter.com/joncalhoun"))
	s.gomega.Expect(links[0].Text).Should(Equal("Check me out on twitter"))

	s.gomega.Expect(links[1].Href).Should(Equal("https://github.com/gophercises"))
	s.gomega.Expect(links[1].Text).Should(Equal("Gophercises is on Github!"))
}

func (s *LinkParserTestSuite) Test_Exercise3ParsedSuccessfully() {
	var links []Link
	var err error
	htmlPath, _ := filepath.Abs("../ex3.html")
	var htmlBytes []byte
	if htmlBytes, err = ioutil.ReadFile(htmlPath); err != nil {
		panic("couldn't read input file" + err.Error())
	}
	if links, err = s.unit.Parse(string(htmlBytes)); err != nil {
		panic(err)
	}

	testHref := []string{"#", "/lost", "https://twitter.com/marcusolsson"}
	testText := []string{"Login", "Lost? Need help?", "@marcusolsson"}

	s.gomega.Expect(len(links)).Should(Equal(3))
	for i, link := range links {
		s.gomega.Expect(link.Href).Should(Equal(testHref[i]))
		s.gomega.Expect(link.Text).Should(Equal(testText[i]))
	}
}

func (s *LinkParserTestSuite) Test_CorrectlyParsesLinksWithComments() {
	var links []Link
	var err error
	htmlPath, _ := filepath.Abs("../ex4.html")
	var htmlBytes []byte
	if htmlBytes, err = ioutil.ReadFile(htmlPath); err != nil {
		panic("couldn't read input file" + err.Error())
	}
	if links, err = s.unit.Parse(string(htmlBytes)); err != nil {
		panic(err)
	}

	testHref := []string{"/dog-cat"}
	testText := []string{"dog cat"}

	s.gomega.Expect(len(links)).Should(Equal(1))
	for i, link := range links {

		s.gomega.Expect(link.Href).Should(Equal(testHref[i]))
		s.gomega.Expect(link.Text).Should(Equal(testText[i]))
	}
}
