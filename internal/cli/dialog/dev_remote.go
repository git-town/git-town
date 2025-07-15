package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	DevRemoteTypeTitle = `Development Remote`
	DevRemoteHelp      = `
Which remote should Git Town use
for development?

Typically that's the "origin" remote.

`
)

func DevRemote(existingValue gitdomain.Remote, remotes gitdomain.Remotes, inputs dialogcomponents.TestInputs) (Option[gitdomain.Remote], dialogdomain.Exit, error) {
	cursor := slice.Index(remotes, existingValue).GetOrElse(0)
	entries := make(list.Entries[Option[gitdomain.Remote]], len(remotes)+1)
	entries[0] = list.Entry[Option[gitdomain.Remote]]{
		Data:     None[gitdomain.Remote](),
		Disabled: false,
		Text:     messages.DialogDefaultText,
	}
	for _, remote := range remotes {
		entries = append(entries, list.Entry[Option[gitdomain.Remote]]{
			Data:     Some(remote),
			Disabled: false,
			Text:     remote.String(),
		})
	}
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, DevRemoteTypeTitle, DevRemoteHelp, inputs)
	fmt.Printf(messages.DevRemote, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
