package dialogdomain

// Exit indicates that the user has aborted a dialog
// and wishes to exit the entire Git Town command
// that displayed the dialog.
type Exit bool

func (self Exit) ShouldExit() bool {
	return bool(self)
}
