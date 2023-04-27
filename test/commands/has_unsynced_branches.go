package commands

import (
	"fmt"
	"strings"
)

// HasUnsyncedBranches indicates whether one or more local branches are out of sync with their tracking branch.
func HasUnsyncedBranches(shell Shell) (bool, error) {
	output, err := shell.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", shell.Dir(), err, output)
	}
	return strings.Contains(output, "["), nil
}
