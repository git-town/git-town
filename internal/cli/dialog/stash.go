package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	stashTitle = `Stash uncommitted changes when creating branches`
	stashHelp  = `
Should the hack, append, and prepend commands
automatically stash uncommitted changes
before creating and switching to a new branch?

Stashing ensures these commands work reliably,
even if the uncommitted changes conflict
with the new branch. The downside is that
stashing and unstashing alters your Git index.

If you'd rather keep your index untouched,
at the cost of potentially dealing with more merge conflicts,
you can disable stashing.

You can also override the default on a case-by-case basis
with the --stash and --no-stash CLI flags.

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
