package dialog

import (
	"cmp"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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

func PerennialRegex(oldValue Option[configdomain.PerennialRegex], inputs components.TestInput) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	value, aborted, err1 := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          PerennialRegexHelp,
		Prompt:        "Perennial regex: ",
		TestInput:     inputs,
		Title:         perennialRegexTitle,
	})
	fmt.Printf(messages.PerennialRegex, components.FormattedSelection(value, aborted))
	perennialRegex, err2 := configdomain.ParsePerennialRegex(value)
	return perennialRegex, aborted, cmp.Or(err1, err2)
}
