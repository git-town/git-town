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

func EnterSyncBeforeShip(existing configdomain.SyncBeforeShip, inputs TestInput) (configdomain.SyncBeforeShip, bool, error) {
	entries := []string{`yes, "git ship" should also sync the branch`, `no, "git ship" should not sync the branch`}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterSyncBeforeShipHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Sync before ship: %s\n", formattedSelection(selection, aborted))
	cutSelection, _, _ := strings.Cut(selection, ",")
	parsedAnswer, err := configdomain.ParseSyncBeforeShip(cutSelection, "user dialog")
	return parsedAnswer, aborted, err
}
