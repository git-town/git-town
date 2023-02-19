package hosting

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v7/src/giturl"
	"golang.org/x/oauth2"
)

// GiteaConfig contains connection information for Gitea-based hosting platforms.
type GiteaConfig struct {
	apiToken   string
	hostname   string
	originURL  string
	owner      string
	repository string
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewGiteaConfig(url giturl.Parts, config config) *GiteaConfig {
	driverType := config.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if driverType != "gitea" && url.Host != "gitea.com" {
		return nil
	}
	return &GiteaConfig{
		originURL:  config.OriginURL(),
		hostname:   url.Host,
		apiToken:   config.GiteaToken(),
		owner:      url.Org,
		repository: url.Repo,
	}
}

// Driver provides a fully configured driver instance for the Gitea API.
func (c GiteaConfig) Driver(log logFn) GiteaDriver {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", c.hostname), httpClient)
	return GiteaDriver{
		client:      giteaClient,
		GiteaConfig: c,
		log:         log,
	}
}

func (c GiteaConfig) HostingServiceName() string {
	return "Gitea"
}

func (c GiteaConfig) NewPullRequestURL(branch string, parentBranch string) (string, error) {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c GiteaConfig) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.owner, c.repository)
}

// GiteaDriver provides access to the Gitea API.
type GiteaDriver struct {
	GiteaConfig
	client *gitea.Client
	log    logFn
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (d *GiteaDriver) apiMergePullRequest(pullRequestNumber int64, commitTitle, commitMessage string) (mergeSha string, err error) {
	_, err = d.client.MergePullRequest(d.owner, d.repository, pullRequestNumber, gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   commitTitle,
		Message: commitMessage,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := d.client.GetPullRequest(d.owner, d.repository, pullRequestNumber)
	if err != nil {
		return "", err
	}
	return *pullRequest.MergedCommitID, nil
}

// retargetPullRequests retargets pullrequests onto a new base branch
// this comes in handy when an ancestor got merged, so that children can be retargeted to the ancestor's own target branch
// example:
//
//	ancestor -> initial
//	children1 -> ancestor  --> initial (retargeted to initial after merge)
//	children2 -> ancestor  --> initial (retargeted to initial after merge)
//
//nolint:unparam
func (d *GiteaDriver) apiRetargetPullRequests(pullRequests []*gitea.PullRequest, newBaseName string) error {
	for _, pullRequest := range pullRequests {
		// RE-ENABLE AFTER https://github.com/go-gitea/gitea/issues/11552 and remove the nolint above
		// if options.LogRequests {
		// 	helpers.PrintLog(fmt.Sprintf("Gitea API: Updating base branch for PR #%d to #%s", *pullRequest.Index, newBaseName))
		// }
		d.log("Gitea API: Updating base branch for PR #%d to #%s", 1, newBaseName)
		d.log("The Gitea API currently does not support retargeting, please restarget #%d manually, see https://github.com/go-gitea/gitea/issues/11552", pullRequest.Index)
		// _, err = d.client.EditPullRequest(d.owner, d.repository, *pullRequest.Index, &gitea.EditPullRequestOption{
		// 	Base: newBaseName
		// })
		// if err != nil {
		//	return err
		// }
	}
	return nil
}

func (d *GiteaDriver) LoadPullRequestInfo(branch, parentBranch string) (*PullRequestInfo, error) {
	if d.apiToken == "" {
		return nil, nil //nolint:nilnil // we really want to return nil here
	}
	openPullRequests, err := d.client.ListRepoPullRequests(d.owner, d.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	baseName := parentBranch
	headName := d.owner + "/" + branch
	pullRequests := filterPullRequests(openPullRequests, baseName, headName)
	if len(pullRequests) < 1 {
		return nil, fmt.Errorf("no pull request from branch %q to branch %q found", branch, parentBranch)
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf("found %d pull requests from branch %q to branch %q", len(pullRequests), branch, parentBranch)
	}
	pullRequest := pullRequests[0]
	return &PullRequestInfo{
		CanMergeWithAPI:      pullRequest.Mergeable,
		DefaultCommitMessage: createDefaultCommitMessage(pullRequest),
		PullRequestNumber:    pullRequest.Index,
	}, nil
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (d *GiteaDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	openPullRequests, err := d.client.ListRepoPullRequests(d.owner, d.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return "", err
	}
	baseName := options.Branch
	newBaseName := options.ParentBranch
	err = d.apiRetargetPullRequests(filterPullRequests(openPullRequests, baseName, ""), newBaseName)
	if err != nil {
		return "", err
	}
	commitMessageParts := strings.SplitN(options.CommitMessage, "\n", 2)
	commitTitle := commitMessageParts[0]
	commitMessage := ""
	if len(commitMessageParts) == 2 {
		commitMessage = commitMessageParts[1]
	}
	return d.apiMergePullRequest(options.PullRequestNumber, commitTitle, commitMessage)
}

func createDefaultCommitMessage(pullRequest *gitea.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", pullRequest.Title, pullRequest.Index)
}

func filterPullRequests(pullRequests []*gitea.PullRequest, baseName, headName string) []*gitea.PullRequest {
	pullRequestsFiltered := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Base.Name != baseName {
			break
		}
		if headName != "" && pullRequest.Head.Name != headName {
			break
		}
		pullRequestsFiltered = append(pullRequestsFiltered, pullRequest)
	}
	return pullRequestsFiltered
}
