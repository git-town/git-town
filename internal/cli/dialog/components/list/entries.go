package list

import (
	"fmt"
)

// Entries provides methods for a collection of Entry instances.
type Entries[S comparable] []Entry[S]

// creates an Entries instance containing the given records
func NewEntries[S ComparableStringer](records ...S) Entries[S] {
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

// provides the position of the given needle in this list
func (self Entries[S]) IndexOfData(needle S) int {
	for e, entry := range self {
		if entry.Data == needle {
			return e
		}
	}
	return -1
}

// narrower type needed to use the NewEntries convenience function
type ComparableStringer interface {
	comparable
	fmt.Stringer
}
