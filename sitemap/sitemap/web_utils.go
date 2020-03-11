package sitemap

import (
	"encoding/xml"
	"fmt"
	"github.com/fatih/set"
	"io/ioutil"
	"net/http"
	"strings"
)

func XmlEncode(links []string) string {
	urlSet := UrlSet{
		Url: []Url{},
	}
	for _, link := range links {
		urlSet.Url = append(urlSet.Url, Url{Loc: link})
	}
	marshalledSiteMap, _ := xml.MarshalIndent(urlSet, "", "   ")
	return xml.Header + string(marshalledSiteMap)
}

func InitSiteMap(domain string) SiteParser {
	siteMap := SiteParser{
		Domain:             domain,
		ToBeProcessedLinks: make(chan string, 1),
		VisitedLinks:       make(chan set.Interface, 1),
	}
	siteMap.ToBeProcessedLinks <- domain
	visitedLinks := set.New(set.ThreadSafe)
	siteMap.VisitedLinks <- visitedLinks
	return siteMap
}

func formatUrl(domain string, url string) (formattedUrl string) {
	formattedUrl = url
	if strings.HasPrefix(url, "/") {
		formattedUrl = domain + url
	}
	return strings.TrimSpace(formattedUrl)
}

func getPage(url string) (string, error) {
	var response *http.Response
	var err error
	if response, err = http.Get(url); err != nil {
		return "", fmt.Errorf("couldn't fetch page at url %s due to: %s", url, err)
	}
	defer response.Body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
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

func isValidLink(domain string, link string) bool {
	return !isUntraversableLink(link) && isSpecifiedCategoryLink(domain, link, "")
}

func getDomain(url string) string {
	domainTokens := strings.Split(url, ".")
	domain := strings.Split(strings.Join(domainTokens[1:], "."), "/")[0]
	return domain
}
