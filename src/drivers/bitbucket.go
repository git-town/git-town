package drivers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/drivers/helpers"
)

// BitbucketCodeHostingDriver provides access to the API of Bitbucket installations.
type BitbucketCodeHostingDriver struct {
	git        gitRunner
	hostname   string
	originURL  string
	repository string
}

// LoadBitbucket provides a Bitbucket driver instance if the given repo configuration is for a Bitbucket repo,
// otherwise nil.
func LoadBitbucket(config config, git gitRunner) *BitbucketCodeHostingDriver {
	driverType := config.CodeHostingDriverName()
	originURL := config.GetRemoteOriginURL()
	hostname := helpers.GetURLHostname(originURL)
	manualHostName := config.CodeHostingOriginHostname()
	if manualHostName != "" {
		hostname = manualHostName
	}
	if driverType != "bitbucket" && hostname != "bitbucket.org" {
		return nil
	}
	return &BitbucketCodeHostingDriver{
		git:        git,
		hostname:   hostname,
		originURL:  originURL,
		repository: helpers.GetURLRepositoryName(originURL),
	}
}

func (d *BitbucketCodeHostingDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	return result, nil
}

func (d *BitbucketCodeHostingDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	branchSha, err := d.git.ShaForBranch(branch)
	if err != nil {
		return "", fmt.Errorf("cannot determine pull request URL from %q to %q: %w", branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{d.repository, branchSha[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{d.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *BitbucketCodeHostingDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *BitbucketCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *BitbucketCodeHostingDriver) HostingServiceName() string {
	return "Bitbucket"
}
