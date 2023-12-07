package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/messages"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// SelectAuthor allows the user to select an author amongst a given list of authors.
func SelectAuthor(branch string, authors []string) (string, error) {
	if len(authors) == 1 {
		return authors[0], nil
	}
	io.Printf("Multiple people authored the %q branch.", branch)
	fmt.Println()
	result := ""
	prompt := &survey.Select{
		Message: "Please choose an author for the squash commit:",
		Options: authors,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf(messages.DialogCannotReadAuthor, err)
	}
	return result, nil
}
