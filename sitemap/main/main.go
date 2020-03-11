package main

import (
	"flag"
	"fmt"
	. "gophercises/sitemap/sitemap"
)

func main() {
	domainPtr := flag.String("domain", "http://www.wikipedia.org", "Domain to be searched")
	flag.Parse()

	siteParser := InitSiteMap(*domainPtr)
	links := siteParser.Parse()
	fmt.Println(XmlEncode(links))
}
