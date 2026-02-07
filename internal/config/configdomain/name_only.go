package configdomain

// NameOnly indicates whether git diff-parent should display only the names of changed files.
type NameOnly bool

func (self NameOnly) DisplayNameOnly() bool {
	return bool(self)
}
