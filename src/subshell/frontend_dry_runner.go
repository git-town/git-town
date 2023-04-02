package subshell

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/kballard/go-shellquote"
)

// FrontendDryRunner prints the given shell commands to the CLI as if they were executed
// but does not execute them.
type FrontendDryRunner struct {
	CurrentBranch   *cache.String
	OmitBranchNames bool
	Stats           Statistics
}

// Run runs the given command in this ShellRunner's directory.
func (r *FrontendDryRunner) Run(executable string, args ...string) error {
	PrintCommand(r.CurrentBranch.Value(), r.OmitBranchNames, executable, args...)
	return nil
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Failed commands abort immediately with the encountered error.
func (r *FrontendDryRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments) in this ShellInDir's directory.
func (r *FrontendDryRunner) RunString(fullCmd string) error {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.Run(cmd, args...)
}
