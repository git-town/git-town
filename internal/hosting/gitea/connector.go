package gitea

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"golang.org/x/oauth2"
)

type Connector struct {
	hostingdomain.Data
	APIToken Option[configdomain.GiteaToken]
	client   *gitea.Client
	log      print.Logger
}

func (self Connector) CanMakeAPICalls() bool {
	return self.APIToken.IsSome() || len(hostingdomain.ReadProposalOverride()) > 0
}

func (self Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := hostingdomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		self.log.Ok()
		if proposalURLOverride == hostingdomain.OverrideNoProposal {
			return None[hostingdomain.Proposal](), nil
		}
		return Some(hostingdomain.Proposal{
			MergeWithAPI: true,
			Number:       123,
			Target:       target,
			Title:        "title",
			URL:          proposalURLOverride,
		}), nil
	}
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err)
		return None[hostingdomain.Proposal](), err
	}
	self.log.Ok()
	pullRequests := FilterPullRequests(openPullRequests, self.Organization, branch, target)
	switch len(pullRequests) {
	case 0:
		return None[hostingdomain.Proposal](), nil
	case 1:
		pullRequest := pullRequests[0]
		return Some(hostingdomain.Proposal{
			MergeWithAPI: pullRequest.Mergeable,
			Number:       int(pullRequest.Index),
			Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
			Title:        pullRequest.Title,
			URL:          pullRequest.HTMLURL,
		}), nil
	default:
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
	}
}

func (self Connector) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _ gitdomain.ProposalTitle, _ gitdomain.ProposalBody) (string, error) {
	toCompare := parentBranch.String() + "..." + branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.HostingGithubMergingViaAPI, number)
	_, _, err := self.client.MergePullRequest(self.Organization, self.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   commitMessageParts.Subject,
		Message: commitMessageParts.Text,
	})
	if err != nil {
		self.log.Failed(err)
		return err
	}
	self.log.Ok()
	self.log.Start(messages.APIProposalLookupStart)
	_, _, err = self.client.GetPullRequest(self.Organization, self.Repository, int64(number))
	self.log.Ok()
	return err
}

func (self Connector) UpdateProposalTarget(_ int, _ gitdomain.LocalBranchName) error {
	// if self.log != nil {
	// 	self.log(message.HostingGiteaUpdateBasebranchViaAPI, number, target)
	// }
	// _, err := self.client.EditPullRequest(self.owner, self.repository, int64(number), gitea.EditPullRequestOption{
	// 	Base: newBaseName,
	// })
	// return err
	return errors.New(messages.HostingGiteaNotImplemented)
}

func FilterPullRequests(pullRequests []*gitea.PullRequest, organization string, branch, target gitdomain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest(nil)
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
func NewConnector(args NewConnectorArgs) Connector {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP("https://"+args.RemoteURL.Host, httpClient)
	return Connector{
		APIToken: args.APIToken,
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: giteaClient,
		log:    args.Log,
	}
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.GiteaToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}
