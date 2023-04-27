package commands

import (
	"fmt"

	"github.com/git-town/git-town/v8/src/stringslice"
)

// StashSize provides the number of stashes in this repository.
func StashSize(repo *Repo) (int, error) {
	output, err := repo.Run("git", "stash", "list")
	if err != nil {
		return 0, fmt.Errorf("cannot determine Git stash: %w", err)
	}
	if output == "" {
		return 0, nil
	}
	return len(stringslice.Lines(output)), nil
}
