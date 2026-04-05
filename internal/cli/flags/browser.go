package flags

import (
	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const browserLong = "browser"

// type-safe access to the CLI arguments of type browserdomain.Browser
func Browser() (AddFunc, ReadBrowserFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().String(browserLong, "", "the browser executable to use")
	}
	readFlag := func(cmd *cobra.Command) (Option[browserdomain.Browser], error) {
		return readStringOptFlag[browserdomain.Browser](cmd.Flags(), browserLong)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadBrowserFlagFunc func(*cobra.Command) (Option[browserdomain.Browser], error)
