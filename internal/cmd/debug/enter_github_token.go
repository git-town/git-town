package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterGitHubToken() *cobra.Command {
	return &cobra.Command{
		Use: "github-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.GitHubToken(dialog.CommonArgs{
				ConfigFile:        configdomain.EmptyPartialConfig(),
				Inputs:            dialogInputs,
				LocalGitConfig:    configdomain.EmptyPartialConfig(),
				UnscopedGitConfig: configdomain.EmptyPartialConfig(),
			})
			return err
		},
	}
}
