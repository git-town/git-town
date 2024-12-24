package configdomain

// defines the type of Git configuration used
type ConfigScope string

const (
	ConfigScopeGlobal ConfigScope = "global"
	ConfigScopeLocal  ConfigScope = "local"
)

func (self ConfigScope) String() string {
	return string(self)
}
