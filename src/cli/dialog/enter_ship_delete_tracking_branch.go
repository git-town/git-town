package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterShipDeleteTrackingBranchHelp = `
Should "git ship" delete the tracking branch?

You want to disable this if your code hosting system
(GitHub, GitLab, etc) deletes head branches when
merging pull requests through its UI.

`

func EnterShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs TestInput) (configdomain.ShipDeleteTrackingBranch, bool, error) {
	entries := []string{`yes, "git ship" should delete tracking branches`, `no, my code hosting platform deletes tracking branches`}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterShipDeleteTrackingBranchHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Ship deletes tracking branches: %s\n", formattedSelection(selection, aborted))
	cutSelection, _, _ := strings.Cut(selection, ",")
	parsedAnswer, err := configdomain.ParseShipDeleteTrackingBranch(cutSelection, "user dialog")
	return parsedAnswer, aborted, err
}
