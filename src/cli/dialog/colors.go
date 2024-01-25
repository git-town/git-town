package dialog

import "github.com/muesli/termenv"

// Typical colors used in BubbleTea dialogs.
type dialogColors struct {
	help      termenv.Style // color of help text
	helpKey   termenv.Style // color of key names in help text
	initial   termenv.Style // color for the row containing the currently checked out branch
	selection termenv.Style // color for the currently selected entry
}

func createColors() dialogColors {
	return dialogColors{
		help:      termenv.String().Faint(),
		helpKey:   termenv.String().Faint().Bold(),
		initial:   termenv.String().Foreground(termenv.ANSIGreen),
		selection: termenv.String().Foreground(termenv.ANSICyan),
	}
}

func formattedSelection(selection string, aborted bool) string {
	if aborted {
		return red().Styled("(aborted)")
	}
	return green().Styled(selection)
}

func green() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIGreen)
}

func red() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIRed)
}
