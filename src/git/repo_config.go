package git

import (
	"github.com/git-town/git-town/v11/src/config"
)

// RepoConfig represents the known state of a Git repository.
type RepoConfig struct {
	*config.GitTown
}
