package dialog

import (
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

func ObservedRegex(args CommonArgs) (Option[configdomain.ObservedRegex], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[configdomain.ObservedRegex]{
		ConfigFileValue: args.ConfigFile.ObservedRegex,
		HelpText:        observedRegexHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.ObservedRegex,
		ParseFunc:       configdomain.ParseObservedRegex,
		PrintResultFunc: dialogcomponents.FormattedSelection,
		Prompt:          "Observed Regex: ",
		ResultMessage:   messages.ObservedRegex,
		Title:           observedRegexTitle,
		UnscopedValue:   args.UnscopedGitConfig.ObservedRegex,
	})
}
