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

func (s *UrlShortenerTestSuite) TestMapHandler() {
	req := createRequest("/urlshort")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	pathsToUrls := map[string]string{
		"/urlshort":       "https://gophercises.com/urlshort-godoc",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	mux := DefaultMux(fallbackLocation)
	mH := MapHandler(pathsToUrls, CreateSelectionFunction(), mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(pathsToUrls[req.URL.Path]))
	s.gomega.Expect(writer.Code).Should(Equal(StatusFound))
}

func (s *UrlShortenerTestSuite) Test_WhenNoMappingsFallbackIsMuxHandler() {
	req := createRequest("/urlshort")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	checkTypesSelectionFunction := func(reqPath string, pathsToUrls map[string]string, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(fallback).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, pathsToUrls, fallback, fallbackLocation)

		s.gomega.Expect(reflect.TypeOf(selectedHandler).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))
		return selectedHandler, redirectLocation
	}

	mH := MapHandler(nil, checkTypesSelectionFunction, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(fallbackLocation))
	s.gomega.Expect(writer.Code).Should(Equal(StatusOK))
}

func (s UrlShortenerTestSuite) Test_WhenUnrecognizedRequestPathFallbackIsMuxHandler() {
	req := createRequest("/some-unrecognized-path")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	checkTypesSelectionFunction := func(reqPath string, pathsToUrls map[string]string, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(fallback).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, pathsToUrls, fallback, fallbackLocation)

		s.gomega.Expect(reflect.TypeOf(selectedHandler).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))
		return selectedHandler, redirectLocation
	}

	pathsToUrls := map[string]string{
		"/urlshort":       "https://gophercises.com/urlshort-godoc",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	mH := MapHandler(pathsToUrls, checkTypesSelectionFunction, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)
}

func (s *UrlShortenerTestSuite) TestYamlHandler() {
	yaml := "mappings:\n - path: /urlshort\n   url: https://gophercises.com/urlshort-godoc\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution"
	fallbackLocation := "/"

	yH, err := YAMLHandler([]byte(yaml), CreateSelectionFunction(), nil, fallbackLocation)
	if err != nil {
		log.Fatal("couldn't create the YAML handler: ", err)
	}
	s.gomega.Expect(yH).ShouldNot(BeNil())
}
