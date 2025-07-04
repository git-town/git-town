package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterBitbucketAppPassword() *cobra.Command {
	return &cobra.Command{
		Use: "bitbucket-app-password",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.BitbucketAppPassword(None[forgedomain.BitbucketAppPassword](), dialogInputs.Next())
			return err
		},
	}
}
