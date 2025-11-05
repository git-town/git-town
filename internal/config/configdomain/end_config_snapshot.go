package configdomain

// EndConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type EndConfigSnapshot struct {
	Global SingleSnapshot
	Local  SingleSnapshot
}
