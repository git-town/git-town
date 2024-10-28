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

func ParseConfigScope(text string) ConfigScope {
	switch text {
	case "local", "":
		return ConfigScopeLocal
	case "global":
		return ConfigScopeGlobal
	default:
		panic("unknown locality: " + text)
	}
}
