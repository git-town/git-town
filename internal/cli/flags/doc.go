// Package flags provides helper methods for working with Cobra flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
package flags

import (
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)

func readBool[T ~bool](flags *pflag.FlagSet, name string) (Option[T], error) {
	if !flags.Changed(name) {
		return None[T](), nil
	}
	value, err := flags.GetBool(detachedLong)
	if err != nil {
		return None[T](), err
	}
	return Some(T(value)), nil
}
