package hosting

import (
	"context"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v7/src/config"
	"golang.org/x/oauth2"
)

type GiteaConnector struct {
	client *gitea.Client
	CommonConfig
	log logFn
}

func (c *GiteaConnector) FindProposal(branch, target string) (*Proposal, error) {
	openPullRequests, err := c.client.ListRepoPullRequests(c.Organization, c.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	pullRequests := FilterGiteaPullRequests(openPullRequests, c.Organization, branch, target)
	if len(pullRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf("found %d pull requests for branch %q", len(pullRequests), branch)
	}
	pullRequest := pullRequests[0]
	return &Proposal{
		CanMergeWithAPI: pullRequest.Mergeable,
		Number:          int(pullRequest.Index),
		Target:          pullRequest.Base.Ref,
		Title:           pullRequest.Title,
	}, nil
}

func (c *GiteaConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *GiteaConnector) HostingServiceName() string {
	return "Gitea"
}

func (c *GiteaConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *GiteaConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.Hostname, c.Organization, c.Repository)
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (c *GiteaConnector) SquashMergeProposal(number int, message string) (mergeSha string, err error) {
	if number <= 0 {
		return "", fmt.Errorf("no pull request number given")
	}
	title, body := ParseCommitMessage(message)
	_, err = c.client.MergePullRequest(c.Organization, c.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := c.client.GetPullRequest(c.Organization, c.Repository, int64(number))
	if err != nil {
		return "", err
	}
	return *pullRequest.MergedCommitID, nil
}

func (c *GiteaConnector) UpdateProposalTarget(number int, target string) error {
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
func NewGiteaConnector(gitConfig gitTownConfig, log logFn) (*GiteaConnector, error) {
	hostingService, err := gitConfig.HostingService()
	if err != nil {
		return nil, err
	}
	url := gitConfig.OriginURL()
	if url == nil || (url.Host != "gitea.com" && hostingService != config.HostingServiceGitea) {
		return nil, nil //nolint:nilnil
	}
	apiToken := gitConfig.GiteaToken()
	hostname := url.Host
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", hostname), httpClient)
	return &GiteaConnector{
		client: giteaClient,
		CommonConfig: CommonConfig{
			APIToken:     apiToken,
			Hostname:     hostname,
			Organization: url.Org,
			Repository:   url.Repo,
		},
		log: log,
	}, nil
}

func FilterGiteaPullRequests(pullRequests []*gitea.PullRequest, organization, branch, target string) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	headName := organization + "/" + branch
	for p := range pullRequests {
		pullRequest := pullRequests[p]
		if pullRequest.Head.Name == headName && pullRequest.Base.Name == target {
			result = append(result, pullRequest)
		}
	}
	return result
}
