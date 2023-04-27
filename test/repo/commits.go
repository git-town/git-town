package repo

import (
	"fmt"

	"github.com/git-town/git-town/v8/test/git"
)

// Commits provides a list of the commits in this Git repository with the given fields.
func Commits(repo *Repo, fields []string, mainBranch string) ([]git.Commit, error) {
	branches, err := repo.ProdGit().LocalBranchesMainFirst(mainBranch)
	if err != nil {
		return []git.Commit{}, fmt.Errorf("cannot determine the Git branches: %w", err)
	}
	result := []git.Commit{}
	for _, branch := range branches {
		commits, err := CommitsInBranch(repo, branch, fields)
		if err != nil {
			return []git.Commit{}, err
		}
		result = append(result, commits...)
	}
	return result, nil
}
