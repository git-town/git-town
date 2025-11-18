package dialogdomain

// Exit indicates that the user has aborted a dialog
// and wishes to exit the entire Git Town command
// that displayed the dialog.
type Exit bool

func (e Exit) ShouldExit() bool {
	return bool(e)
}
