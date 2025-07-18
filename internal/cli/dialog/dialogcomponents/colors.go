package dialogcomponents

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
)

// FormattedOptionalSelection provides the given optional dialog choice in a printable format.
func FormattedOptionalSelection(value fmt.Stringer, has bool, exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	if has {
		return colors.Green().Styled(value.String())
	}
	return colors.Green().Styled("(use global setting)")
}

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
	if selection == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled(selection)
}
