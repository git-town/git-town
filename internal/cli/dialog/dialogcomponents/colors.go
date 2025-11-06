package dialogcomponents

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FormattedOption provides the given optional dialog choice in a printable format.
func FormattedOption[T fmt.Stringer](userInput Option[T], hasGlobal bool, exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	if input, hasInput := userInput.Get(); hasInput {
		return colors.Green().Styled(input.String())
	}
	if hasGlobal {
		return colors.Green().Styled("(use global setting)")
	}
	return colors.Green().Styled("(not provided)")
}

// FormattedSecret provides the given API token in a printable format.
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
