package list

import (
	"fmt"
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
