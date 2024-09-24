package debug

import (
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterBitbucketUsername() *cobra.Command {
	return &cobra.Command{
		Use: "bitbucket-username",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.BitbucketUsername(None[configdomain.BitbucketUsername](), dialogInputs.Next())
			return err
		},
	}
}
