package list

// Entry is an entry in a List instance.
type Entry[S any] struct {
	Data S
	//exhaustruct:optional
	Disabled bool
	Text     string
}
