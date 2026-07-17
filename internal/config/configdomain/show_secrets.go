package configdomain

// ShowSecrets indicates whether a Git Town command should display sensitive information like tokens in its output.
type ShowSecrets bool

func (self ShowSecrets) ShouldShowSecrets() bool {
	return bool(self)
}
