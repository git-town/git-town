package dialog

import (
	"cmp"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
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

func PerennialRegex(args TextArgs[configdomain.PerennialRegex]) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	input, exit, err1 := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "perennial-regex",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          PerennialRegexHelp,
		Prompt:        messages.PerennialRegexPrompt,
		TestInputs:    args.Inputs,
		Title:         perennialRegexTitle,
	})
	newValue, err2 := configdomain.ParsePerennialRegex(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.PerennialRegex]()
	}
	fmt.Printf(messages.PerennialRegexResult, dialogcomponents.FormattedSelection(newValue.String(), exit))
	return newValue, exit, cmp.Or(err1, err2)
}
