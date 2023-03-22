// Package flags provides helper methods for working with Cobra flags.
package flags

import (
	"github.com/spf13/cobra"
)

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)
