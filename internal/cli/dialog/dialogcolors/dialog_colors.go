package dialogcolors

import (
	"github.com/git-town/git-town/v22/pkg/colors"
	"github.com/muesli/termenv"
)

// DialogColors defines the colors used in dialogs.
type DialogColors struct {
	EntryNumber termenv.Style // color for the number of entries
	Help        termenv.Style // color of help text
	HelpKey     termenv.Style // color of key names in help text
	Initial     termenv.Style // color for the row containing the currently checked out branch
	Selection   termenv.Style // color for the currently selected entry
	Title       termenv.Style // color for the title of the current screen
}

func NewDialogColors() DialogColors {
	return DialogColors{
		EntryNumber: colors.Faint(),
		Help:        colors.Faint(),
		HelpKey:     colors.FaintBold(),
		Initial:     colors.Green(),
		Selection:   colors.Cyan(),
		Title:       colors.Bold(),
	}
}
