package main

import (
	"bufio"
	"fmt"
	"gophercises/hr1/strings"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	st := strings.StringTokenizer{}
	description :=
		`This program counts how many simple words are in camel case phrases.
Note: Camel case phrases start with a camel case word but may contain non camel case words.`
	fmt.Println(description)
	for i := 0; i < 5; i++ {
		scanner.Scan()
		input := scanner.Text()
		fmt.Printf("input: %s; word count: %d\n", input, st.Tokenize(input))
	}
}
