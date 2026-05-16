package gitlab

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func parseMergeRequest(mergeRequest *gitlab.BasicMergeRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Active:       mergeRequest.State == "opened",
		MergeWithAPI: true,
		Number:       forgedomain.ProposalNumber(mergeRequest.IID),
		Source:       gitdomain.LocalBranchNameOrPanic(mergeRequest.SourceBranch),
		Target:       gitdomain.LocalBranchNameOrPanic(mergeRequest.TargetBranch),
		Title:        gitdomain.ProposalTitle(mergeRequest.Title),
		Body:         gitdomain.NewProposalBodyOpt(mergeRequest.Description),
		URL:          mergeRequest.WebURL,
	}
}
