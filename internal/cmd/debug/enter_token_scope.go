package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterTokenScope() *cobra.Command {
	return &cobra.Command{
		Use: "token-scope",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.TokenScope(configdomain.ConfigScopeGlobal, dialogTestInputs.Next())
			return err
		},
	}
}
