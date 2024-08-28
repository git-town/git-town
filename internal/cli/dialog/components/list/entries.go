package list

import (
	"fmt"

	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// Entries provides methods for a collection of Entry instances.
type Entries[S fmt.Stringer] []Entry[S]

// creates an Entries instance containing the given records
func NewEntries[S fmt.Stringer](records ...S) Entries[S] {
	result := make([]Entry[S], len(records))
	for r, record := range records {
		result[r] = Entry[S]{
			Data:    record,
			Enabled: true,
			Text:    record.String(),
		}
	}
	return result
}

// indicates whether all entries in this list are disabled
func (self Entries[S]) AllDisabled() bool {
	for _, entry := range self {
		if entry.Enabled {
			return false
		}
	}
	return true
}

// provides the index of the element that serializes to the given value
func (self Entries[S]) Index(value string) Option[int] {
	for e, entry := range self {
		if entry.Data.String() == value {
			return Some(e)
		}
	}
	return None[int]()
}
