package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	commitTitleTitle = `Select commit title for proposal`
	commitTitleHelp  = `
Select which commit's title should be used
for the proposal title, or select "none"
to use the default behavior.
`
)

// CommitTitle lets the user select a commit title from the given commits.
func CommitTitle(commits gitdomain.Commits, inputs dialogcomponents.Inputs) (Option[gitdomain.CommitMessage], dialogdomain.Exit, error) {
	if len(commits) == 0 {
		return None[gitdomain.CommitMessage](), false, nil
	}
	if len(commits) == 1 {
		return Some(commits[0].Message), false, nil
	}
	entries := make(list.Entries[Option[gitdomain.CommitMessage]], len(commits)+1)
	entries[0] = list.Entry[Option[gitdomain.CommitMessage]]{
		Data: None[gitdomain.CommitMessage](),
		Text: "(use default)",
	}
	for c, commit := range commits {
		parts := commit.Message.Parts()
		entries[c+1] = list.Entry[Option[gitdomain.CommitMessage]]{
			Data: Some(commit.Message),
			Text: parts.Subject,
		}
	}
	selection, exit, err := dialogcomponents.RadioList(entries, 1, commitTitleTitle, commitTitleHelp, inputs, "commit-title")
	if err != nil || exit {
		return None[gitdomain.CommitMessage](), exit, err
	}
	if selected, has := selection.Get(); has {
		fmt.Printf(messages.CommitTitleSelected, selected.Parts().Subject)
	} else {
		fmt.Printf("Commit title: (none)\n")
	}
	return selection, exit, err
}
