package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

const squashCommitAuthorTitle = `Squash commit author`

// SquashCommitAuthor allows the user to select an author amongst a given list of authors.
func SquashCommitAuthor(branch gitdomain.LocalBranchName, authors []gitdomain.Author, inputs dialogcomponents.Inputs) (gitdomain.Author, dialogdomain.Exit, error) {
	if len(authors) == 1 {
		return authors[0], false, nil
	}
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(authors...), 0, squashCommitAuthorTitle, fmt.Sprintf(messages.BranchAuthorMultiple, branch), inputs, "squash-commit-author")
	fmt.Printf(messages.SquashCommitAuthorSelection, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
