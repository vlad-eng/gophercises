package sitemap

import (
	"encoding/xml"
	"fmt"
	"github.com/fatih/set"
	. "gophercises/link/parser"
	"time"
)

const timeout = 10

type Url struct {
	XMLName xml.Name `xml:url`
	Loc     string   `xml:loc`
}

type UrlSet struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 Urlset"`
	Url     []Url    `xml:url`
}

type SiteParser struct {
	Domain             string
	ToBeProcessedLinks chan string
	VisitedLinks       chan set.Interface
}

func (s *SiteParser) Parse() []string {
	fmt.Println("Processing, please wait...")

	links := make([]string, 0)
	for {
		select {
		case link := <-s.ToBeProcessedLinks:
			if link != "" {
				links = append(links, link)
			}
			pageToProcess := formatUrl(s.Domain, link)
			go s.processPage(pageToProcess)

		case <-time.After(timeout * time.Second):
			return links
		}
	}
}

func (s *SiteParser) processPage(pageUrl string) {
	var currentHtml string
	var err error
	if currentHtml, err = getPage(pageUrl); err != nil {
		fmt.Println(err)
		return
	}
	pageLinks, _ := s.getPageLinks(currentHtml)
	visitedLinks := <-s.VisitedLinks

	for _, link := range pageLinks {
		if !visitedLinks.Has(link) {
			visitedLinks.Add(link)
			s.ToBeProcessedLinks <- link
		}
	}
	s.VisitedLinks <- visitedLinks
}

func (s *SiteParser) getPageLinks(htmlPage string) ([]string, error) {
	pageParser := PageParser{}
	links := make([]Link, 0)
	pageLinks := make([]string, 0)
	var err error
	if links, err = pageParser.Parse(htmlPage); err != nil {
		return nil, fmt.Errorf("error while parsing a htmlPage page: %s", err)
	}
	for _, link := range links {
		if isValidLink(s.Domain, link.Href) {
			pageLinks = append(pageLinks, link.Href)
		}
	}
	return pageLinks, nil
}
