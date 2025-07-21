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
	contributionRegexTitle = `Contribution branch regex`
	contributionRegexHelp  = `
Branches matching this regular expression
will be treated as contribution branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "contribution".

If you're not sure what to enter here,
it's safe to leave it blank.
`
)

func ContributionRegex(args TextArgs[configdomain.ContributionRegex]) (Option[configdomain.ContributionRegex], dialogdomain.Exit, error) {
	input, exit, err1 := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "contribution-regex",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          contributionRegexHelp,
		Prompt:        messages.ContributionRegexPrompt,
		TestInputs:    args.Inputs,
		Title:         contributionRegexTitle,
	})
	newValue, err2 := configdomain.ParseContributionRegex(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.ContributionRegex]()
	}
	fmt.Printf(messages.ContributionRegexResult, dialogcomponents.FormattedSelection(newValue.String(), exit))
	return newValue, exit, cmp.Or(err1, err2)
}
