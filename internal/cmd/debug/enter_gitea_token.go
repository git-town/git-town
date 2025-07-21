package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterGiteaToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitea-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.GiteaToken(dialog.Args[forgedomain.GiteaToken]{
				Global: None[forgedomain.GiteaToken](),
				Inputs: dialogInputs,
				Local:  None[forgedomain.GiteaToken](),
			})
			return err
		},
	}
}
