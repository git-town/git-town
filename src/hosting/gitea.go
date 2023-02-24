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
	CommonConfig
	log logFn
}

func (c *GiteaConnector) FindProposal(branch, target string) (*Proposal, error) {
	openPullRequests, err := c.client.ListRepoPullRequests(c.organization, c.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	pullRequests := FilterGiteaPullRequests(openPullRequests, c.organization, branch, target)
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
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.organization, c.repository)
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (c *GiteaConnector) SquashMergeProposal(number int, message string) (mergeSha string, err error) {
	if number <= 0 {
		return "", fmt.Errorf("no pull request number given")
	}
	title, body := ParseCommitMessage(message)
	_, err = c.client.MergePullRequest(c.organization, c.repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := c.client.GetPullRequest(c.organization, c.repository, int64(number))
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
func NewGiteaConnector(url giturl.Parts, config gitConfig, log logFn) *GiteaConnector {
	hostingService := config.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if hostingService != "gitea" && url.Host != "gitea.com" {
		return nil
	}
	apiToken := config.GiteaToken()
	hostname := url.Host
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", hostname), httpClient)
	return &GiteaConnector{
		client: giteaClient,
		CommonConfig: CommonConfig{
			apiToken:     apiToken,
			hostname:     hostname,
			originURL:    config.OriginURL(),
			organization: url.Org,
			repository:   url.Repo,
		},
		log: log,
	}
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
