package flags

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	commitMessageFileLong  = "message-file"
	commitMessageFileShort = "f"
)

// provides type-safe access to the CLI arguments of type gitdomain.MessageFile
func CommitMessageFile() (AddFunc, ReadCommitMessageFileFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(commitMessageFileLong, commitMessageFileShort, "", "Read the commit message from the given file (use \"-\" to read from STDIN)")
	}
	readFlag := func(cmd *cobra.Command) (Option[gitdomain.CommitMessageFile], error) {
		return readStringOptFlag[gitdomain.CommitMessageFile](cmd.Flags(), commitMessageFileLong)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFileFlagFunc reads gitdomain.CommitMessageFile from the CLI args.
type ReadCommitMessageFileFlagFunc func(*cobra.Command) (Option[gitdomain.CommitMessageFile], error)
