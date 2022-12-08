package hosting

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/giturl"
)

// BitbucketDriver provides access to the API of Bitbucket installations.
type BitbucketDriver struct {
	git          gitRunner
	hostname     string
	originURL    string
	organization string
	repository   string
}

// NewBitbucketDriver provides a Bitbucket driver instance if the given repo configuration is for a Bitbucket repo,
// otherwise nil.
func NewBitbucketDriver(config config, git gitRunner) *BitbucketDriver {
	driverType := config.HostingService()
	originURL := config.OriginURL()
	url := giturl.Parse(originURL)
	if url == nil {
		return nil
	}
	manualOrigin := config.OriginOverride()
	if manualOrigin != "" {
		url.Host = manualOrigin
	}
	if driverType != "bitbucket" && url.Host != "bitbucket.org" {
		return nil
	}
	return &BitbucketDriver{
		git:          git,
		hostname:     url.Host,
		organization: url.Org,
		originURL:    originURL,
		repository:   url.Repo,
	}
}

func (d *BitbucketDriver) LoadPullRequestInfo(branch, parentBranch string) (PullRequestInfo, error) {
	return PullRequestInfo{}, nil
}

func (d *BitbucketDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	branchSha, err := d.git.ShaForBranch(branch)
	if err != nil {
		return "", fmt.Errorf("cannot determine pull request URL from %q to %q: %w", branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{d.organization + "/" + d.repository, branchSha[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{d.organization + "/" + d.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *BitbucketDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.organization, d.repository)
}

func (d *BitbucketDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *BitbucketDriver) HostingServiceName() string {
	return "Bitbucket"
}
