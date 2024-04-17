package list

type Status int

const (
	StatusActive  Status = iota // the user is currently entering data into the dialog
	StatusDone                  // the user has made a selection
	StatusAborted               // the user has aborted the dialog
)
