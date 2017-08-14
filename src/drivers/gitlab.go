package drivers

import (
	"errors"
	"fmt"
	"net/url"
)

// GitlabCodeHostingDriver provides tools for working with repositories
// on Gitlab.
type GitlabCodeHostingDriver struct{}

// CanMergePullRequest returns whether or not MergePullRequest should be called when shipping
func (driver *GitlabCodeHostingDriver) CanMergePullRequest(options MergePullRequestOptions) (bool, error) {
	return false, nil
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Gitlab
func (driver *GitlabCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
}

// GetRepositoryURL returns the URL of the given repository on Gitlab
func (driver *GitlabCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://gitlab.com/" + repository
}

// MergePullRequest is unimplemented
func (driver *GitlabCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the Gitlab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}
