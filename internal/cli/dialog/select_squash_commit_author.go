package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const squashCommitAuthorTitle = `Squash commit author`

// SelectSquashCommitAuthor allows the user to select an author amongst a given list of authors.
func SelectSquashCommitAuthor(branch gitdomain.LocalBranchName, authors []gitdomain.Author, dialogTestInputs dialogcomponents.TestInput) (gitdomain.Author, dialogdomain.Exit, error) {
	if len(authors) == 1 {
		return authors[0], false, nil
	}
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(authors...), 0, squashCommitAuthorTitle, fmt.Sprintf(messages.BranchAuthorMultiple, branch), dialogTestInputs)
	fmt.Printf(messages.SquashCommitAuthorSelection, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
