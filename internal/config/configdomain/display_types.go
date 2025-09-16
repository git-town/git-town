package configdomain

// whether to display branch types in the CLI output
type DisplayTypes bool

func (self DisplayTypes) ShouldDisplayTypes() bool {
	return bool(self)
}
