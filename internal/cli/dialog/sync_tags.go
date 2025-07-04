package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	syncTagsTitle = `Sync-tags strategy`
	SyncTagsHelp  = `
Should "git town sync" sync Git tags with origin?

`
)

func SyncTags(existing configdomain.SyncTags, inputs dialogcomponents.TestInput) (configdomain.SyncTags, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.SyncTags]{
		{
			Data: true,
			Text: "yes, sync Git tags",
		},
		{
			Data: false,
			Text: "no, don't sync Git tags",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncTagsTitle, SyncTagsHelp, inputs)
	fmt.Printf(messages.SyncTags, dialogcomponents.FormattedSelection(format.Bool(selection.IsTrue()), exit))
	return selection, exit, err
}
