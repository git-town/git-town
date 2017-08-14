package drivers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

// BitbucketCodeHostingDriver provides functionality for working with
// repositories hosted on Bitbucket
type BitbucketCodeHostingDriver struct{}

// CanMergePullRequest returns whether or not MergePullRequest should be called when shipping
func (driver BitbucketCodeHostingDriver) CanMergePullRequest(options MergePullRequestOptions) (bool, error) {
	return false, nil
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Bitbucket
func (driver BitbucketCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("source", strings.Join([]string{repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{repository, "", parentBranch}, ":"))
	return fmt.Sprintf("https://bitbucket.org/%s/pull-request/new?%s", repository, query.Encode())
}

// GetRepositoryURL returns the URL where the given repository can be found
// on Bitbucket.com
func (driver BitbucketCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://bitbucket.org/" + repository
}

// MergePullRequest is unimplemented
func (driver BitbucketCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the BitBucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}
