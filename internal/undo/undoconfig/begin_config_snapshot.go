package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type BeginConfigSnapshot struct {
	Global   configdomain.SingleSnapshot
	Local    configdomain.SingleSnapshot
	Unscoped configdomain.SingleSnapshot
}

func EmptyBeginConfigSnapshot() BeginConfigSnapshot {
	return BeginConfigSnapshot{
		Global:   map[configdomain.Key]string{},
		Local:    map[configdomain.Key]string{},
		Unscoped: map[configdomain.Key]string{},
	}
}
