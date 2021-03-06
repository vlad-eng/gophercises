package urlshort

import (
	. "net/http"
)

type HandlerSelectionFunction func(string, interface{}, Handler, string) (Handler, string)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(dataStore interface{}, selectionFunction HandlerSelectionFunction, fallback Handler, fallbackLocation string) HandlerFunc {
	handlerFunc := func(w ResponseWriter, r *Request) {
		reqPath := r.URL.Path

		sw := ServerWriter{
			w: w,
		}

		redirectHandler, redirectLocation := selectionFunction(reqPath, dataStore, fallback, fallbackLocation)
		sw.location = redirectLocation
		redirectHandler.ServeHTTP(sw.w, r)
		if sw.wroteHeader == false {
			sw.w.Header().Set("Location", sw.location)
			sw.wroteHeader = true
		}
	}
	return handlerFunc
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, selectionFunction HandlerSelectionFunction, fallback Handler, fallbackLocation string) (HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap, err := buildMapFromYaml(parsedYaml)
	return MapHandler(pathMap, selectionFunction, fallback, fallbackLocation), err
}

func DBHandler(selectionFunction HandlerSelectionFunction, fallback Handler, fallbackLocation string) (HandlerFunc, error) {
	db := useDB("mysql", "paths")
	pathMap, err := buildMapFromDB(db, "paths")
	return MapHandler(pathMap, selectionFunction, fallback, fallbackLocation), err
}
