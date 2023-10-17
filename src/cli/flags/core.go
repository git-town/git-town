// Package flags provides helper methods for working with Cobra flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
package flags

import (
	"github.com/spf13/cobra"
)

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)
