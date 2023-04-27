package commands

import (
	"fmt"
	"strings"
)

// ShaForCommit provides the SHA for the commit with the given name.
func ShaForCommit(repo *Repo, name string) (string, error) {
	output, err := repo.Run("git", "log", "--reflog", "--format=%H", "--grep=^"+name+"$")
	if err != nil {
		return "", fmt.Errorf("cannot determine the SHA of commit %q: %w", name, err)
	}
	result := output
	if result == "" {
		return "", fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	result = strings.Split(result, "\n")[0]
	return result, nil
}
