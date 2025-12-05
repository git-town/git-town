package configdomain

// Redact indicates whether a Git Town command should redact API tokens from the output.
type Redact bool

func (r Redact) ShouldRedact() bool {
	return bool(r)
}
