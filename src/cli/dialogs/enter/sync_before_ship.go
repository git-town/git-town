package enter

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const syncBeforeShipHelp = `
Should "git ship" sync branches before shipping them?

Guidance: enable when shipping branches locally on your machine
and disable when shipping feature branches via the code hosting
API or web UI.

When enabled, branches are always fully up to date when shipped
and you get a chance to resolve merge conflicts
between the feature branch to ship and the main development branch
on the feature branch. This helps keep the main branch green.
But this also triggers another CI run and delays shipping.

`

const (
	SyncBeforeShipEntryYes syncBeforeShipEntry = `yes, "git ship" should also sync the branch`
	SyncBeforeShipEntryNo  syncBeforeShipEntry = `no, "git ship" should not sync the branch`
)

func SyncBeforeShip(existing configdomain.SyncBeforeShip, inputs dialog.TestInput) (configdomain.SyncBeforeShip, bool, error) {
	entries := []syncBeforeShipEntry{
		SyncBeforeShipEntryYes,
		SyncBeforeShipEntryNo,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := dialog.RadioList(entries, defaultPos, syncBeforeShipHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Sync before ship: %s\n", dialog.FormattedSelection(selection.Short(), aborted))
	return selection.SyncBeforeShip(), aborted, err
}

type syncBeforeShipEntry string

func (self syncBeforeShipEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self syncBeforeShipEntry) String() string {
	return string(self)
}

func (self syncBeforeShipEntry) SyncBeforeShip() configdomain.SyncBeforeShip {
	switch self {
	case SyncBeforeShipEntryYes:
		return configdomain.SyncBeforeShip(true)
	case SyncBeforeShipEntryNo:
		return configdomain.SyncBeforeShip(false)
	}
	panic("unhandled syncBeforeShipEntry: " + self)
}
