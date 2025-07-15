package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterPerennialRegex() *cobra.Command {
	return &cobra.Command{
		Use: "perennial-regex",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.PerennialRegex(dialog.CommonArgs{
				ConfigFile:        configdomain.PartialConfig{},
				Inputs:            dialogInputs,
				LocalGitConfig:    configdomain.PartialConfig{},
				UnscopedGitConfig: configdomain.PartialConfig{},
			})
			return err
		},
	}
}
