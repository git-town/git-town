package undoconfig

import (
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/pkg/keys"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Global configdomain.SingleSnapshot
	Local  configdomain.SingleSnapshot
}

func EmptyConfigSnapshot() ConfigSnapshot {
	return ConfigSnapshot{
		Global: map[keys.Key]string{},
		Local:  map[keys.Key]string{},
	}
}
