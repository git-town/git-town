package dialogcomponents

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
	if selection == "" {
		return colors.Green().Styled("(not provided)")
	}
	return colors.Green().Styled(selection)
}

// FormattedSelection provides the given dialog choice in a printable format.
func FormattedOptionalSelection(selection Option[fmt.Stringer], exit dialogdomain.Exit) string {
	if exit {
		return colors.Red().Styled("(aborted)")
	}
	if selected, has := selection.Get(); has {
		return colors.Green().Styled(selected.String())
	}
	return colors.Green().Styled("(use global setting)")
}
