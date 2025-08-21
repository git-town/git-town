package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	stashTitle = `Stash uncommitted changes when creating branches`
	stashHelp  = `
Should the "hack", "append", and "prepend" commands
stash away uncommitted changes
before creating the new branch and switching to it?

Stashing uncommitted changes makes these commands work reliably,
even if the uncommitted changes conflict with other branches.
But stashing and unstashing changes your Git index.
So if you don't want your staged changes get messed up
when creating new branches, you could disable this setting
to keep your index as-is.

You can also enable or disable stashing using the
--stash and --no-stash CLI flags.
`
)

func Stash(args Args[configdomain.Stash]) (Option[configdomain.Stash], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.Stash]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.Stash]]{
			Data: None[configdomain.Stash](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.Stash]]{
		{
			Data: Some(configdomain.Stash(true)),
			Text: "yes, stash uncommitted changes before creating new branches",
		},
		{
			Data: Some(configdomain.Stash(false)),
			Text: "no, keep my Git index intact",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, stashTitle, stashHelp, args.Inputs, "stash")
	fmt.Printf(messages.StashResult, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
