package configdomain

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type BeginConfigSnapshot struct {
	Global   SingleSnapshot
	Local    SingleSnapshot
	Unscoped SingleSnapshot
}
