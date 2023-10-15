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

func (self *GiteaConnector) FindProposal(branch, target domain.LocalBranchName) (*Proposal, error) {
	openPullRequests, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	pullRequests := FilterGiteaPullRequests(openPullRequests, self.Organization, branch, target)
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
		Target:          domain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:           pullRequest.Title,
	}, nil
}

func (self *GiteaConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self *GiteaConnector) HostingServiceName() string {
	return "Gitea"
}

func (self *GiteaConnector) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	toCompare := parentBranch.String() + "..." + branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self *GiteaConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.Hostname, self.Organization, self.Repository)
}

func (self *GiteaConnector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
	if number <= 0 {
		return domain.EmptySHA(), fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	title, body := ParseCommitMessage(message)
	_, err = self.client.MergePullRequest(self.Organization, self.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return domain.EmptySHA(), err
	}
	pullRequest, err := self.client.GetPullRequest(self.Organization, self.Repository, int64(number))
	if err != nil {
		return domain.EmptySHA(), err
	}
	return domain.NewSHA(*pullRequest.MergedCommitID), nil
}

func (self *GiteaConnector) UpdateProposalTarget(_ int, _ domain.LocalBranchName) error {
	// TODO: update the client and uncomment
	// if self.log != nil {
	// 	self.log(message.HostingGiteaUpdateBasebranchViaAPI, number, target)
	// }
	// _, err := self.client.EditPullRequest(self.owner, self.repository, int64(number), gitea.EditPullRequestOption{
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

func FilterGiteaPullRequests(pullRequests []*gitea.PullRequest, organization string, branch, target domain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	headName := organization + "/" + branch.String()
	for p := range pullRequests {
		pullRequest := pullRequests[p]
		if pullRequest.Head.Name == headName && pullRequest.Base.Name == target.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}
