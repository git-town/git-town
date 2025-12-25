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
	syncTagsTitle = `Sync-tags strategy`
	SyncTagsHelp  = `
Should "git town sync" sync Git tags with origin?

`
)

func SyncTags(args Args[configdomain.SyncTags]) (Option[configdomain.SyncTags], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncTags]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.SyncTags]]{
			Data: None[configdomain.SyncTags](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.SyncTags]]{
		{
			Data: Some(configdomain.SyncTags(true)),
			Text: "yes, sync Git tags",
		},
		{
			Data: Some(configdomain.SyncTags(false)),
			Text: "no, don't sync Git tags",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncTagsTitle, SyncTagsHelp, args.Inputs, "sync-tags")
	fmt.Printf(messages.SyncTags, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
