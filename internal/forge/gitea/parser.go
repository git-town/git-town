package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        gitdomain.ProposalTitle(pullRequest.Title),
		Body:         NewOption(gitdomain.ProposalBody(pullRequest.Body)),
		URL:          pullRequest.HTMLURL,
	}
}
