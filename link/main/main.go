package main

import (
	"fmt"
	"github.com/golang/go/src/pkg/strconv"
	. "gophercises/link/parser"
	"io/ioutil"
	"path/filepath"
)

func main() {
	linkParser := PageParser{}
	var extractedLinks []Link
	var err error
	var htmlBytes []byte

	for i := 1; i < 4; i++ {
		htmlPath, _ := filepath.Abs("ex" + strconv.Itoa(i) + ".html")
		if htmlBytes, err = ioutil.ReadFile(htmlPath); err != nil {
			panic("couldn't read input file" + err.Error())
		}
		htmlString := string(htmlBytes)
		if extractedLinks, err = linkParser.Parse(htmlString); err != nil {
			panic(err)
		}
		for _, link := range extractedLinks {
			fmt.Printf("Link: %s\n", link)
		}
	}

}
