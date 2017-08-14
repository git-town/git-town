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
type BitbucketCodeHostingDriver struct {
	repository string
}

// NewBitbucketCodeHostingDriver returns a new BitbucketCodeHostingDriver instance
func NewBitbucketCodeHostingDriver(repository string) *BitbucketCodeHostingDriver {
	return &BitbucketCodeHostingDriver{repository: repository}
}

// CanMergePullRequest returns whether or not MergePullRequest should be called when shipping
func (driver *BitbucketCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, error) {
	return false, nil
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Bitbucket
func (driver *BitbucketCodeHostingDriver) GetNewPullRequestURL(branch, parentBranch string) string {
	query := url.Values{}
	query.Add("source", strings.Join([]string{driver.repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{driver.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("https://bitbucket.org/%s/pull-request/new?%s", driver.repository, query.Encode())
}

// GetRepositoryURL returns the URL where the given repository can be found
// on Bitbucket.com
func (driver *BitbucketCodeHostingDriver) GetRepositoryURL() string {
	return "https://bitbucket.org/" + driver.repository
}

// MergePullRequest is unimplemented
func (driver *BitbucketCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the BitBucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}
