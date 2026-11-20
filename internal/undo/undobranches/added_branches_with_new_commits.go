package undobranches

import "github.com/git-town/git-town/v24/internal/git/gitdomain"

// AddedBranchesWithNewCommits provides the branches that the last Git Town command created
// and that point to a different commit now than when that command finished.
// Undoing that command deletes these branches, which would lose these new commits.
func AddedBranchesWithNewCommits(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, touchedBranches []gitdomain.BranchName, currentBranches gitdomain.BranchInfos) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot).KeepOnly(touchedBranches)
	for _, addedBranch := range branchSpans.Changes().LocalAdded {
		endBranch, hasEndBranch := endBranchesSnapshot.Branches.FindByLocalName(addedBranch).Get()
		if !hasEndBranch {
			continue
		}
		endSHA, hasEndSHA := endBranch.LocalSHA().Get()
		if !hasEndSHA {
			continue
		}
		currentBranch, hasCurrentBranch := currentBranches.FindByLocalName(addedBranch).Get()
		if !hasCurrentBranch {
			continue
		}
		currentSHA, hasCurrentSHA := currentBranch.LocalSHA().Get()
		if hasCurrentSHA && currentSHA != endSHA {
			result = append(result, addedBranch)
		}
	}
	return result
}
