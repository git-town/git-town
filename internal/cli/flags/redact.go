package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const redactLong = "redact"

// Redact provides type-safe access to the CLI arguments for the redact flag.
func Redact() (AddFunc, ReadRedactFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(redactLong, false, "hide API tokens from the output")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Redact, error) {
		return readBoolFlag[configdomain.Redact](cmd.Flags(), redactLong)
	}
	return addFlag, readFlag
}

// ReadRedactFlagFunc is the type signature for the function that reads the "redact" flag from the args to the given Cobra command.
type ReadRedactFlagFunc func(*cobra.Command) (configdomain.Redact, error)
