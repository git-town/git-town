package repo

import (
	"fmt"
	"strings"
)

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func FileContentInCommit(repo *Repo, sha string, filename string) (string, error) {
	output, err := repo.Run("git", "show", sha+":"+filename)
	if err != nil {
		return "", fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	result := output
	if strings.HasPrefix(result, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		result = ""
	}
	return result, nil
}
