package undoconfig

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Global configdomain.SingleSnapshot
	Local  configdomain.SingleSnapshot
}

func EmptyConfigSnapshot() ConfigSnapshot {
	return ConfigSnapshot{
		Global: map[configdomain.Key]string{},
		Local:  map[configdomain.Key]string{},
	}
}
