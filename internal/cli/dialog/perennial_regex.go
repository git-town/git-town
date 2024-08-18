package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

const (
	perennialRegexTitle = `Regular expression for perennial branches`
	PerennialRegexHelp  = `
All branches whose names match this regular expression
are also considered perennial branches.

If you are not sure, leave this empty.

`
)

// PerennialRegex lets the user enter the GitHub API token.
func PerennialRegex(oldValue Option[configdomain.PerennialRegex], inputs components.TestInput) (Option[configdomain.PerennialRegex], bool, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          PerennialRegexHelp,
		Prompt:        "Perennial regex: ",
		TestInput:     inputs,
		Title:         perennialRegexTitle,
	})
	fmt.Printf(messages.PerennialRegex, components.FormattedSelection(value, aborted))
	return configdomain.ParsePerennialRegex(value), aborted, err
}
