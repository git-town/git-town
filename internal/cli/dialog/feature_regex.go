package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	featureRegexTitle = `Feature branch regex`
	featureRegexHelp  = `
Branches matching this regular expression
will be treated as feature branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "feature".

`
)

func FeatureRegex(args CommonArgs) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[configdomain.FeatureRegex]{
		ConfigFileValue: args.ConfigFile.FeatureRegex,
		DialogName:      "feature-regex",
		HelpText:        featureRegexHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.FeatureRegex,
		ParseFunc:       configdomain.ParseFeatureRegex,
		PrintResultFunc: dialogcomponents.FormattedSelection,
		Prompt:          "Feature Regex: ",
		ResultMessage:   messages.FeatureRegex,
		Title:           featureRegexTitle,
		UnscopedValue:   args.UnscopedGitConfig.FeatureRegex,
	})
}
