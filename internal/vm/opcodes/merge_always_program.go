package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// MergeAlwaysProgram merges the feature branch into its parent by always creating a merge comment (merge --no-ff).
type MergeAlwaysProgram struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
}

func (self *MergeAlwaysProgram) Abort() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeAlwaysProgram) AutomaticUndoError() error {
	return errors.New(messages.ShipExitMergeError)
}

func (self *MergeAlwaysProgram) Run(args shared.RunArgs) error {
	// Reverting parent is intentionally not supported due to potential confusion
	// caused by reverted merge commit. See
	// <https://github.com/git/git/blob/master/Documentation/howto/revert-a-faulty-merge.adoc>
	// for more information.
	useMessage := configdomain.UseCustomMessageOr(self.CommitMessage, configdomain.EditDefaultMessage())
	return args.Git.MergeNoFastForward(args.Frontend, useMessage, self.Branch)
}

func (self *MergeAlwaysProgram) ShouldUndoOnError() bool {
	return true
}
