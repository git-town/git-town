package flags

import (
	"github.com/spf13/cobra"
)

// AddFunc defines the type signature for helper functions that add a CLI flag to a Cobra command.
type AddFunc func(*cobra.Command)
