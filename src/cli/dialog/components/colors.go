package components

import "github.com/muesli/termenv"

// Typical colors used in BubbleTea dialogs.
type dialogColors struct {
	Help      termenv.Style // color of help text
	HelpKey   termenv.Style // color of key names in help text
	Initial   termenv.Style // color for the row containing the currently checked out branch
	Selection termenv.Style // color for the currently selected entry
	Title     termenv.Style // color for the title of the current screen
}

// FormattedToken provides the given API token in a printable format.
func FormattedSecret(secret string, aborted bool) string {
	if aborted {
		return Red().Styled("(aborted)")
	}
	if secret == "" {
		return green().Styled("(not provided)")
	}
	return green().Styled("(provided)")
}

// FormattedSelection provides the given dialog choice in a printable format.
func FormattedSelection(selection string, aborted bool) string {
	if aborted {
		return Red().Styled("(aborted)")
	}
	return green().Styled(selection)
}

// FormattedToken provides the given API token in a printable format.
func FormattedToken(token string, aborted bool) string {
	if aborted {
		return Red().Styled("(aborted)")
	}
	if token == "" {
		return green().Styled("(not provided)")
	}
	return green().Styled(token)
}

func createColors() dialogColors {
	return dialogColors{
		Help:      termenv.String().Faint(),
		HelpKey:   termenv.String().Faint().Bold(),
		Initial:   termenv.String().Foreground(termenv.ANSIGreen),
		Selection: termenv.String().Foreground(termenv.ANSICyan),
		Title:     termenv.String().Bold(),
	}
}

func green() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIGreen)
}

func Red() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIRed)
}
