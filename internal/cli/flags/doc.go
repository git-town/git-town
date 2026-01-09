// Package flags provides helper methods for working with Cobra flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
package flags

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const negate = "no-"

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)

// defines the negated form of the flag with the given name and default value
// You must define the non-negated version yourself, because there are too many things to configure there
func defineNegatedFlag(flags *pflag.FlagSet, name, desc string) {
	negateName := negate + name
	flags.Bool(negateName, false, desc)
}

// provides the value of the CLI flag with the given name and bool-based type
func readBoolFlag[T ~bool](flags *pflag.FlagSet, name string) (T, error) { //nolint:ireturn
	value, err := flags.GetBool(name)
	return T(value), err
}

// provides the value of the CLI flag with the given name and optional bool-based type
func readBoolOptFlag[T ~bool](flags *pflag.FlagSet, name string) (Option[T], error) {
	if flags.Changed(name) {
		value, err := flags.GetBool(name)
		return Some(T(value)), err
	}
	return None[T](), nil
}

// provides the value of the CLI flag with the given name and optional int-based type
func readIntOptFlag[T ~int](flags *pflag.FlagSet, name string) (Option[T], error) {
	if flags.Changed(name) {
		value, err := flags.GetInt(name)
		return Some(T(value)), err
	}
	return None[T](), nil
}

// provides the value of the CLI flag with the given name and optional negatable bool-based type
func readNegatableFlag[T ~bool](flags *pflag.FlagSet, name string) (Option[T], error) {
	if value, err := readBoolOptFlag[T](flags, name); value.IsSome() || err != nil {
		return value, err
	}
	valueOpt, err := readBoolOptFlag[T](flags, negate+name)
	if value, has := valueOpt.Get(); has {
		return Some(T(!bool(value))), err
	}
	return None[T](), err
}

// provides the value of the CLI flag with the given name and optional string-based type
func readStringOptFlag[T ~string](flags *pflag.FlagSet, name string) (Option[T], error) {
	value, err := flags.GetString(name)
	return NewOption(T(value)), err
}
