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

// IndexWithText provides the index of the element with the given text.
func (self Entries[S]) IndexWithText(text string) (found bool, index int) {
	for e := range self {
		if self[e].Text == text {
			return true, e
		}
	}
	return false, 0
}

// IndexWithTextOr provides the index of the element with the given text
// or the given default index if the element isn't in this collection.
func (self Entries[S]) IndexWithTextOr(text string, defaultIndex int) int {
	found, index := self.IndexWithText(text)
	if found {
		return index
	}
	return defaultIndex
}
