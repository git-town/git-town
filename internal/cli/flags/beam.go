package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const beamLong = "beam"

// type-safe access to the CLI arguments of type configdomain.Beam
func Beam() (AddFunc, ReadBeamFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(beamLong, "b", false, "beam some commits from this branch to the new branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Beam, error) {
		value, err := cmd.Flags().GetBool(beamLong)
		return configdomain.Beam(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadBeamFlagFunc func(*cobra.Command) (configdomain.Beam, error)
