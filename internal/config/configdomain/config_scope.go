package configdomain

// defines the type of Git configuration used
type ConfigScope int

const (
	// the global Git configuration
	ConfigScopeGlobal ConfigScope = iota

	// the local Git configuration
	ConfigScopeLocal
)

func (self ConfigScope) String() string {
	switch self {
	case ConfigScopeGlobal:
		return "global"
	case ConfigScopeLocal:
		return "local"
	}
	panic("unknown scope")
}
