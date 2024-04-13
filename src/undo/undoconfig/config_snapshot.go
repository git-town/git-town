package undoconfig

import "github.com/git-town/git-town/v14/src/config/gitconfig"

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Global gitconfig.SingleSnapshot
	Local  gitconfig.SingleSnapshot
}

func EmptyConfigSnapshot() ConfigSnapshot {
	return ConfigSnapshot{
		Global: map[gitconfig.Key]string{},
		Local:  map[gitconfig.Key]string{},
	}
}
