package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
)

const (
	syncTagsTitle = `Sync-tags strategy`
	SyncTagsHelp  = `
Should "git sync" sync tags with origin?

`
)

const (
	SyncTagsEntryYes syncTagsEntry = `yes, sync Git tags`
	SyncTagsEntryNo  syncTagsEntry = `no, don't sync Git tags`
)

func SyncTags(existing configdomain.SyncTags, inputs components.TestInput) (configdomain.SyncTags, bool, error) {
	entries := list.NewEntries(
		SyncTagsEntryYes,
		SyncTagsEntryNo,
	)
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, syncTagsTitle, SyncTagsHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.SyncTags, components.FormattedSelection(selection.Data.Short(), aborted))
	return selection.Data.SyncTags(), aborted, err
}

type syncTagsEntry string

func (self syncTagsEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self syncTagsEntry) String() string {
	return string(self)
}

func (self syncTagsEntry) SyncTags() configdomain.SyncTags {
	switch self {
	case SyncTagsEntryYes:
		return configdomain.SyncTags(true)
	case SyncTagsEntryNo:
		return configdomain.SyncTags(false)
	}
	panic("unhandled syncTagsEntry: " + self)
}
