package configdomain

// EndConfigSnapshot is a snapshot of the entire Git configuration after a Git Town command finished.
type EndConfigSnapshot struct {
	Global SingleSnapshot
	Local  SingleSnapshot
}
