package adventure

import (
	. "net/http"
	"strings"
)

type StoryServer struct {
	PreviousRequest string
	Story           *Story
}

func (s *StoryServer) ServeHTTP(w ResponseWriter, r *Request) {
	var currentKey string
	if len(s.PreviousRequest) > 0 {
		currentKey = s.PreviousRequest[1:]
	}
	nextKey := r.URL.Path[1:]

	if page, err := s.Story.NextPage(currentKey, nextKey); err != nil {
		Error(w, "couldn't retrieve story page for the chosen option: "+err.Error(), StatusInternalServerError)
	} else {
		renderedPage := StoryTemplate{
			structure: htmlStructure,
			title:     page.Title,
			body:      strings.Join(page.Story, "\n"),
		}
		links := make([]string, len(page.Options))
		for i, option := range page.Options {
			links[i] = option.Arc
		}
		renderedPage.links = links
		w.Write([]byte(renderedPage.toHTML()))
	}
}
