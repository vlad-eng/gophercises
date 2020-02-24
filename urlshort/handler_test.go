package urlshort

import (
	. "database/sql"
	"fmt"
	"github.com/golang/go/src/pkg/database/sql"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"log"
	. "net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

const tableName = "testPaths"

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

/*
	NOTE:
 	Too many params, alternatives?
*/
func createTable(driverName string, dbName string, tableName string, columns []string, testData [][]string) *DB {
	db := useDB(driverName, dbName)
	//var newRow *Rows
	var err error
	_, _ = db.Query("DROP TABLE " + tableName)
	if _, err = db.Query("CREATE TABLE " + dbName + "." + tableName + "(id INTEGER, " + columns[0] + " VARCHAR(50), " + columns[1] + " VARCHAR(100));"); err != nil {
		fmt.Println(err.Error())
	}

	//var insert *Rows
	var testRow []string
	var i int
	if testData != nil && len(testData) > 0 {
		for i, testRow = range testData {
			insertCommand := "INSERT INTO " + tableName + "(id" + ", " + columns[0] + ", " + columns[1] +
				") VALUES(" + strconv.FormatInt(int64(i+1), 10) + ", \"" + testRow[0] + "\", \"" + testRow[1] + "\")"
			if _, err = db.Query(insertCommand); err != nil {
				panic(err.Error())
			}
		}
	}

	return db
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
	mH := MapHandler(pathsToUrls, CreateSelectionFunction(tableName), mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(pathsToUrls[req.URL.Path]))
	s.gomega.Expect(writer.Code).Should(Equal(StatusFound))
}

func (s *UrlShortenerTestSuite) TestDBHandler_WhenMappingsInDB() {
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)
	checkTypesSelectionFunction := func(reqPath string, dataStore interface{}, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(dataStore).String()).Should(Equal(reflect.TypeOf(&sql.DB{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, dataStore, tableName, fallback, fallbackLocation)

		s.gomega.Expect(reflect.TypeOf(selectedHandler).String()).ShouldNot(Equal(reflect.TypeOf(&ServeMux{}).String()))
		return selectedHandler, redirectLocation
	}

	driverName := "mysql"
	dbName := "mysql"
	tableName := "testpaths"
	columns := []string{"path", "redirectionUrl"}
	testData := [][]string{{"/urlshort", "https://gophercises.com/urlshort-godoc"},
		{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"}}
	ds := createTable(driverName, dbName, tableName, columns, testData)

	mH := MapHandler(ds, checkTypesSelectionFunction, mux, fallbackLocation)
	req := createRequest("/urlshort")
	writer := httptest.NewRecorder()
	mH.ServeHTTP(writer, req)

	defer ds.Close()
	s.gomega.Expect(mH).ShouldNot(BeNil())
}

func (s *UrlShortenerTestSuite) Test_WhenNoMappingsFallbackIsMuxHandler() {
	req := createRequest("/urlshort")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	checkTypesSelectionFunction := func(reqPath string, dataStore interface{}, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(fallback).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, dataStore, tableName, fallback, fallbackLocation)

		s.gomega.Expect(reflect.TypeOf(selectedHandler).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))
		return selectedHandler, redirectLocation
	}

	dataStore := make(map[string]string)
	mH := MapHandler(dataStore, checkTypesSelectionFunction, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(fallbackLocation))
	s.gomega.Expect(writer.Code).Should(Equal(StatusOK))
}

func (s *UrlShortenerTestSuite) Test_WhenUnrecognizedRequestPathFallbackIsMuxHandler() {
	req := createRequest("/some-unrecognized-path")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	checkTypesSelectionFunction := func(reqPath string, dataStore interface{}, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(fallback).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, dataStore, tableName, fallback, fallbackLocation)

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

func (s *UrlShortenerTestSuite) Test_WhenNoDBMappingsRequestPathFallbackIsMuxHandler() {
	req := createRequest("/urlshort")
	writer := httptest.NewRecorder()
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	checkTypesSelectionFunction := func(reqPath string, dataStore interface{}, fallback Handler, fallbackLocation string) (Handler, string) {
		s.gomega.Expect(reflect.TypeOf(fallback).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))

		selectedHandler, redirectLocation := SelectHandler(reqPath, dataStore, tableName, fallback, fallbackLocation)

		s.gomega.Expect(reflect.TypeOf(selectedHandler).String()).Should(Equal(reflect.TypeOf(&ServeMux{}).String()))
		return selectedHandler, redirectLocation
	}

	driverName := "mysql"
	dbName := "mysql"
	tableName := "testpaths"
	columns := []string{"path", "redirectionUrl"}
	testData := make([][]string, 0)
	ds := createTable(driverName, dbName, tableName, columns, testData)

	mH := MapHandler(ds, checkTypesSelectionFunction, mux, fallbackLocation)
	mH.ServeHTTP(writer, req)

	defer ds.Close()
	s.gomega.Expect(mH).ShouldNot(BeNil())
	s.gomega.Expect(writer.Header().Get("Location")).Should(Equal(fallbackLocation))
	s.gomega.Expect(writer.Code).Should(Equal(StatusOK))
}

func (s *UrlShortenerTestSuite) TestYamlHandler() {
	yaml := "mappings:\n - path: /urlshort\n   url: https://gophercises.com/urlshort-godoc\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution"
	fallbackLocation := "/"

	yH, err := YAMLHandler([]byte(yaml), CreateSelectionFunction(tableName), nil, fallbackLocation)
	if err != nil {
		log.Fatal("couldn't create the YAML handler: ", err)
	}
	s.gomega.Expect(yH).ShouldNot(BeNil())
}
