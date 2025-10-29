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
	featureRegexTitle = `Feature branch regex`
	FeatureRegexHelp  = `
Branches matching this regular expression
will be treated as feature branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "feature".

If you're not sure what to enter here,
it's safe to leave it blank.

`
)

func FeatureRegex(args Args[configdomain.FeatureRegex]) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	input, exit, errInput := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "feature-regex",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          FeatureRegexHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.FeatureRegexPrompt,
		Title:         featureRegexTitle,
	})
	newValue, errNewValue := configdomain.ParseFeatureRegex(input, "dialog")
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.FeatureRegex]()
	}
	fmt.Printf(messages.FeatureRegexResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, cmp.Or(errInput, errNewValue)
}
