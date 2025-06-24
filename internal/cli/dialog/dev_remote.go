package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	DevRemoteTypeTitle = `Development Remote`
	DevRemoteHelp      = `
Which remote should Git Town use for development?

Typically that's the "origin" remote.

`
)

func DevRemote(existingValue gitdomain.Remote, options gitdomain.Remotes, inputs components.TestInput) (gitdomain.Remote, dialogdomain.Exit, error) {
	cursor := slice.Index(options, existingValue).GetOrElse(0)
	selection, aborted, err := components.RadioList(list.NewEntries(options...), cursor, DevRemoteTypeTitle, DevRemoteHelp, inputs)
	fmt.Printf(messages.DevRemote, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
