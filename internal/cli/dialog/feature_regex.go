package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

const (
	featureRegexTitle = `Feature branch regex`
	FeatureRegexHelp  = `
Branches matching this regular expression will be treated as feature branches.
This setting only applies if the "default-branch-type"
is set to something other than "feature".

`
)

func FeatureRegex(existingValue Option[configdomain.FeatureRegex], inputs components.TestInput) (Option[configdomain.FeatureRegex], bool, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: existingValue.String(),
		Help:          FeatureRegexHelp,
		Prompt:        "Feature regex: ",
		TestInput:     inputs,
		Title:         featureRegexTitle,
	})
	if err != nil {
		return None[configdomain.FeatureRegex](), false, err
	}
	fmt.Printf(messages.FeatureRegex, components.FormattedSelection(value, aborted))
	featureRegex, err := configdomain.ParseFeatureRegex(value)
	return featureRegex, aborted, err
}
