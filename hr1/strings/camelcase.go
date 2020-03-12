package strings

import (
	"fmt"
	"github.com/fatih/set"
	"strings"
	"unicode"
)

type StringTokenizer struct {
}

func (st *StringTokenizer) Tokenize(camelcase string) int {
	if len(camelcase) == 0 {
		return 0
	}
	if unicode.IsUpper(rune(camelcase[0])) {
		panic(fmt.Errorf("not a camel case string"))
	}
	wordCount := 1
	camelcase = strings.TrimSpace(camelcase)
	isAnalyzableChar := true
	for i, ch := range camelcase {
		if isAnalyzableChar && unicode.IsUpper(ch) == true {
			wordCount++
		}
		if strings.EqualFold(string(ch), " ") {
			isAnalyzableChar = false
		}
		if isAdditionalCamelCaseString(camelcase, i) {
			wordCount++
			isAnalyzableChar = true
		}
	}
	return wordCount
}

func isAdditionalCamelCaseString(camelcase string, index int) bool {
	separators := *getSeparators()
	return index > 0 &&
		isSeparator(string(camelcase[index-1]), separators) &&
		unicode.IsLower(rune(camelcase[index])) == true
}

func isSeparator(ch string, p set.Interface) bool {
	for _, separator := range p.List() {
		if strings.EqualFold(separator.(string), ch) {
			return true
		}
	}
	return false
}

func getSeparators() *set.Interface {
	separators := set.New(set.ThreadSafe)
	separators.Add(" ")
	separators.Add(",")
	separators.Add(";")
	separators.Add(".")
	separators.Add("\t")
	return &separators
}
