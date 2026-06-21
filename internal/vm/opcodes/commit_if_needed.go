package opcodes

import (
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/vm/shared"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// CommitIfNeeded commits all open changes as a new commit,
// but only if there are changes to commit.
type CommitIfNeeded struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
}

func (self *CommitIfNeeded) Run(args shared.RunArgs) error {
	repoStatus, err := args.Git.RepoStatus(args.Backend)
	if err != nil {
		return err
	}
	if !repoStatus.NeedsToCommit() {
		return nil
	}
	args.PrependOpcodes(&Commit{
		AuthorOverride:                 self.AuthorOverride,
		FallbackToDefaultCommitMessage: self.FallbackToDefaultCommitMessage,
		Message:                        self.Message,
	})
	return nil
}
