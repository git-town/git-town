package enter

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterShipDeleteTrackingBranchHelp = `
Should "git ship" delete the tracking branch?

You want to disable this if your code hosting system
(GitHub, GitLab, etc) deletes head branches when
merging pull requests through its UI.

`

const (
	ShipDeleteTrackingBranchEntryYes shipDeleteTrackingBranchEntry = `yes, "git ship" should delete tracking branches`
	ShipDeleteTrackingBranchEntryNo  shipDeleteTrackingBranchEntry = `no, my code hosting platform deletes tracking branches`
)

func EnterShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs dialogcomponents.TestInput) (configdomain.ShipDeleteTrackingBranch, bool, error) {
	entries := []shipDeleteTrackingBranchEntry{
		ShipDeleteTrackingBranchEntryYes,
		ShipDeleteTrackingBranchEntryNo,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := dialogcomponents.RadioList(entries, defaultPos, enterShipDeleteTrackingBranchHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Ship deletes tracking branches: %s\n", dialogcomponents.FormattedSelection(selection.Short(), aborted))
	return selection.ShipDeleteTrackingBranch(), aborted, err
}

type shipDeleteTrackingBranchEntry string

func (self shipDeleteTrackingBranchEntry) ShipDeleteTrackingBranch() configdomain.ShipDeleteTrackingBranch {
	switch self {
	case ShipDeleteTrackingBranchEntryYes:
		return configdomain.ShipDeleteTrackingBranch(true)
	case ShipDeleteTrackingBranchEntryNo:
		return configdomain.ShipDeleteTrackingBranch(false)
	}
	panic("unhandled shipDeleteTrackingBranchEntry: " + self)
}

func (self shipDeleteTrackingBranchEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self shipDeleteTrackingBranchEntry) String() string {
	return string(self)
}
