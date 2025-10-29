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
	observedRegexTitle = `Observed branch regex`
	observedRegexHelp  = `
Branches matching this regular expression
will be treated as observed branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "observed".

If you're not sure what to enter here,
it's safe to leave it blank.

`
)

func ObservedRegex(args Args[configdomain.ObservedRegex]) (Option[configdomain.ObservedRegex], dialogdomain.Exit, error) {
	input, exit, errInput := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "observed-regex",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          observedRegexHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.ObservedRegexPrompt,
		Title:         observedRegexTitle,
	})
	newValue, errNewValue := configdomain.ParseObservedRegex(input, "dialog")
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.ObservedRegex]()
	}
	fmt.Printf(messages.ObservedRegexResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, cmp.Or(errInput, errNewValue)
}
