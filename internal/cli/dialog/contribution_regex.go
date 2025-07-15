package dialog

import (
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

func ContributionRegex(args CommonArgs) (Option[configdomain.ContributionRegex], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[configdomain.ContributionRegex]{
		ConfigFileValue: args.ConfigFile.ContributionRegex,
		HelpText:        contributionRegexHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.ContributionRegex,
		ParseFunc:       configdomain.ParseContributionRegex,
		Prompt:          "Contribution Regex: ",
		ResultMessage:   messages.ContributionRegex,
		Title:           contributionRegexTitle,
		UnscopedValue:   args.UnscopedGitConfig.ContributionRegex,
	})
}
