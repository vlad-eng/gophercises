package sitemap

import (
	"encoding/xml"
	"fmt"
	"github.com/fatih/set"
	"github.com/golang/go/src/pkg/strings"
	. "gophercises/link/parser"
	. "io/ioutil"
	. "net/http"
	"time"
)

type Url struct {
	XMLName xml.Name `xml:url`
	Loc     string   `xml:loc`
}

type UrlSet struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 Urlset"`
	Url     []Url    `xml:url`
}

type SiteParser struct {
}

//TODO: Check declaring strings and usage in loops - multiple declarations or reassignments?
func (s *SiteParser) Parse(domain string, outputLinks chan Link, doneChannel chan bool) {
	linkChannel := make(chan Link, 1000)
	linkChannel <- Link{Href: domain, Text: ""}
	finishedRoutinesChannel := make(chan int, 1)
	finishedRoutinesChannel <- 0
	visitedLinks := set.New(set.ThreadSafe)
	visitedChannel := make(chan set.Interface, 1)
	visitedChannel <- visitedLinks

	var link Link
	var startedRoutines int
	for {
		select {
		case link = <-linkChannel:
			startedRoutines++
			if link.Text != "" {
				outputLinks <- link
			}
			go processPage(domain, link.Href, visitedChannel, linkChannel, finishedRoutinesChannel)

		case <-time.After(10 * time.Second):
			doneChannel <- true
		}
	}
}

func (s *SiteParser) Format(links []Link) string {
	siteMap := UrlSet{
		Url: []Url{},
	}

	for _, link := range links {
		siteMap.Url = append(siteMap.Url, Url{Loc: link.Href})
	}

	marshalledSiteMap, _ := xml.MarshalIndent(siteMap, "", "   ")
	return xml.Header + string(marshalledSiteMap)
}

func processPage(domain string, currentPage string, visitedChannel chan set.Interface, linkChannel chan Link, finishedRoutines chan int) {
	if strings.HasPrefix(currentPage, "/") {
		currentPage = domain + currentPage
	}
	var err error
	currentPage = strings.TrimSpace(currentPage)
	var currentHtml string
	currentHtml, _ = getHtml(currentPage)

	childrenLinks := make([]Link, 0)
	if childrenLinks, err = getInternalLinks(domain, currentHtml); err != nil {
		panic(fmt.Errorf("%s", err))
	}
	visitedLinks := <-visitedChannel

	for _, childLink := range childrenLinks {
		if !visitedLinks.Has(childLink.Href) {
			visitedLinks.Add(childLink.Href)
			fmt.Printf("siteLink: %s\n", childLink.Href)
			linkChannel <- childLink
		}
	}
	visitedChannel <- visitedLinks
	previouslyFinishedRoutines := <-finishedRoutines
	finishedRoutines <- previouslyFinishedRoutines + 1
}

func getInternalLinks(domain string, html string) ([]Link, error) {
	pageParser := PageParser{}
	links := make([]Link, 0)
	domainLinks := make([]Link, 0)
	var err error
	if links, err = pageParser.Parse(html); err != nil {
		return nil, fmt.Errorf("error while parsing a html page: %s", err)
	}
	for _, link := range links {
		if !isUntraversableLink(link.Href) && isSpecifiedCategoryLink(domain, link.Href, "") {
			domainLinks = append(domainLinks, link)
		}
	}
	return domainLinks, nil
}

func getHtml(url string) (string, error) {
	var response *Response
	var err error
	if response, err = Get(url); err != nil {
		return "", fmt.Errorf("couldn't fetch page at url %s due to: %s", url, err)
	}
	defer response.Body.Close()

	var bytes []byte
	if bytes, err = ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("couldn't read received response: %s", err)
	}
	html := string(bytes)
	return html, nil
}

func getUntraversablePrefixes() set.Interface {
	unTraversableLinks := set.New(set.ThreadSafe)
	unTraversableLinks.Add("#")
	return unTraversableLinks
}

func isUntraversableLink(href string) bool {
	untraversablePrefixes := getUntraversablePrefixes().List()
	var isUntraversable bool
	for _, prefix := range untraversablePrefixes {
		isUntraversable = isUntraversable || strings.HasPrefix(href, prefix.(string))
	}
	return isUntraversable
}

func isInternalLink(url string, href string) bool {
	return isSpecifiedCategoryLink(url, href, "")
}

func isSpecifiedCategoryLink(url string, href string, category string) bool {
	domain := getDomain(url)
	if category != "" {
		domain = domain + "/" + category
	}
	if strings.HasPrefix(href, "http") &&
		strings.Contains(href, domain) {
		return true
	}
	return false
}

func getDomain(url string) string {
	domainTokens := strings.Split(url, ".")
	domain := strings.Split(strings.Join(domainTokens[1:], "."), "/")[0]
	return domain
}
