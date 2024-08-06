package configdomain

// defines the type of Git configuration used
type ConfigScope int

const (
	// the global Git configuration
	ConfigScopeGlobal ConfigScope = iota

	// the local Git configuration
	ConfigScopeLocal
)
