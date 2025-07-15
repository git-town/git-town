package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	perennialRegexTitle = `Perennial branch Regex`
	perennialRegexHelp  = `
Any branch name matching this regular expression
will be treated as a perennial branch.

Example: ^release-.+

If you're not sure what to enter here,
it's safe to leave it blank.

`
)

func PerennialRegex(args CommonArgs) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	return ConfigStringDialog(ConfigStringDialogArgs[configdomain.PerennialRegex]{
		ConfigFileValue: args.ConfigFile.PerennialRegex,
		HelpText:        perennialRegexHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.PerennialRegex,
		ParseFunc:       configdomain.ParsePerennialRegex,
		PrintResultFunc: dialogcomponents.FormattedSelection,
		Prompt:          "Perennial Regex: ",
		ResultMessage:   messages.PerennialRegex,
		Title:           perennialRegexTitle,
		UnscopedValue:   args.UnscopedGitConfig.PerennialRegex,
	})
}

type CommonArgs struct {
	ConfigFile        configdomain.PartialConfig
	Inputs            dialogcomponents.TestInputs
	LocalGitConfig    configdomain.PartialConfig
	UnscopedGitConfig configdomain.PartialConfig
}
