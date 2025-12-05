package configdomain

// Redact indicates whether a Git Town command should redact sensitive information from the output.
type Redact bool

func (r Redact) ShouldRedact() bool {
	return bool(r)
}
