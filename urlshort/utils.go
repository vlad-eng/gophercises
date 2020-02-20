package urlshort

import (
	"fmt"
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

func buildMap(yaml string) (map[string]string, error) {
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

func SelectHandler(reqPath string, pathsToUrls map[string]string, fallbackHandler Handler, fallbackLocation string) (Handler, string) {
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

func CreateSelectionFunction() HandlerSelectionFunction {
	return func(reqPath string, pathsToUrls map[string]string, fallback Handler, fallbackLocation string) (Handler, string) {
		return SelectHandler(reqPath, pathsToUrls, fallback, fallbackLocation)
	}
}
