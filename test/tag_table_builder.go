package test

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v7/test/helpers"
)

// TagTableBuilder collects data about tags in Git repositories
// in the same way that our Gherkin tables describing tags in repos are organized.
type TagTableBuilder struct {
	// tagsToLocations stores data about what locations tags are.
	//
	// Structure:
	//   tag1: [local]
	//   tag2: [local, remote]
	tagToLocations map[string]helpers.OrderedStringSet
}

// NewTagTableBuilder provides a fully initialized instance of TagTableBuilder.
func NewTagTableBuilder() TagTableBuilder {
	return TagTableBuilder{
		tagToLocations: make(map[string]helpers.OrderedStringSet),
	}
}

// Add registers the given tag from the given location into this table.
func (builder *TagTableBuilder) Add(tag, location string) {
	locations, exists := builder.tagToLocations[tag]
	if exists {
		builder.tagToLocations[tag] = locations.Add(location)
	} else {
		builder.tagToLocations[tag] = helpers.NewOrderedStringSet(location)
	}
}

// AddMany registers the given tags from the given location into this table.
func (builder *TagTableBuilder) AddMany(tags []string, location string) {
	for _, tag := range tags {
		builder.Add(tag, location)
	}
}

// Table provides the data accumulated by this TagTableBuilder as a DataTable.
func (builder *TagTableBuilder) Table() (result DataTable) {
	result.AddRow("NAME", "LOCATION")
	tags := make([]string, len(builder.tagToLocations))
	index := 0
	for tag := range builder.tagToLocations {
		tags[index] = tag
		index++
	}
	sort.Strings(tags)
	for _, tag := range tags {
		result.AddRow(tag, strings.Join(builder.tagToLocations[tag].Slice(), ", "))
	}
	return result
}
