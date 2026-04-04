package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	upLong  = "up"
	upShort = "u"
)

// Up provides type-safe access to the CLI arguments of type configdomain.Up.
// The flag can be used in two ways:
// - --up (uses default value of 1)
// - --up=2 (uses the specified integer value)
func Up() (AddFunc, ReadUpFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().UintP(upLong, upShort, 0, "commit into the given ancestor branch")
		cmd.Flags().Lookup(upLong).NoOptDefVal = "1"
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Up], error) {
		return readUintOptFlag[configdomain.Up](cmd.Flags(), upLong)
	}
	return addFlag, readFlag
}

// ReadUpFlagFunc is the type signature for the function that reads the "up" flag from the args to the given Cobra command.
type ReadUpFlagFunc func(*cobra.Command) (Option[configdomain.Up], error)
