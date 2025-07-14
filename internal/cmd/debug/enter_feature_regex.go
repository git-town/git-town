package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterFeatureRegex() *cobra.Command {
	return &cobra.Command{
		Use: "feature-regex",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.FeatureRegex]{
				ConfigFileValue: None[configdomain.FeatureRegex](),
				HelpText:        dialog.FeatureRegexHelp,
				Inputs:          dialogInputs,
				LocalValue:      None[configdomain.FeatureRegex](),
				ParseFunc:       configdomain.ParseFeatureRegex,
				Prompt:          "Your feature regex: ",
				ResultMessage:   messages.FeatureRegex,
				Title:           dialog.FeatureRegexTitle,
				UnscopedValue:   None[configdomain.FeatureRegex](),
			})
			return err
		},
	}
}
