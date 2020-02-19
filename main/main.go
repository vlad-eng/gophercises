package main

import (
	"fmt"
	"gophercises/urlshort"
	"net/http"
)

func main() {
	fallbackLocation := "/"
	mux := urlshort.DefaultMux(fallbackLocation)

	// Build the YAMLHandler using the mux as the fallback
	yaml := "mappings:\n - path: /urlshort\n   url: https://github.com/gophercises/urlshort\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution"
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mux, fallbackLocation)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on localhost:8088")
	http.ListenAndServe("localhost:8088", yamlHandler)
}
