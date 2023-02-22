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
	Config
	log logFn
}

func (c *GiteaConnector) ChangeRequestDetails(number int) (*ChangeRequestInfo, error) {
	return nil, fmt.Errorf("TODO: implement")
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
		return nil, nil //nolint:nilnil
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

func (c *GiteaConnector) ChangeRequests() ([]ChangeRequestInfo, error) {
	pullRequests, err := c.client.ListRepoPullRequests(c.owner, c.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	result := make([]ChangeRequestInfo, len(pullRequests))
	if err != nil {
		return result, err
	}
	for p := range pullRequests {
		result[p] = parseGiteaPullRequest(pullRequests[p])
	}
	return result, nil
}

func (c *GiteaConnector) DefaultCommitMessage(crInfo ChangeRequestInfo) string {
	return fmt.Sprintf("%s (#%d)", crInfo.Title, crInfo.Number)
}

func (c *GiteaConnector) HostingServiceName() string {
	return "Gitea"
}

func (c *GiteaConnector) NewChangeRequestURL(branch, parentBranch string) (string, error) {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *GiteaConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.owner, c.repository)
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
	return fmt.Errorf("updating Gitea pull requests is currently not supported")
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewGiteaConnector(url giturl.Parts, config gitConfig, log logFn) *GiteaConnector {
	hostingService := config.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if hostingService != "gitea" && url.Host != "gitea.com" {
		return nil
	}
	hostingConfig := Config{
		apiToken:   config.GiteaToken(),
		hostname:   url.Host,
		originURL:  config.OriginURL(),
		owner:      url.Org,
		repository: url.Repo,
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: hostingConfig.apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", hostingConfig.hostname), httpClient)
	return &GiteaConnector{
		client: giteaClient,
		Config: hostingConfig,
		log:    log,
	}
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

func parseGiteaPullRequest(pullRequest *gitea.PullRequest) ChangeRequestInfo {
	return ChangeRequestInfo{
		CanMergeWithAPI: pullRequest.Mergeable,
		Number:          int(pullRequest.Index),
		Title:           pullRequest.Title,
	}
}
