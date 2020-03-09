package main

import (
	"flag"
	"fmt"
	. "gophercises/link/parser"
	. "gophercises/sitemap/sitemap"
)

func main() {
	domainPtr := flag.String("domain", "http://www.wikipedia.org", "Domain to be searched")
	flag.Parse()

	siteParser := SiteParser{}
	visitedLinks := make(chan Link, 0)
	doneChannel := make(chan bool, 1)
	go siteParser.Parse(*domainPtr, visitedLinks, doneChannel)
	links := make([]Link, 0)
	for {
		select {
		case link := <-visitedLinks:
			links = append(links, link)

		case <-doneChannel:
			fmt.Println(siteParser.Format(links))
			return
		}
	}
}
