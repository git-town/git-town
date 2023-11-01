package undo

import (
	"github.com/git-town/git-town/v10/src/config"
)

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Cwd       string // the current working directory
	GitConfig config.GitConfig
}
