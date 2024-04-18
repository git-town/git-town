package list

import "fmt"

// Entry is an entry in a List instance.
type Entry[S fmt.Stringer] struct {
	Data    S
	Enabled bool
	Text    string
}

func (self Entry[S]) String() string {
	return self.Text
}
