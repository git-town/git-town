package commands

import (
	"fmt"
	"strings"
)

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func FilesInCommit(repo Repo, sha string) ([]string, error) {
	output, err := repo.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return []string{}, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(output, "\n"), nil
}
