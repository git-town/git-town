package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
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

func FeatureRegex(existingValue Option[configdomain.FeatureRegex], inputs components.TestInput) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: existingValue.String(),
		Help:          FeatureRegexHelp,
		Prompt:        "Feature regex: ",
		TestInput:     inputs,
		Title:         featureRegexTitle,
	})
	fmt.Printf(messages.FeatureRegex, components.FormattedSelection(value, aborted))
	featureRegex, err := configdomain.ParseFeatureRegex(value)
	return featureRegex, aborted, err
}
