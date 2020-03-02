package main

import (
	"fmt"
	. "gophercises/chooseadventure/adventure"
	"net/http"
)

func main() {
	var story *Story
	var err error
	if story, err = Decode("adventure/gopher.json"); err != nil {
		fmt.Printf("couldn't load from the given file due to: %s\n", err)
	}

	storyServer := StoryServer{
		PreviousRequest: "",
		Story:           story,
	}
	http.ListenAndServe("localhost:8088", &storyServer)
}
