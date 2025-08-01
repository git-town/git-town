package list

import (
	"fmt"

	"github.com/git-town/git-town/v21/pkg/equal"
)

// Entries provides methods for a collection of Entry instances.
type Entries[S any] []Entry[S]

// creates an Entries instance containing the given records
func NewEntries[S fmt.Stringer](records ...S) Entries[S] {
	result := make([]Entry[S], len(records))
	for r, record := range records {
		result[r] = Entry[S]{
			Data:     record,
			Disabled: false,
			Text:     record.String(),
		}
	}
	return result
}

// indicates whether all entries in this list are disabled
func (self Entries[S]) AllDisabled() bool {
	for _, entry := range self {
		if !entry.Disabled {
			return false
		}
	}
	return true
}

// provides the index of the first entry that is not disabled
func (self Entries[S]) FirstEnabled() int {
	for e, entry := range self {
		if !entry.Disabled {
			return e
		}
	}
	return 0
}

// provides the position of the given needle in this list
func (self Entries[S]) IndexOf(needle S) int {
	for e, entry := range self {
		if equal.Equal(entry.Data, needle) {
			return e
		}
	}
	return 0
}

// provides the position of the given needle in this list
func (self Entries[S]) IndexOfFunc(needle S, equalFn func(a, b S) bool) int {
	for e, entry := range self {
		if equalFn(entry.Data, needle) {
			return e
		}
	}
	return 0
}
