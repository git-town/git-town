package flags

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const showSecretsLong = "show-secrets"

// ShowSecrets provides type-safe access to the CLI arguments for displaying sensitive information.
func ShowSecrets() (AddFunc, ReadShowSecretsFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(showSecretsLong, false, "display sensitive information like tokens unredacted")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.ShowSecrets, error) {
		return readBoolFlag[configdomain.ShowSecrets](cmd.Flags(), showSecretsLong)
	}
	return addFlag, readFlag
}

// ReadShowSecretsFlagFunc is the type signature for the function that reads the "show-secrets" flag from the args to the given Cobra command.
type ReadShowSecretsFlagFunc func(*cobra.Command) (configdomain.ShowSecrets, error)
