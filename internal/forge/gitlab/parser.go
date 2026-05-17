package gitlab

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func parseMergeRequest(mergeRequest *gitlab.BasicMergeRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Active:       mergeRequest.State == "opened",
		MergeWithAPI: true,
		Number:       forgedomain.ProposalNumber(mergeRequest.IID),
		Source:       gitdomain.LocalBranchNameOrPanic(stringss.Trimmed(mergeRequest.SourceBranch)), // we can assume the GitLab API provides correct strings
		Target:       gitdomain.LocalBranchNameOrPanic(stringss.Trimmed(mergeRequest.TargetBranch)), // we can assume the GitLab API provides correct strings
		Title:        gitdomain.ProposalTitle(mergeRequest.Title),
		Body:         gitdomain.NewProposalBodyOpt(mergeRequest.Description),
		URL:          mergeRequest.WebURL,
	}
}
