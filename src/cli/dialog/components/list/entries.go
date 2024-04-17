package list

import "fmt"

type Entries[S fmt.Stringer] []Entry[S]

func (self Entries[S]) allDisabled() bool {
	for _, entry := range self {
		if entry.Enabled {
			return false
		}
	}
	return true
}
