package sitemap

import (
	"fmt"
	"github.com/fatih/set"
	"github.com/golang/go/src/pkg/strings"
	. "gophercises/link/parser"
	. "io/ioutil"
	. "net/http"
)

type SiteParser struct {
}

//TODO: Check declaring strings and usage in loops - multiple declarations or reassignments?
func (s *SiteParser) Parse(domain string) ([]Link, error) {
	pageLinks := make([]Link, 0)
	toBeVisitedPages := make(chan string, 1000)
	toBeVisitedPages <- domain
	visitedLinks := set.New(set.ThreadSafe)
	linkChannel := make(chan Link)
	finishedRoutinesChannel := make(chan int, 1)
	finishedRoutinesChannel <- 0

	var currentPage string
	var link Link
	var startedRoutines int
	var finishedRoutines int
	for {
		select {
		case currentPage = <-toBeVisitedPages:
			visitedLinks.Add(currentPage)
			startedRoutines++
			go processPage(domain, currentPage, visitedLinks, toBeVisitedPages, linkChannel, finishedRoutinesChannel)

		case link = <-linkChannel:
			pageLinks = append(pageLinks, link)
			for link = range linkChannel {
				pageLinks = append(pageLinks, link)
			}
		case finishedRoutines = <-finishedRoutinesChannel:
			finishedRoutinesChannel <- finishedRoutines
			if startedRoutines > 0 && startedRoutines == finishedRoutines && len(linkChannel) == 0 {
				return pageLinks, nil
			}
		}
	}
}

func processPage(domain string, currentPage string, visitedLinks set.Interface, toBeVisitedPages chan string, linkChannel chan Link, finishedRoutines chan int) {
	if strings.HasPrefix(currentPage, "/") {
		currentPage = domain + currentPage
	}
	var err error
	currentPage = strings.TrimSpace(currentPage)
	var currentHtml string
	currentHtml = getHtml(currentPage)

	childrenLinks := make([]Link, 0)
	if childrenLinks, err = getInternalLinks(domain, currentHtml); err != nil {
		panic(fmt.Errorf("%s", err))
	}

	for _, childLink := range childrenLinks {
		if !visitedLinks.Has(childLink.Href) {
			visitedLinks.Add(childLink.Href)
			toBeVisitedPages <- childLink.Href
			fmt.Printf("siteLink: %s\n", childLink.Href)
			linkChannel <- childLink
		}
	}
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

func getHtml(url string) string {
	var response *Response
	var err error
	if response, err = Get(url); err != nil {
		//return "", fmt.Errorf("couldn't fetch page at url %s due to: %s", url, err)
		panic(fmt.Errorf("couldn't fetch page at url %s due to: %s", url, err))
	}
	defer response.Body.Close()

	var bytes []byte
	if bytes, err = ReadAll(response.Body); err != nil {
		//return "", fmt.Errorf("couldn't read received response: %s", err)
		panic(fmt.Errorf("couldn't read received response: %s", err))
	}
	html := string(bytes)
	//htmlChannel <- html
	//doneChannel <- true
	return html
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
