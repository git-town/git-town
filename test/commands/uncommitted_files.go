package commands

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v8/src/stringslice"
)

// UncommittedFiles provides the names of the files not committed into Git.
func UncommittedFiles(repo *Repo) ([]string, error) {
	output, err := repo.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine uncommitted files in %q: %w", repo.Dir(), err)
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
