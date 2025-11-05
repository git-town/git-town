package configdomain

import "fmt"

// BeginConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type BeginConfigSnapshot struct {
	Global   SingleSnapshot
	Local    SingleSnapshot
	Unscoped SingleSnapshot
}

// ByScope looks up one of the contained SingleSnapshots by scope.
func (self BeginConfigSnapshot) ByScope(scope ConfigScope) SingleSnapshot {
	switch scope {
	case ConfigScopeGlobal:
		return self.Global
	case ConfigScopeLocal:
		return self.Local
	}
	panic(fmt.Sprintf("unknown config scope: %q", scope))
}
