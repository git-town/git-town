package commands

import (
	"fmt"

	"github.com/git-town/git-town/v8/src/config"
)

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func DeleteMainBranchConfiguration(shell Shell) error {
	_, err := shell.Run("git", "config", "--unset", config.MainBranchKey)
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w", err)
	}
	return nil
}
