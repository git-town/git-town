package commands

import (
	"fmt"
	"strings"

	prodgit "github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/stringslice"
	"github.com/git-town/git-town/v8/test/subshell"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	subshell.Mocking
	Config prodgit.RepoConfig
	*prodgit.BackendCommands
}

// StageFiles adds the file with the given name to the Git index.
func (r *TestCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := r.Run("git", args...)
	return err
}

// StashSize provides the number of stashes in this repository.
func (r *TestCommands) StashSize() (int, error) {
	output, err := r.Run("git", "stash", "list")
	if err != nil {
		return 0, fmt.Errorf("cannot determine Git stash: %w", err)
	}
	if output == "" {
		return 0, nil
	}
	return len(stringslice.Lines(output)), nil
}

// Tags provides a list of the tags in this repository.
func (r *TestCommands) Tags() ([]string, error) {
	output, err := r.Run("git", "tag")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine tags in repo %q: %w", r.WorkingDir, err)
	}
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// UncommittedFiles provides the names of the files not committed into Git.
func (r *TestCommands) UncommittedFiles() ([]string, error) {
	output, err := r.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine uncommitted files in %q: %w", r.WorkingDir, err)
	}
	result := []string{}
	for _, line := range stringslice.Lines(output) {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, parts[1])
	}
	return result, nil
}
