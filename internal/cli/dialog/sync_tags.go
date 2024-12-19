package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/cli/format"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	syncTagsTitle = `Sync-tags strategy`
	SyncTagsHelp  = `
Should "git town sync" sync tags with origin?

`
)

func SyncTags(existing configdomain.SyncTags, inputs components.TestInput) (configdomain.SyncTags, bool, error) {
	entries := list.Entries[configdomain.SyncTags]{
		{
			Data:    true,
			Enabled: true,
			Text:    "yes, sync Git tags",
		},
		{
			Data:    false,
			Enabled: true,
			Text:    "no, don't sync Git tags",
		},
	}
	defaultPos := list.DialogPosition(entries, existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, syncTagsTitle, SyncTagsHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.SyncTags, components.FormattedSelection(format.Bool(selection.IsTrue()), aborted))
	return selection, aborted, err
}
