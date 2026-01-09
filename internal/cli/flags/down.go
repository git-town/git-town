package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	downLong  = "down"
	downShort = "d"
)

// Down provides type-safe access to the CLI arguments of type configdomain.Down.
// The flag can be used in two ways:
// - --down (uses default value of 1)
// - --down=2 (uses the specified integer value)
func Down() (AddFunc, ReadDownFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().IntP(downLong, downShort, 0, "commit into the given ancestor branch")
		cmd.Flags().Lookup(downLong).NoOptDefVal = "1"
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Down], error) {
		return readIntOptFlag[configdomain.Down](cmd.Flags(), downLong)
	}
	return addFlag, readFlag
}

// ReadDownFlagFunc is the type signature for the function that reads the "down" flag from the args to the given Cobra command.
type ReadDownFlagFunc func(*cobra.Command) (Option[configdomain.Down], error)
