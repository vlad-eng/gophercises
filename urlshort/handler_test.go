package urlshort

import (
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"log"
	. "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type UrlShortenerTestSuite struct {
	suite.Suite
	gomega *GomegaWithT
}

func TestUrlShortener(t *testing.T) {
	gomega := NewGomegaWithT(t)
	testSuite := UrlShortenerTestSuite{gomega: gomega}
	suite.Run(t, &testSuite)
}

func createRequest() *Request {
	var req *Request
	var err error
	req, err = NewRequest("GET", "localhost:8088", nil)
	req.URL.Path = "/urlshort-godoc"
	if err != nil {
		log.Fatal("couldn't create request due to: ", err)
	}

	return req
}

func (s *UrlShortenerTestSuite) TestMapHandler() {
	req := createRequest()
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux := DefaultMux(fallbackLocation)
	mH := MapHandler(pathsToUrls, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(pathsToUrls[req.URL.Path]))
	s.gomega.Expect(writer.Code).Should(Equal(StatusFound))
}

func (s *UrlShortenerTestSuite) TestFallbackMuxHandler() {
	req := createRequest()
	writer := httptest.NewRecorder()
	fallbackLocation := "/"

	mux := DefaultMux(fallbackLocation)

	mH := MapHandler(nil, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(fallbackLocation))
	s.gomega.Expect(writer.Code).Should(Equal(StatusOK))
}

func (s *UrlShortenerTestSuite) TestYamlHandler() {
	yaml := "mappings:\n - path: /urlshort\n   url: https://gophercises.com/urlshort-godoc\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution"
	fallbackLocation := "/"

	yH, err := YAMLHandler([]byte(yaml), nil, fallbackLocation)
	if err != nil {
		log.Fatal("couldn't create the YAML handler: ", err)
	}
	s.gomega.Expect(yH).ShouldNot(BeNil())
}

func (s *UrlShortenerTestSuite) Test_WhenNoMappingsFallbackIsMuxHandler() {
	var fallbackHandler *ServeMux
	fallbackHandler = NewServeMux()
	fallbackLocation := "/"
	reqPath := "/somepath"
	redirectHandler, redirectLocation := selectHandler(reqPath, nil, fallbackHandler, fallbackLocation)

	s.gomega.Expect(reflect.TypeOf(redirectHandler).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))
	s.gomega.Expect(redirectLocation).Should(Equal(fallbackLocation))
}
