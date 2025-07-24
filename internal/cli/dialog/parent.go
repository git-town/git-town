package dialog

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)")

const (
	ParentBranchTitleTemplate = `Parent branch for %s`
	parentBranchHelpTemplate  = `
Please select the parent of branch %q
or enter its number.


`
)

// ParentOutcome describes the selection that the user made in the `Parent` dialog.
type ParentOutcome int

const (
	ParentOutcomeExit            ParentOutcome = iota // the user exited the dialog
	ParentOutcomePerennialBranch                      // the user chose the "perennial branch" option
	ParentOutcomeSelectedParent                       // the user selected one of the branches
)
