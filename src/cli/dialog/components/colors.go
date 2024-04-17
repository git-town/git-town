package components

import (
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/muesli/termenv"
)

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
		return colors.Red().Styled("(aborted)")
	}
	if secret == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled("(provided)")
}

// FormattedSelection provides the given dialog choice in a printable format.
func FormattedSelection(selection string, aborted bool) string {
	if aborted {
		return colors.Red().Styled("(aborted)")
	}
	return colors.Green().Styled(selection)
}

// FormattedToken provides the given API token in a printable format.
func FormattedToken(token string, aborted bool) string {
	if aborted {
		return colors.Red().Styled("(aborted)")
	}
	if token == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled(token)
}
