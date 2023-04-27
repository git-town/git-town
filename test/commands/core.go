package commands

import (
	prodgit "github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/test/subshell"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	subshell.Mocking
	Config prodgit.RepoConfig
	*prodgit.BackendCommands
}

type Shell interface {
	Run(string, ...string) (string, error)
	RunMany([][]string) error
	Dir() string
}
