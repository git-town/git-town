package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterSyncBeforeShipHelp = `
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
	syncBeforeShipEntryYes syncBeforeShipEntry = `yes, "git ship" should also sync the branch`
	syncBeforeShipEntryNo  syncBeforeShipEntry = `no, "git ship" should not sync the branch`
)

func EnterSyncBeforeShip(existing configdomain.SyncBeforeShip, inputs TestInput) (configdomain.SyncBeforeShip, bool, error) {
	entries := []syncBeforeShipEntry{syncBeforeShipEntryYes, syncBeforeShipEntryNo}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(entries, defaultPos, enterSyncBeforeShipHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection.String(), ",")
	fmt.Printf("Sync before ship: %s\n", formattedSelection(cutSelection, aborted))
	return selection.ToSyncBeforeShip(), aborted, err
}

type syncBeforeShipEntry string

func (self syncBeforeShipEntry) String() string {
	return string(self)
}

func (self syncBeforeShipEntry) ToSyncBeforeShip() configdomain.SyncBeforeShip {
	switch self {
	case syncBeforeShipEntryYes:
		return configdomain.SyncBeforeShip(true)
	case syncBeforeShipEntryNo:
		return configdomain.SyncBeforeShip(false)
	}
	panic("unhandled syncBeforeShipEntry: " + self)
}
