package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/messages"
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
func PerennialRegex(oldValue configdomain.PerennialRegex, inputs components.TestInput) (configdomain.PerennialRegex, bool, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          PerennialRegexHelp,
		Prompt:        messages.PerennialRegex,
		TestInput:     inputs,
		Title:         perennialRegexTitle,
	})
	fmt.Printf(messages.PerennialRegex, components.FormattedSelection(value, aborted))
	return configdomain.PerennialRegex(value), aborted, err
}
