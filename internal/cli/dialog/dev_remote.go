package dialog

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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

func DevRemote(remotes gitdomain.Remotes, args Args[gitdomain.Remote]) (Option[gitdomain.Remote], dialogdomain.Exit, error) {
	options := list.Entries[Option[gitdomain.Remote]]{}
	global, hasGlobal := args.Global.Get()
	if hasGlobal {
		options = append(options, list.Entry[Option[gitdomain.Remote]]{
			Data: None[gitdomain.Remote](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	} else {
		options = append(options, list.Entry[Option[gitdomain.Remote]]{
			Data: None[gitdomain.Remote](),
			Text: fmt.Sprintf(messages.DialogUseDefaultValue, args.Defaults),
		})
	}
	for _, remote := range remotes {
		options = append(options, list.Entry[Option[gitdomain.Remote]]{
			Data: Some(remote),
			Text: remote.String(),
		})
	}
	local, hasLocal := args.Local.Get()
	if hasLocal && !slices.Contains(remotes, local) {
		options = append(options, list.Entry[Option[gitdomain.Remote]]{
			Data: Some(local),
			Text: local.String(),
		})
	}
	if len(options) == 1 {
		return options[0].Data, false, nil
	}
	cursor := 0
	if hasLocal {
		cursor = options.IndexOf(Some(local))
	}
	selection, exit, err := dialogcomponents.RadioList(options, cursor, DevRemoteTypeTitle, DevRemoteHelp, args.Inputs, "dev-remote")
	fmt.Printf(messages.DevRemote, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
