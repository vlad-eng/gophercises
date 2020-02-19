package main

import (
	"fmt"
	"net/http"
	"urlshort"
)

func main() {
	fallbackLocation := "/"
	mux := urlshort.DefaultMux(fallbackLocation)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux, fallbackLocation)
	// Build the YAMLHandler using the mapHandler as the fallback
	yaml := "mappings:\n - path: /urlshort\n   url: https://gophercises.com/urlshort-godoc\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution"
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler, fallbackLocation)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on localhost:8088")
	http.ListenAndServe("localhost:8088", yamlHandler)
}
