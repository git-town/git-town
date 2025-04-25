package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

const (
	perennialRegexTitle = `Perennial branch Regex`
	PerennialRegexHelp  = `
Any branch name matching this regular expression
will be treated as a perennial branch.

Example: ^release-.+

If you're not sure what to enter here,
it's safe to leave it blank.

`
)

func PerennialRegex(oldValue Option[configdomain.PerennialRegex], inputs components.TestInput) (Option[configdomain.PerennialRegex], bool, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          PerennialRegexHelp,
		Prompt:        "Perennial regex: ",
		TestInput:     inputs,
		Title:         perennialRegexTitle,
	})
	if err != nil {
		return None[configdomain.PerennialRegex](), false, err
	}
	fmt.Printf(messages.PerennialRegex, components.FormattedSelection(value, aborted))
	perennialRegex, err := configdomain.ParsePerennialRegex(value)
	return perennialRegex, aborted, err
}
