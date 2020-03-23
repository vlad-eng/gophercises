package normalizer

import (
	"regexp"
)

type Type string

const (
	LandLine Type = "LandLine"
	Mobile   Type = "Mobile"
)

const NormalizationPattern = "\\d"
const RemovedCharsPattern = "\\D"

type PhoneNumber struct {
	ID         int `sql:"primary_key;unique"`
	Number     string
	Owner      string
	NumberType string
}

func (p *PhoneNumber) Normalize() {
	var regExpr *regexp.Regexp
	regExpr = regexp.MustCompile(RemovedCharsPattern)
	p.Number = regExpr.ReplaceAllString(p.Number, "")
}
