package main

import (
	"flag"
	"fmt"
	. "gophercises/sitemap/sitemap"
)

func main() {
	domainPtr := flag.String("url", "http://www.wikipedia.org", "The URL you want to build a sitemap for")
	flag.Parse()

	siteParser := InitSiteParser(*domainPtr)
	links := siteParser.Parse()
	fmt.Println(XmlEncode(links))
}
