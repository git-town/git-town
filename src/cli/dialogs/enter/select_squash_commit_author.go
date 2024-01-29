package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

// SelectSquashCommitAuthor allows the user to select an author amongst a given list of authors.
func SelectSquashCommitAuthor(branch gitdomain.LocalBranchName, authors []string, dialogTestInputs components.TestInput) (string, bool, error) {
	if len(authors) == 1 {
		return authors[0], false, nil
	}
	authorsList := squashCommitAuthorList(authors)
	selection, aborted, err := components.RadioList(authorsList, 0, fmt.Sprintf(messages.BranchAuthorMultiple, branch), dialogTestInputs)
	fmt.Printf("Selected squash commit author: %s\n", components.FormattedSelection(selection.String(), aborted))
	return selection.String(), aborted, err
}

type squashCommitAuthor string

func (self squashCommitAuthor) String() string {
	return string(self)
}

func squashCommitAuthorList(authors []string) []squashCommitAuthor {
	result := make([]squashCommitAuthor, len(authors))
	for a, author := range authors {
		result[a] = squashCommitAuthor(author)
	}
	return result
}
