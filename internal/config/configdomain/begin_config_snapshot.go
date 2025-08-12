package configdomain

import "fmt"

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type BeginConfigSnapshot struct {
	Global   SingleSnapshot
	Local    SingleSnapshot
	Unscoped SingleSnapshot
}

func (self BeginConfigSnapshot) ByScope(scope ConfigScope) SingleSnapshot {
	switch scope {
	case ConfigScopeGlobal:
		return self.Global
	case ConfigScopeLocal:
		return self.Local
	}
	panic(fmt.Sprintf("unknown config scope: %q", scope))
}
