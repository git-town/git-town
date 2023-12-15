package gitea

import (
	"context"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/git-town/git-town/v11/src/messages"
	"golang.org/x/oauth2"
)

type Connector struct {
	client   *gitea.Client
	APIToken configdomain.GiteaToken
	common.Config
	log common.Log
}

func (self *Connector) DefaultProposalMessage(proposal domain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self *Connector) FindProposal(branch, target domain.LocalBranchName) (*domain.Proposal, error) {
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return nil, err
	}
	pullRequests := FilterPullRequests(openPullRequests, self.Organization, branch, target)
	if len(pullRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
	}
	pullRequest := pullRequests[0]
	return &domain.Proposal{
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Target:       domain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
	}, nil
}

func (self *Connector) HostingServiceName() string {
	return "Gitea"
}

func (self *Connector) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	toCompare := parentBranch.String() + "..." + branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self *Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self *Connector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
	if number <= 0 {
		return domain.EmptySHA(), fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	title, body := common.CommitMessageParts(message)
	_, _, err = self.client.MergePullRequest(self.Organization, self.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   title,
		Message: body,
	})
	if err != nil {
		return domain.EmptySHA(), err
	}
	pullRequest, _, err := self.client.GetPullRequest(self.Organization, self.Repository, int64(number))
	if err != nil {
		return domain.EmptySHA(), err
	}
	return domain.NewSHA(*pullRequest.MergedCommitID), nil
}

func (self *Connector) UpdateProposalTarget(_ int, _ domain.LocalBranchName) error {
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

func FilterPullRequests(pullRequests []*gitea.PullRequest, organization string, branch, target domain.LocalBranchName) []*gitea.PullRequest {
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

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "gitea.com" && args.HostingService != configdomain.HostingGitea) {
		return nil, nil //nolint:nilnil
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", args.OriginURL.Host), httpClient)
	return &Connector{
		APIToken: args.APIToken,
		client:   giteaClient,
		Config: common.Config{
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		log: args.Log,
	}, nil
}

type NewConnectorArgs struct {
	OriginURL      *giturl.Parts
	HostingService configdomain.Hosting
	APIToken       configdomain.GiteaToken
	Log            common.Log
}
