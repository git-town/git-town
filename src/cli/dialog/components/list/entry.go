package list

import "fmt"

type Entry[S fmt.Stringer] struct {
	// TODO: Checked bool in a RadioListEntry struct
	Data    S
	Enabled bool
	Text    string
}
