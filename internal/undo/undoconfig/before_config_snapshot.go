package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type BeforeConfigSnapshot struct {
	Global   configdomain.SingleSnapshot
	Local    configdomain.SingleSnapshot
	Unscoped configdomain.SingleSnapshot
}

func EmptyBeforeConfigSnapshot() BeforeConfigSnapshot {
	return BeforeConfigSnapshot{
		Global:   map[configdomain.Key]string{},
		Local:    map[configdomain.Key]string{},
		Unscoped: map[configdomain.Key]string{},
	}
}
