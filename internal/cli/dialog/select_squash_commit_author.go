package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const squashCommitAuthorTitle = `Squash commit author`

// SelectSquashCommitAuthor allows the user to select an author amongst a given list of authors.
func SelectSquashCommitAuthor(branch gitdomain.LocalBranchName, authors []gitdomain.Author, dialogTestInputs components.TestInput) (gitdomain.Author, dialogdomain.Aborted, error) {
	if len(authors) == 1 {
		return authors[0], false, nil
	}
	selection, aborted, err := components.RadioList(list.NewEntries(authors...), 0, squashCommitAuthorTitle, fmt.Sprintf(messages.BranchAuthorMultiple, branch), dialogTestInputs)
	fmt.Printf(messages.SquashCommitAuthorSelection, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
