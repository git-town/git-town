package tags

import (
	"strings"

	"github.com/cucumber/messages-go/v16"
)

// ApplyTagFilter will apply a filter string on the
// array of pickles and returned the filtered list.
func ApplyTagFilter(filter string, pickles []*messages.Pickle) []*messages.Pickle {
	if filter == "" {
		return pickles
	}

	var result = []*messages.Pickle{}

	for _, pickle := range pickles {
		if match(filter, pickle.Tags) {
			result = append(result, pickle)
		}
	}

	return result
}

// Based on http://behat.readthedocs.org/en/v2.5/guides/6.cli.html#gherkin-filters
func match(filter string, tags []*messages.PickleTag) (ok bool) {
	ok = true

	for _, andTags := range strings.Split(filter, "&&") {
		var okComma bool

		for _, tag := range strings.Split(andTags, ",") {
			tag = strings.TrimSpace(tag)
			tag = strings.Replace(tag, "@", "", -1)

			okComma = contains(tags, tag) || okComma

			if tag[0] == '~' {
				tag = tag[1:]
				okComma = !contains(tags, tag) || okComma
			}
		}

		ok = ok && okComma
	}

	return
}

func contains(tags []*messages.PickleTag, tag string) bool {
	for _, t := range tags {
		tagName := strings.Replace(t.Name, "@", "", -1)

		if tagName == tag {
			return true
		}
	}

	return false
}
