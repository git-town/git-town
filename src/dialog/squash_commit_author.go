package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// DetermineSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func DetermineSquashCommitAuthor(branch string, authors []string, repo *git.ProdRepo) (string, error) {
	if len(authors) == 1 {
		return authors[0], nil
	}
	cli.Printf(squashCommitAuthorHeaderTemplate, branch)
	fmt.Println()
	return askForAuthor(authors)
}

// Helpers

const squashCommitAuthorHeaderTemplate = "Multiple people authored the %q branch."

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
