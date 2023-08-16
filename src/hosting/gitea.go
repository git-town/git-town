package hosting

import (
	"context"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/messages"
	"golang.org/x/oauth2"
)

type GiteaConnector struct {
	client *gitea.Client
	CommonConfig
	log Log
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
		return nil, fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
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

func (c *GiteaConnector) SquashMergeProposal(number int, message string) (mergeSha domain.SHA, err error) {
	if number <= 0 {
		return domain.SHA{}, fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	title, body := ParseCommitMessage(message)
	_, err = c.client.MergePullRequest(c.Organization, c.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return domain.SHA{}, err
	}
	pullRequest, err := c.client.GetPullRequest(c.Organization, c.Repository, int64(number))
	if err != nil {
		return domain.SHA{}, err
	}
	return domain.NewSHA(*pullRequest.MergedCommitID), nil
}

func (c *GiteaConnector) UpdateProposalTarget(_ int, _ string) error {
	// TODO: update the client and uncomment
	// if c.log != nil {
	// 	c.log(message.HostingGiteaUpdateBasebranchViaAPI, number, target)
	// }
	// _, err := c.client.EditPullRequest(c.owner, c.repository, int64(number), gitea.EditPullRequestOption{
	// 	Base: newBaseName,
	// })
	// return err
	return fmt.Errorf(messages.HostingGiteaNotImplemented)
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewGiteaConnector(args NewGiteaConnectorArgs) (*GiteaConnector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "gitea.com" && args.HostingService != config.HostingGitea) {
		return nil, nil //nolint:nilnil
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", args.OriginURL.Host), httpClient)
	return &GiteaConnector{
		client: giteaClient,
		CommonConfig: CommonConfig{
			APIToken:     args.APIToken,
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		log: args.Log,
	}, nil
}

type NewGiteaConnectorArgs struct {
	OriginURL      *giturl.Parts
	HostingService config.Hosting
	APIToken       string
	Log            Log
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
