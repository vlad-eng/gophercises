package adventure

import (
	"fmt"
	. "html/template"
	"strings"
)

const htmlStructure = `
	<html>
		<head>
			<title> %s </title>
		</head>
		<body>
			<p>
				%s
			</p>
			<p>
				%s
			</p>
		</body>
	</html>
`

type StoryTemplate struct {
	structure string
	title     string
	body      string
	links     []string
}

func (t *StoryTemplate) toHTML() HTML {
	return HTML(t.String())
}

func (t *StoryTemplate) String() string {
	links := t.TemplateLinks()
	return fmt.Sprintf(htmlStructure, t.title, t.body, links)
}

func (t *StoryTemplate) TemplateLinks() string {
	linkStrings := make([]string, len(t.links))
	for i, link := range t.links {
		linkStrings[i] = "<a href='/" + link + "'>" + link + "</a>"
	}
	return strings.Join(linkStrings, "\t")
}
