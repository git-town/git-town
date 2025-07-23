package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterForgeType() *cobra.Command {
	return &cobra.Command{
		Use: "forge-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.ForgeType(dialog.Args[forgedomain.ForgeType]{
				Global: None[forgedomain.ForgeType](),
				Inputs: inputs,
				Local:  None[forgedomain.ForgeType](),
			})
			return err
		},
	}
}
