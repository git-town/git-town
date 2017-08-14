package drivers

import (
	"errors"
	"fmt"
	"net/url"
)

// GitlabCodeHostingDriver provides tools for working with repositories
// on Gitlab.
type GitlabCodeHostingDriver struct {
	repository string
}

// NewGitlabCodeHostingDriver returns a new GitlabCodeHostingDriver instance
func NewGitlabCodeHostingDriver(repository string) *GitlabCodeHostingDriver {
	return &GitlabCodeHostingDriver{repository: repository}
}

// CanMergePullRequest returns whether or not MergePullRequest should be called when shipping
func (driver *GitlabCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, error) {
	return false, nil
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Gitlab
func (driver *GitlabCodeHostingDriver) GetNewPullRequestURL(branch, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", driver.repository, query.Encode())
}

// GetRepositoryURL returns the URL of the given repository on Gitlab
func (driver *GitlabCodeHostingDriver) GetRepositoryURL() string {
	return "https://gitlab.com/" + driver.repository
}

// MergePullRequest is unimplemented
func (driver *GitlabCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the Gitlab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}
