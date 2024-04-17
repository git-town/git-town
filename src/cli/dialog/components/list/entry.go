package list

import "fmt"

// Entry is an entry in a List instance.
type Entry[S fmt.Stringer] struct {
	Data    S
	Enabled bool
	Text    string
}
