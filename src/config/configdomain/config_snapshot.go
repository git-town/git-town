package configdomain

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	GitConfig FullCache
}
