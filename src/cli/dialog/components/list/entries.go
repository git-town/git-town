package list

import "fmt"

// Entries provides methods for a collection of Entry instances.
type Entries[S fmt.Stringer] []Entry[S]

func (self Entries[S]) allDisabled() bool {
	for _, entry := range self {
		if entry.Enabled {
			return false
		}
	}
	return true
}

// NewEnabledListEntries creates enabled BubbleListEntries for the given data types.
func NewEnabledListEntries[S fmt.Stringer](records []S) []Entry[S] {
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
