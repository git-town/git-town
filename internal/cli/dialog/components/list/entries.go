package list

import "fmt"

// Entries provides methods for a collection of Entry instances.
type Entries[S fmt.Stringer] []Entry[S]

// NewEnabledListEntries creates Entries for the given data types.
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

// AllDisabled indicates whether all entries in this list are disabled.
func (self Entries[S]) AllDisabled() bool {
	for _, entry := range self {
		if entry.Enabled {
			return false
		}
	}
	return true
}

// IndexWithText provides the index of the element with the given text.
func (self Entries[S]) IndexWithText(text string) (found bool, index int) {
	for e, entry := range self {
		if entry.Data.String() == text {
			return true, e
		}
	}
	return false, 0
}
