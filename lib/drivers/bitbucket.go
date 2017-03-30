package drivers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/lib/git"
)

type BitbucketCodeHostingDriver struct{}

func (driver BitbucketCodeHostingDriver) GetNewPullRequestUrl(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("source", strings.Join([]string{repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{repository, "", parentBranch}, ":"))
	return fmt.Sprintf("https://bitbucket.org/%s/pull-request/new?%s", repository, query.Encode())
}

func (driver BitbucketCodeHostingDriver) GetRepositoryUrl(repository string) string {
	return "https://bitbucket.org/" + repository
}
