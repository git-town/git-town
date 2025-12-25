package dialog

import (
	"cmp"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

func PerennialRegex(args Args[configdomain.PerennialRegex]) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	input, exit, errInput := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "perennial-regex",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          PerennialRegexHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.PerennialRegexPrompt,
		Title:         perennialRegexTitle,
	})
	newValue, errNewValue := configdomain.ParsePerennialRegex(input, "dialog")
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.PerennialRegex]()
	}
	fmt.Printf(messages.PerennialRegexResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, cmp.Or(errInput, errNewValue)
}
