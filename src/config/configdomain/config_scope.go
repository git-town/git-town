package configdomain

type ConfigScope int

const (
	ConfigScopeGlobal ConfigScope = iota
	ConfigScopeLocal
)
