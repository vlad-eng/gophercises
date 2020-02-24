package main

import (
	"fmt"
	. "gophercises/urlshort"
	. "net/http"
)

func main() {
	fallbackLocation := "/"
	mux := DefaultMux(fallbackLocation)

	// Build the YAMLHandler using the mux as the fallback
	//yaml := "mappings:\n - path: /urlshort\n   url: https://github.com/gophercises/urlshort\n - path: /urlshort-final\n   url: https://github.com/gophercises/urlshort/tree/solution\n"
	selectionFunction := CreateSelectionFunction("paths")

	//Using the mysql db handler
	dbHandler, err := DBHandler(selectionFunction, mux, fallbackLocation)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on localhost:8088")
	ListenAndServe("localhost:8088", dbHandler)
}
