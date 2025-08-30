package helpers

import messages "github.com/cucumber/messages/go/v21"

// hasTag indicates whether the given feature has a tag with the given name.
func HasTag(tags []*messages.PickleTag, name string) bool {
	for _, tag := range tags {
		if tag.Name == name {
			return true
		}
	}
	return false
}
