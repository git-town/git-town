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

`
)

func FeatureRegex(existingValue Option[configdomain.FeatureRegex], inputs dialogcomponents.TestInputs) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	value, exit, err1 := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: existingValue.String(),
		Help:          FeatureRegexHelp,
		Prompt:        "Feature regex: ",
		TestInputs:    inputs,
		Title:         featureRegexTitle,
	})
	fmt.Printf(messages.FeatureRegex, dialogcomponents.FormattedSelection(value, exit))
	featureRegex, err2 := configdomain.ParseFeatureRegex(value)
	return featureRegex, exit, cmp.Or(err1, err2)
}
