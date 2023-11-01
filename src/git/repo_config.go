package git

import (
	"github.com/git-town/git-town/v10/src/config"
)

// RepoConfig represents the known state of a Git repository.
type RepoConfig struct {
	*config.GitTown
	DryRun bool // single source of truth for whether to dry-run Git commands in this repo
}
