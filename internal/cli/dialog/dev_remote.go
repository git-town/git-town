package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	DevRemoteTypeTitle = `Development Remote`
	DevRemoteHelp      = `
Which remote should Git Town use
for development?

Typically that's the "origin" remote.

`
)

func DevRemote(existingValue gitdomain.Remote, options gitdomain.Remotes, inputs dialogcomponents.TestInput) (gitdomain.Remote, dialogdomain.Exit, error) {
	cursor := slice.Index(options, existingValue).GetOrElse(0)
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(options...), cursor, DevRemoteTypeTitle, DevRemoteHelp, inputs)
	fmt.Printf(messages.DevRemote, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
