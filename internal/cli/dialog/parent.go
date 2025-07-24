package dialog

// ParentOutcome describes the selection that the user made in the `Parent` dialog.
type ParentOutcome int

const (
	ParentOutcomeExit            ParentOutcome = iota // the user exited the dialog
	ParentOutcomePerennialBranch                      // the user chose the "perennial branch" option
	ParentOutcomeSelectedParent                       // the user selected one of the branches
)
