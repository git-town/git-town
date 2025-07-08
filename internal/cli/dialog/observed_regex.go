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
	observedRegexTitle = `Observed branch regex`
	observedRegexHelp  = `
Branches matching this regular expression
will be treated as observed branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "observed".

`
)

func ObservedRegex(existingValue Option[configdomain.ObservedRegex], inputs dialogcomponents.TestInput) (Option[configdomain.ObservedRegex], dialogdomain.Exit, error) {
	value, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: existingValue.String(),
		Help:          observedRegexHelp,
		Prompt:        "Observed regex: ",
		TestInput:     inputs,
		Title:         observedRegexTitle,
	})
	fmt.Printf(messages.ObservedRegex, dialogcomponents.FormattedSelection(value, exit))
	if err != nil {
		return None[configdomain.ObservedRegex](), exit, err
	}
	observedRegex, err := configdomain.ParseObservedRegex(value)
	return observedRegex, false, err
}
