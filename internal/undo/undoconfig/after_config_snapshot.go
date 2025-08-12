package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type AfterConfigSnapshot struct {
	Global configdomain.SingleSnapshot
	Local  configdomain.SingleSnapshot
}

func EmptyAfterConfigSnapshot() AfterConfigSnapshot {
	return AfterConfigSnapshot{
		Global: map[configdomain.Key]string{},
		Local:  map[configdomain.Key]string{},
	}
}
