package urlshort

import (
	"fmt"
	. "gopkg.in/yaml.v2"
	"net/http"
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

	pathToUrl := make(map[string]string)
	pathToUrl["path"] = config.Mappings[0].Path
	pathToUrl["url"] = config.Mappings[0].Url

	return pathToUrl, nil
}

func parseYAML(yaml []byte) (string, error) {
	return string(yaml), nil
}

func DefaultMux(fallbackLocation string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(fallbackLocation, muxFallbackHandler)
	return mux
}

func muxFallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
