package debug

import (
	"os"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterGitLabToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitlab-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.GitLabToken(None[configdomain.GitLabToken](), dialogInputs.Next())
			return err
		},
	}
}
