package commands

import (
	"github.com/git-town/git-town/v8/src/git"
)

// Repo is a repository clone on which the test commands execute.
type Repo interface {
	Run(string, ...string) (string, error)
	RunMany([][]string) error
	Dir() string

	// TODO: clean this up after extracting production git commands into stand-alone functions
	Config() *git.RepoConfig
	ProdGit() *git.BackendCommands
}
