package prompt

import (
	"fmt"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/git"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func GetSquashCommitAuthor(branchName string) string {
	authors := getBranchAuthors(branchName)
	if len(authors) == 1 {
		return authors[0]
	}
	cfmt.Printf(squashCommitAuthorHeaderTemplate, branchName)
	fmt.Println()
	return askForAuthor(authors)
}

// Helpers

var squashCommitAuthorHeaderTemplate = "Multiple people authored the '%s' branch."

func askForAuthor(authors []string) string {
	result := ""
	prompt := &survey.Select{
		Message: "Please choose an author for the squash commit:",
		Options: authors,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}

func getBranchAuthors(branchName string) (result []string) {
	// Returns lines of "<number of commits>\t<name and email>"
	output := command.New("git", "shortlog", "-s", "-n", "-e", git.GetMainBranch()+".."+branchName).Output()
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return
}
