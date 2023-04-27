package datatable

import "github.com/git-town/git-town/v8/src/git"

type Shell interface {
	Run(string, ...string) (string, error)
	RunMany([][]string) error
	Dir() string

	// TODO: clean this up after extracting production git commands into stand-alone functions
	Config() *git.RepoConfig
	ProdGit() *git.BackendCommands
}
