package debug

import (
	"os"

	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterPerennialRegex() *cobra.Command {
	return &cobra.Command{
		Use: "perennial-regex",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.PerennialRegex(None[configdomain.PerennialRegex](), dialogInputs.Next())
			return err
		},
	}
}
