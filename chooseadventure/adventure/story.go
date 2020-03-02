package adventure

import "fmt"

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Page struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Story map[string]Page

func (s *Story) NextPage(currentArc string, nextArc string) (*Page, error) {
	if currentArc != "" {
		currentPage, ok := (*s)[currentArc]
		if !ok {
			return nil, fmt.Errorf("no page for option: %s", currentArc)
		}
		currentOptions := currentPage.Options
		var foundNextArc bool
		for _, option := range currentOptions {
			if option.Arc == nextArc {
				foundNextArc = true
			}
		}
		if !foundNextArc {
			return nil, fmt.Errorf("unknown option: %s", nextArc)
		}
	}
	nextPage, ok := (*s)[nextArc]
	if !ok {
		return nil, fmt.Errorf("no page for option: %s", nextArc)
	}
	return &nextPage, nil
}
