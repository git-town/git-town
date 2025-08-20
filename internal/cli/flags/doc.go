// Package flags provides helper methods for working with Cobra flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
package flags

import (
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const negate = "no-"

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)

func defineNegatableBoolFlag(flags *pflag.FlagSet, name string) {
	flags.Bool(name, true, "whether to auto-resolve phantom merge conflicts")
	negateName := negate + name
	flags.Bool(negateName, false, "")
	_ = flags.MarkHidden(negateName)
}

func readBoolFlag[T ~bool](flags *pflag.FlagSet, name string) (T, error) { //nolint:ireturn
	value, err := flags.GetBool(name)
	return T(value), err
}

func readBoolOptFlag[T ~bool](flags *pflag.FlagSet, name string) (Option[T], error) {
	value, err := flags.GetBool(name)
	return NewOption(T(value)), err
}

func readNegatableBoolFlag[T ~bool](flags *pflag.FlagSet, name string) (Option[T], error) {
	if value, err := readBoolOptFlag[T](flags, name); value.IsSome() || err != nil {
		return value, err
	}
	valueOpt, err := readBoolOptFlag[T](flags, negate+name)
	if value, has := valueOpt.Get(); has {
		return Some(T(!bool(value))), err
	}
	return None[T](), err
}

func readStringOptFlag[T ~string](flags *pflag.FlagSet, name string) (Option[T], error) {
	value, err := flags.GetString(name)
	return NewOption(T(value)), err
}
