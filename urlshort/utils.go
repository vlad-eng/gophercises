package urlshort

import (
	. "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "gopkg.in/yaml.v2"
	. "net/http"
)

type YamlMapping struct {
	Path string
	Url  string
}

type YamlConfig struct {
	Mappings []YamlMapping `yaml:mappings`
}

func buildMapFromYaml(yaml string) (map[string]string, error) {
	config := YamlConfig{}
	config.Mappings = make([]YamlMapping, 0)

	err := Unmarshal([]byte(yaml), &config)
	if err != nil {
		return nil, err
	}

	pathToUrl := make(map[string]string, len(config.Mappings))
	for _, mapping := range config.Mappings {
		pathToUrl[mapping.Path] = mapping.Url
	}

	return pathToUrl, nil
}

func parseYAML(yaml []byte) (string, error) {
	return string(yaml), nil
}

func DefaultMux(fallbackLocation string) *ServeMux {
	mux := NewServeMux()
	mux.HandleFunc(fallbackLocation, muxFallbackHandler)
	return mux
}

func muxFallbackHandler(w ResponseWriter, r *Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func selectHandler(reqPath string, pathsToUrls map[string]string, fallbackHandler Handler, fallbackLocation string) (Handler, string) {
	var redirectHandler Handler
	var redirectLocation string
	if pathsToUrls[reqPath] != "" {
		redirectionUrl := pathsToUrls[reqPath]
		redirectHandler = RedirectHandler(redirectionUrl, StatusFound)
		redirectLocation = redirectionUrl
	} else {
		redirectHandler = fallbackHandler
		redirectLocation = fallbackLocation
	}
	return redirectHandler, redirectLocation
}

func useDB(driverName, dbName string) *DB {
	var db *DB
	var err error
	if db, err = Open(driverName, "root:password@/"+dbName); err != nil {
		panic(err.Error())
	}
	return db
}

func buildMapFromDB(db *DB, tableName string) (map[string]string, error) {
	var err error
	pathsToUrls := make(map[string]string, 0)
	db, err = Open("mysql", "root:password@/mysql")
	defer db.Close()

	var mappings *Rows
	mappings, err = db.Query("SELECT path, redirectionUrl from " + tableName)

	var mappingCounter int
	var path string
	var redirectionUrl string
	for mappings.Next() {
		err = mappings.Scan(&path, &redirectionUrl)
		pathsToUrls[path] = redirectionUrl
		mappingCounter++
	}
	return pathsToUrls, err
}

/*
	NOTES:
	1) Too many params
	2) tableName is not needed in info in string map and not in DB
*/
func SelectHandler(reqPath string, dataStore interface{}, tableName string, fallbackHandler Handler, fallbackLocation string) (Handler, string) {
	db, isDBStore := dataStore.(*DB)
	var pathsToUrls = make(map[string]string, 0)
	if isDBStore {
		var err error
		pathsToUrls, err = buildMapFromDB(db, tableName)
		if err != nil {
			panic(err)
		}
	} else {
		var isMapStore bool
		pathsToUrls, isMapStore = dataStore.(map[string]string)
		if !isMapStore {
			panic("unknown data store")
		}
	}

	return selectHandler(reqPath, pathsToUrls, fallbackHandler, fallbackLocation)
}

/*
	NOTE: tableName has nothing to do with the purpose of the function
*/
func CreateSelectionFunction(tableName string) HandlerSelectionFunction {
	return func(reqPath string, dataStore interface{}, fallback Handler, fallbackLocation string) (Handler, string) {
		return SelectHandler(reqPath, dataStore, tableName, fallback, fallbackLocation)
	}
}
