package dialog

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	if len(remotes) == 0 {
		return None[gitdomain.Remote](), false, nil
	}
	options := list.Entries[Option[gitdomain.Remote]]{}
	global, hasGlobal := args.Global.Get()
	if hasGlobal {
		options = append(options, list.Entry[Option[gitdomain.Remote]]{
			Data: None[gitdomain.Remote](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
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
		return None[gitdomain.Remote](), false, nil
	}
	var cursor int
	switch {
	case hasLocal:
		cursor = options.IndexOf(Some(local))
	case hasGlobal:
		cursor = 0
	default:
		cursor = options.IndexOf(Some(config.DefaultNormalConfig().DevRemote))
	}
	selection, exit, err := dialogcomponents.RadioList(options, cursor, DevRemoteTypeTitle, DevRemoteHelp, args.Inputs, "dev-remote")
	fmt.Printf(messages.DevRemote, dialogcomponents.FormattedSelection(selection.GetOrZero().String(), exit))
	return selection, exit, err
}
