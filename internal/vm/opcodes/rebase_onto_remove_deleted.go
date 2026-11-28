package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// RebaseOntoRemoveDeleted rebases the current branch against the target branch, while moving the target branch onto the Onto branch.
// If there are merge conflicts,
type RebaseOntoRemoveDeleted struct {
	BranchToRebaseOnto gitdomain.LocalBranchName
	CommitsToRemove    gitdomain.BranchName
}

func (self *RebaseOntoRemoveDeleted) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoRemoveDeleted) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoRemoveDeleted) Run(args shared.RunArgs) error {
	rebaseAlreadyInProgress, rebaseProbeErr := args.Git.HasRebaseInProgress(args.Backend)
	rebaseErr := args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseOnto.Location(), self.CommitsToRemove.Location())
	if rebaseErr == nil || args.Config.Value.NormalConfig.AutoResolve.NoAutoResolve() {
		return rebaseErr
	}
	// Here the rebase-onto has failed.
	// Recover automatically only for a conflicted rebase that this opcode newly started.
	// A rebase that was already active belongs to an earlier operation and must never be resolved here.
	// Other failures, for example a stale --onto ref, must surface unchanged.
	if rebaseProbeErr != nil || rebaseAlreadyInProgress {
		return rebaseErr
	}
	rebaseInProgress, err := args.Git.HasRebaseInProgress(args.Backend)
	if err != nil || !rebaseInProgress {
		return rebaseErr
	}
	conflictingFiles, err := args.Git.FileConflicts(args.Backend)
	if err != nil || len(conflictingFiles) == 0 {
		return rebaseErr
	}
	for _, conflictingFile := range conflictingFiles {
		if conflictingFile.BaseChange.IsNone() || !self.conflictCausedByRemovedCommits(args, conflictingFile) {
			return rebaseErr
		}
	}
	opcodes := []shared.Opcode{}
	for _, conflictingFile := range conflictingFiles {
		if conflictingChange, has := conflictingFile.CurrentBranchChange.Get(); has {
			opcodes = append(
				opcodes,
				&ConflictResolve{
					FilePath:   conflictingChange.FilePath,
					Resolution: gitdomain.ConflictResolutionTheirs,
				},
				&FileStage{
					FilePath: conflictingChange.FilePath,
				},
			)
		} else if baseChange, has := conflictingFile.BaseChange.Get(); has {
			opcodes = append(opcodes, &FileRemove{
				FilePath: baseChange.FilePath,
			})
		}
	}
	opcodes = append(opcodes, &RebaseContinue{})
	args.PrependOpcodes(opcodes...)
	return nil
}

// conflictCausedByRemovedCommits indicates whether the conflict baseline and target state prove
// that removing the given commits caused this conflict. Independent edit conflicts do not satisfy this test:
// in those conflicts the removed commits did not change the file at all, or the target contains its own version.
func (self *RebaseOntoRemoveDeleted) conflictCausedByRemovedCommits(args shared.RunArgs, conflictingFile git.FileConflict) bool {
	filePath, hasFilePath := conflictFilePath(conflictingFile).Get()
	if !hasFilePath {
		return false
	}
	mergeBaseText, err := args.Backend.QueryTrim("git", "merge-base", self.BranchToRebaseOnto.String(), self.CommitsToRemove.String())
	if err != nil {
		return false
	}
	mergeBase, err := gitdomain.NewSHA(mergeBaseText)
	if err != nil {
		return false
	}
	mergeBaseBlob, err := args.Git.ContentBlobInfoLiteral(args.Backend, mergeBase.Location(), filePath)
	if err != nil {
		return false
	}
	removedBlob, err := args.Git.ContentBlobInfoLiteral(args.Backend, self.CommitsToRemove.Location(), filePath)
	if err != nil {
		return false
	}
	targetBlob, err := args.Git.ContentBlobInfoLiteral(args.Backend, self.BranchToRebaseOnto.Location(), filePath)
	if err != nil {
		return false
	}
	if !sameBlob(conflictingFile.BaseChange, removedBlob) {
		return false
	}
	if sameBlob(mergeBaseBlob, removedBlob) {
		return false
	}
	return sameBlob(targetBlob, mergeBaseBlob) || sameBlob(targetBlob, removedBlob)
}

func conflictFilePath(conflictingFile git.FileConflict) Option[string] {
	if incomingChange, hasIncomingChange := conflictingFile.IncomingChange.Get(); hasIncomingChange {
		return Some(incomingChange.FilePath)
	}
	if currentBranchChange, hasCurrentBranchChange := conflictingFile.CurrentBranchChange.Get(); hasCurrentBranchChange {
		return Some(currentBranchChange.FilePath)
	}
	if baseChange, hasBaseChange := conflictingFile.BaseChange.Get(); hasBaseChange {
		return Some(baseChange.FilePath)
	}
	return None[string]()
}

func sameBlob(first, second Option[git.Blob]) bool {
	firstBlob, hasFirstBlob := first.Get()
	secondBlob, hasSecondBlob := second.Get()
	if hasFirstBlob != hasSecondBlob {
		return false
	}
	if !hasFirstBlob {
		return true
	}
	return firstBlob.SHA == secondBlob.SHA && firstBlob.Permission == secondBlob.Permission
}
