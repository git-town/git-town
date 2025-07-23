package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterGitHubToken() *cobra.Command {
	return &cobra.Command{
		Use: "github-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.GitHubToken(dialog.Args[forgedomain.GitHubToken]{
				Global: None[forgedomain.GitHubToken](),
				Inputs: inputs,
				Local:  None[forgedomain.GitHubToken](),
			})
			return err
		},
	}
}
