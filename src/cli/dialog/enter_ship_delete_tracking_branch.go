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

const (
	shipDeleteTrackingBranchEntryYes shipDeleteTrackingBranchEntry = `yes, "git ship" should delete tracking branches`
	shipDeleteTrackingBranchEntryNo  shipDeleteTrackingBranchEntry = `no, my code hosting platform deletes tracking branches`
)

func EnterShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs TestInput) (configdomain.ShipDeleteTrackingBranch, bool, error) {
	entries := []shipDeleteTrackingBranchEntry{shipDeleteTrackingBranchEntryYes, shipDeleteTrackingBranchEntryNo}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(entries, defaultPos, enterShipDeleteTrackingBranchHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Ship deletes tracking branches: %s\n", formattedSelection(selection.Short(), aborted))
	return selection.ToShipDeleteTrackingBranch(), aborted, err
}

type shipDeleteTrackingBranchEntry string

func (self shipDeleteTrackingBranchEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self shipDeleteTrackingBranchEntry) String() string {
	return string(self)
}

func (self shipDeleteTrackingBranchEntry) ToShipDeleteTrackingBranch() configdomain.ShipDeleteTrackingBranch {
	switch self {
	case shipDeleteTrackingBranchEntryYes:
		return configdomain.ShipDeleteTrackingBranch(true)
	case shipDeleteTrackingBranchEntryNo:
		return configdomain.ShipDeleteTrackingBranch(false)
	}
	panic("unhandled shipDeleteTrackingBranchEntry: " + self)
}
