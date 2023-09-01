// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

type querier interface {
	Query(executable string, args ...string) (string, error)
}

type querierRunner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}
