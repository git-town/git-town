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

func FeatureRegex(args DialogArgs[configdomain.FeatureRegex]) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	input, exit, err1 := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "feature-regex",
		ExistingValue: args.Local.Or(args.Global).String(),
		Help:          FeatureRegexHelp,
		Prompt:        messages.FeatureRegexPrompt,
		TestInputs:    args.Inputs,
		Title:         featureRegexTitle,
	})
	newValue, err2 := configdomain.ParseFeatureRegex(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.FeatureRegex]()
	}
	fmt.Printf(messages.FeatureRegex, dialogcomponents.FormattedSelection(newValue.String(), exit))
	return newValue, exit, cmp.Or(err1, err2)
}
