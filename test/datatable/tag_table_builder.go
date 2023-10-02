package datatable

import (
	"sort"

	"github.com/git-town/git-town/v9/test/helpers"
)

// TagTableBuilder collects data about tags in Git repositories
// in the same way that our Gherkin tables describing tags in repos are organized.
type TagTableBuilder struct {
	// tagsToLocations stores data about what locations tags are.
	//
	// Structure:
	//   tag1: [local]
	//   tag2: [local, remote]
	tagToLocations map[string]helpers.OrderedSet[string]
}

// NewTagTableBuilder provides a fully initialized instance of TagTableBuilder.
func NewTagTableBuilder() TagTableBuilder {
	return TagTableBuilder{
		tagToLocations: make(map[string]helpers.OrderedSet[string]),
	}
}

// Add registers the given tag from the given location into this table.
func (ttb *TagTableBuilder) Add(tag, location string) {
	locations, exists := ttb.tagToLocations[tag]
	if exists {
		ttb.tagToLocations[tag] = locations.Add(location)
	} else {
		ttb.tagToLocations[tag] = helpers.NewOrderedSet(location)
	}
}

// AddMany registers the given tags from the given location into this table.
func (ttb *TagTableBuilder) AddMany(tags []string, location string) {
	for _, tag := range tags {
		ttb.Add(tag, location)
	}
}

// Table provides the data accumulated by this TagTableBuilder as a DataTable.
func (ttb *TagTableBuilder) Table() DataTable {
	result := DataTable{}
	result.AddRow("NAME", "LOCATION")
	tags := make([]string, len(ttb.tagToLocations))
	index := 0
	for tag := range ttb.tagToLocations {
		tags[index] = tag
		index++
	}
	sort.Strings(tags)
	for _, tag := range tags {
		result.AddRow(tag, ttb.tagToLocations[tag].Join(", "))
	}
	return result
}
