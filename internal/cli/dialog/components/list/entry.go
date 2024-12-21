package list

// Entry is an entry in a List instance.
type Entry[S comparable] struct {
	Data     S
	Disabled bool `exhaustruct:"optional"`
	Text     string
}
