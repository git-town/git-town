package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
)

func FilterPullRequests(pullRequests []*gitea.PullRequest, branch, target gitdomain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() && pullRequest.Base.Name == target.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func FilterPullRequests2(pullRequests []*gitea.PullRequest, branch gitdomain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func parsePullRequest(pullRequest *gitea.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Active:       pullRequest.State == "open",
		MergeWithAPI: pullRequest.Mergeable,
		Number:       forgedomain.ProposalNumber(pullRequest.Index),
		Source:       gitdomain.LocalBranchNameOrPanic(stringss.Trimmed(pullRequest.Head.Ref)), // we can assume the Gitea API provides correct strings
		Target:       gitdomain.LocalBranchNameOrPanic(stringss.Trimmed(pullRequest.Base.Ref)), // we can assume the Gitea API provides correct strings
		Title:        gitdomain.ProposalTitle(pullRequest.Title),
		Body:         gitdomain.NewProposalBodyOpt(pullRequest.Body),
		URL:          pullRequest.HTMLURL,
	}
}
