package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterBitbucketUsername() *cobra.Command {
	return &cobra.Command{
		Use: "bitbucket-username",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.BitbucketUsername(dialog.Args[forgedomain.BitbucketUsername]{
				Defaults: None[forgedomain.BitbucketUsername](),
				Global:   None[forgedomain.BitbucketUsername](),
				Inputs:   inputs,
				Local:    None[forgedomain.BitbucketUsername](),
			})
			return err
		},
	}
}
