package flags

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	bodyFileLong  = "body-file" // long form of the "body-file" CLI flag
	bodyFileShort = "f"         // short form of the "body-file" CLI flag
)

// provides type-safe access to the CLI arguments of type gitdomain.ProposalBodyFile
func ProposalBodyFile() (AddFunc, ReadProposalBodyFileFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(bodyFileLong, bodyFileShort, "", "Read the proposal body from the given file (use \"-\" to read from STDIN)")
	}
	readFlag := func(cmd *cobra.Command) (Option[gitdomain.ProposalBodyFile], error) {
		return readStringOptFlag[gitdomain.ProposalBodyFile](cmd.Flags(), bodyFileLong)
	}
	return addFlag, readFlag
}

// ReadProposalBodyFileFlagFunc reads gitdomain.ProposalBodyFile from the CLI args.
type ReadProposalBodyFileFlagFunc func(*cobra.Command) (Option[gitdomain.ProposalBodyFile], error)
