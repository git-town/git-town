package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterDevRemote() *cobra.Command {
	return &cobra.Command{
		Use: "dev-remote",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.DevRemote(gitdomain.Remotes{gitdomain.RemoteOrigin, "fork"}, dialog.Args[gitdomain.Remote]{
				Global: None[gitdomain.Remote](),
				Inputs: inputs,
				Local:  Some(gitdomain.RemoteOrigin),
			})
			return err
		},
	}
}
