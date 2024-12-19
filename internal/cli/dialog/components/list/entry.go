package list

// Entry is an entry in a List instance.
type Entry[S comparable] struct {
	Data     S
	Disabled bool // TODO: reverse to Disabled, make optional, then remove from all the instantiations of Entry
	Text     string
}
