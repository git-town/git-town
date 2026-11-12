package tags

import (
	"strings"

	messages "github.com/cucumber/messages/go/v34"
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
func match(filters string, tags []*messages.PickleTag) (ok bool) {
	ok = false
	for _, filter := range strings.Split(filters, ",") {
		ok = matchAnd(filter, tags) || ok
	}

	return
}

func matchAnd(filter string, tags []*messages.PickleTag) bool {
	and := len(filter) > 0
	for _, tag := range strings.Split(filter, "&&") {
		tag = strings.TrimSpace(tag)
		tag = strings.Replace(tag, "@", "", -1)
		if tag[0] == '~' {
			tag = tag[1:]
			and = !contains(tags, tag) && and
		} else {
			and = contains(tags, tag) && and
		}
	}
	return and
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
