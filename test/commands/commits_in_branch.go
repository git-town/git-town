package commands

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v8/src/stringslice"
	"github.com/git-town/git-town/v8/test/git"
)

// CommitsInBranch provides all commits in the given Git branch.
func CommitsInBranch(repo Repo, branch string, fields []string) ([]git.Commit, error) {
	output, err := repo.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return []git.Commit{}, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if stringslice.Contains(fields, "FILE NAME") {
			filenames, err := FilesInCommit(repo, commit.SHA)
			if err != nil {
				return []git.Commit{}, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if stringslice.Contains(fields, "FILE CONTENT") {
			filecontent, err := FileContentInCommit(repo, commit.SHA, commit.FileName)
			if err != nil {
				return []git.Commit{}, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}
