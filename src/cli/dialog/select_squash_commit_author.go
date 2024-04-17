package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

const squashCommitAuthorTitle = `Squash commit author`

// SelectSquashCommitAuthor allows the user to select an author amongst a given list of authors.
func SelectSquashCommitAuthor(branch gitdomain.LocalBranchName, authors []string, dialogTestInputs components.TestInput) (string, bool, error) {
	if len(authors) == 1 {
		return authors[0], false, nil
	}
	authorsList := squashCommitAuthorList(authors)
	selection, aborted, err := components.RadioList(list.NewEnabledListEntries(authorsList), 0, squashCommitAuthorTitle, fmt.Sprintf(messages.BranchAuthorMultiple, branch), dialogTestInputs)
	fmt.Printf(messages.SquashCommitAuthorSelection, components.FormattedSelection(selection.String(), aborted))
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
