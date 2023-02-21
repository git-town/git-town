package hosting

import (
	"context"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v7/src/giturl"
	"golang.org/x/oauth2"
)

type GiteaConnector struct {
	client *gitea.Client
	giteaConfig
	log logFn
}

func (c *GiteaConnector) ChangeRequestForBranch(branch string) (*ChangeRequestInfo, error) {
	openPullRequests, err := c.client.ListRepoPullRequests(c.owner, c.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	headName := c.owner + "/" + branch
	pullRequests := filterPullRequests(openPullRequests, headName)
	if len(pullRequests) == 0 {
		return nil, nil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf("found %d pull requests for branch %q", len(pullRequests), branch)
	}
	pullRequest := pullRequests[0]
	return &ChangeRequestInfo{
		CanMergeWithAPI: pullRequest.Mergeable,
		Number:          int(pullRequest.Index),
		Title:           pullRequest.Title,
	}, nil
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (c *GiteaConnector) SquashMergeChangeRequest(number int, message string) (mergeSha string, err error) {
	title, body := parseCommitMessage(message)
	_, err = c.client.MergePullRequest(c.owner, c.repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := c.client.GetPullRequest(c.owner, c.repository, int64(number))
	if err != nil {
		return "", err
	}
	return *pullRequest.MergedCommitID, nil
}

func (c *GiteaConnector) UpdateChangeRequestTarget(number int, target string) error {
	// TODO: update the client and uncomment
	// if c.log != nil {
	// 	c.log("Gitea API: Updating base branch for PR #%d to #%s", number, target)
	// }
	// _, err := c.client.EditPullRequest(c.owner, c.repository, int64(number), gitea.EditPullRequestOption{
	// 	Base: newBaseName,
	// })
	// return err
	return fmt.Errorf("Updating Gitea pull requests is currently not supported")
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewGiteaConnector(url giturl.Parts, config gitConfig, log logFn) *GiteaConnector {
	manualHostName := config.OriginOverride()
	if config.HostingService() != "gitea" && manualHostName != "gitea.com" {
		return nil
	}
	if manualHostName != "" {
		url.Host = manualHostName
	}
	giteaConfig := giteaConfig{
		Config: Config{
			apiToken:   config.GiteaToken(),
			hostname:   url.Host,
			originURL:  config.OriginURL(),
			owner:      url.Org,
			repository: url.Repo,
		},
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: giteaConfig.apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", giteaConfig.hostname), httpClient)
	return &GiteaConnector{
		client:      giteaClient,
		giteaConfig: giteaConfig,
		log:         log,
	}
}

// GiteaConfig contains connection information for Gitea-based hosting platforms.
type giteaConfig struct {
	Config
}

func (c *giteaConfig) HostingServiceName() string {
	return "Gitea"
}

func (c *giteaConfig) NewChangeRequestURL(branch, parentBranch string) (string, error) {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *giteaConfig) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.owner, c.repository)
}

func (c *giteaConfig) DefaultCommitMessage(crInfo ChangeRequestInfo) string {
	return fmt.Sprintf("%s (#%d)", crInfo.Title, crInfo.Number)
}

func filterPullRequests(pullRequests []*gitea.PullRequest, branch string) []*gitea.PullRequest {
	pullRequestsFiltered := []*gitea.PullRequest{}
	// TODO: don't copy the entire pullRequest struct here, use the index
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch {
			pullRequestsFiltered = append(pullRequestsFiltered, pullRequest)
		}
	}
	return pullRequestsFiltered
}
