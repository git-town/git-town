package undo

import (
	"github.com/git-town/git-town/v9/src/config"
)

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Cwd       string // the current working directory
	GitConfig config.GitConfig
}

func (cs ConfigSnapshot) Diff(other ConfigSnapshot) ConfigDiffs {
	return ConfigDiffs{
		Global: NewConfigDiff(cs.GitConfig.Global, other.GitConfig.Global),
		Local:  NewConfigDiff(cs.GitConfig.Local, other.GitConfig.Local),
	}
}
