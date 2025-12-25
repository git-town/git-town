package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const noVerifyLong = "no-verify"

// type-safe access to the CLI arguments of type configdomain.CommitHook
func NoVerify() (AddFunc, ReadNoVerifyFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(noVerifyLong, "", false, "do not run pre-commit hooks")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.CommitHook, error) {
		result, err := readBoolFlag[configdomain.CommitHook](cmd.Flags(), noVerifyLong)
		return !result, err
	}
	return addFlag, readFlag
}

// ReadNoVerifyFlagFunc is the type signature for the function that reads the "no-verify" flag from the args to the given Cobra command.
type ReadNoVerifyFlagFunc func(*cobra.Command) (configdomain.CommitHook, error)
