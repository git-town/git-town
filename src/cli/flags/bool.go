package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// type-safe access to boolean Cobra command-line flags
func Bool(name, short, desc string, persistent FlagType) (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		switch persistent {
		case FlagTypePersistent:
			cmd.PersistentFlags().BoolP(name, short, false, desc)
		case FlagTypeNonPersistent:
			cmd.Flags().BoolP(name, short, false, desc)
		}
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(name)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a boolean %q flag", cmd.Name(), name))
		}
		return value
	}
	return addFlag, readFlag
}

// ReadBoolFlagFunc defines the type signature for helper functions that provide the value a boolean CLI flag associated with a Cobra command.
type ReadBoolFlagFunc func(*cobra.Command) bool

type FlagType int

const (
	FlagTypePersistent FlagType = iota
	FlagTypeNonPersistent
)
