package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/run"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// DetermineSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func DetermineSquashCommitAuthor(branchName string, repo *git.ProdRepo) (string, error) {
	authors, err := loadBranchAuthors(branchName, repo)
	if err != nil {
		return "", err
	}
	if len(authors) == 1 {
		return authors[0], nil
	}
	cli.Printf(squashCommitAuthorHeaderTemplate, branchName)
	fmt.Println()
	return askForAuthor(authors)
}

// Helpers

var squashCommitAuthorHeaderTemplate = "Multiple people authored the %q branch."

func askForAuthor(authors []string) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: "Please choose an author for the squash commit:",
		Options: authors,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read author from CLI: %w", err)
	}
	return result, nil
}

func loadBranchAuthors(branchName string, repo *git.ProdRepo) ([]string, error) {
	// Returns lines of "<number of commits>\t<name and email>"
	lines, err := run.Exec("git", "shortlog", "-s", "-n", "-e", repo.Config.MainBranch()+".."+branchName)
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range lines.OutputLines() {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return result, nil
}
