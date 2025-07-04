package dialogcomponents

import (
	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
)

// FormattedToken provides the given API token in a printable format.
func FormattedSecret(secret string, exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	if secret == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled("(provided)")
}

// FormattedSelection provides the given dialog choice in a printable format.
func FormattedSelection(selection string, exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	return colors.Green().Styled(selection)
}

// FormattedToken provides the given API token in a printable format.
func FormattedToken(token string, exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	if token == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled(token)
}
