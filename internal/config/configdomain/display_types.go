package configdomain

// whether to display branch types in the CLI output
type DisplayTypes bool

func (self DisplayTypes) IsTrue() bool {
	return bool(self)
}
