package drivers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/src/drivers/helpers"
)

// bitbucketCodeHostingDriver provides access to the API of Bitbucket installations.
type bitbucketCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
	shell      shell
}

// LoadBitbucket provides a Bitbucket driver instance if the given repo configuration is for a Bitbucket repo,
// otherwise nil.
func LoadBitbucket(config config, shell shell) CodeHostingDriver {
	driverType := config.GetCodeHostingDriverName()
	originURL := config.GetRemoteOriginURL()
	hostname := helpers.GetURLHostname(originURL)
	configuredHostName := config.GetCodeHostingOriginHostname()
	if configuredHostName != "" {
		hostname = configuredHostName
	}
	if driverType != "bitbucket" && hostname != "bitbucket.org" {
		return nil
	}
	return &bitbucketCodeHostingDriver{
		originURL:  originURL,
		hostname:   hostname,
		repository: helpers.GetURLRepositoryName(originURL),
		shell:      shell,
	}
}

func (d *bitbucketCodeHostingDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	return result, nil
}

func (d *bitbucketCodeHostingDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	branchSha, err := d.shell.ShaForBranch(branch)
	if err != nil {
		return "", fmt.Errorf("cannot determine pull request URL from %q to %q: %w", branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{d.repository, branchSha[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{d.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *bitbucketCodeHostingDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *bitbucketCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *bitbucketCodeHostingDriver) HostingServiceName() string {
	return "Bitbucket"
}
