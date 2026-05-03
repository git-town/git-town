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
		defineNegatedFlag(cmd.Flags(), browserLong, "don't open any browser windows")
	}
	readFlag := func(cmd *cobra.Command) (Option[browserdomain.BrowserExecutable], error) {
		negatedOpt, err := readBoolOptFlag[bool](cmd.Flags(), negate+browserLong)
		if err != nil {
			return None[browserdomain.BrowserExecutable](), err
		}
		negated, hasNegated := negatedOpt.Get()
		if hasNegated && negated {
			return Some(browserdomain.NoBrowser), nil
		}
		value, err := readStringOptFlag[string](cmd.Flags(), browserLong)
		if err != nil {
			return None[browserdomain.BrowserExecutable](), err
		}
		return browserdomain.ParseBrowserOpt(value)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadBrowserFlagFunc func(*cobra.Command) (Option[browserdomain.BrowserExecutable], error)
