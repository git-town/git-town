package dialog

import (
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

`
)

func ContributionRegex(existingValue Option[configdomain.ContributionRegex], inputs dialogcomponents.TestInput) (Option[configdomain.ContributionRegex], dialogdomain.Exit, error) {
	value, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: existingValue.String(),
		Help:          contributionRegexHelp,
		Prompt:        "Contribution regex: ",
		TestInput:     inputs,
		Title:         contributionRegexTitle,
	})
	fmt.Printf(messages.ContributionRegex, dialogcomponents.FormattedSelection(value, exit))
	if err != nil {
		return None[configdomain.ContributionRegex](), exit, err
	}
	contributionRegex, err := configdomain.ParseContributionRegex(value)
	return contributionRegex, false, err
}
