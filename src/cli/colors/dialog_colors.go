package colors

import "github.com/muesli/termenv"

// Typical colors used in BubbleTea dialogs.
type DialogColors struct {
	Help      termenv.Style // color of help text
	HelpKey   termenv.Style // color of key names in help text
	Initial   termenv.Style // color for the row containing the currently checked out branch
	Selection termenv.Style // color for the currently selected entry
	Title     termenv.Style // color for the title of the current screen
}

func NewDialogColors() DialogColors {
	return DialogColors{
		Help:      Faint(),
		HelpKey:   FaintBold(),
		Initial:   Green(),
		Selection: Cyan(),
		Title:     Bold(),
	}
}
