package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterFeatureRegex() *cobra.Command {
	return &cobra.Command{
		Use: "feature-regex",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.FeatureRegex(dialog.TextArgs[configdomain.FeatureRegex]{
				Global: None[configdomain.FeatureRegex](),
				Inputs: dialogInputs,
				Local:  None[configdomain.FeatureRegex](),
			})
			return err
		},
	}
}
