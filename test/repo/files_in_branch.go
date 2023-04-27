package repo

import (
	"fmt"
	"strings"
)

// FilesInBranch provides the list of the files present in the given branch.
func FilesInBranch(repo *Repo, branch string) ([]string, error) {
	output, err := repo.Run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine files in branch %q in repo %q: %w", branch, repo.Dir(), err)
	}
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result, err
}
