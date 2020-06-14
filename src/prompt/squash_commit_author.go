package prompt

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/run"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func GetSquashCommitAuthor(branchName string, repo *git.ProdRepo) (string, error) {
	authors, err := getBranchAuthors(branchName, repo)
	if err != nil {
		return "", err
	}
	if len(authors) == 1 {
		return authors[0], nil
	}
	cli.Printf(squashCommitAuthorHeaderTemplate, branchName)
	fmt.Println()
	return askForAuthor(authors), nil
}

// Helpers

var squashCommitAuthorHeaderTemplate = "Multiple people authored the %q branch."

func askForAuthor(authors []string) string {
	result := ""
	prompt := &survey.Select{
		Message: "Please choose an author for the squash commit:",
		Options: authors,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		panic(err)
	}
	return result
}

func getBranchAuthors(branchName string, repo *git.ProdRepo) (result []string, err error) {
	// Returns lines of "<number of commits>\t<name and email>"
	lines, err := run.Exec("git", "shortlog", "-s", "-n", "-e", repo.Config.GetMainBranch()+".."+branchName)
	if err != nil {
		return result, err
	}
	for _, line := range lines.OutputLines() {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return result, nil
}
