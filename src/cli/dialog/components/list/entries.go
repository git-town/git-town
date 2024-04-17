package list

import "fmt"

// Entries provides methods for a collection of Entry instances.
type Entries[S fmt.Stringer] []Entry[S]

// NewEnabledListEntries creates enabled BubbleListEntries for the given data types.
func NewEntries[S fmt.Stringer](records ...S) []Entry[S] {
	result := make([]Entry[S], len(records))
	for r, record := range records {
		result[r] = Entry[S]{
			Data: record,
			Text: record.String(),
		}
	}
	return result
}
