package commands

import (
	"fmt"

	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/git"
)

// CreateCommit creates a commit with the given properties in this Git repo.
func CreateCommit(repo *Repo, commit git.Commit) error {
	err := repo.ProdGit().CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = fs.CreateFile(repo.Dir(), commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	_, err = repo.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w", err)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	_, err = repo.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w", err)
	}
	return nil
}
