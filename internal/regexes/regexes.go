package regexes

import "regexp"

// allows filtering by one of multiple regular expressions
type Regexes []*regexp.Regexp

func NewRegexes(texts []string) (Regexes, error) {
	result := make(Regexes, len(texts))
	for t, text := range texts {
		regex, err := regexp.Compile(text)
		if err != nil {
			return result, err
		}
		result[t] = regex
	}
	return result, nil
}

// indicates whether the given text matches any of the internal regexes
func (self Regexes) Matches(text string) bool {
	if len(self) == 0 {
		return true
	}
	for _, re := range self {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}
