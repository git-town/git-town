package gitea

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/git/giturl"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"golang.org/x/oauth2"
)

type Connector struct {
	hostingdomain.Data
	APIToken Option[configdomain.GiteaToken]
	client   *gitea.Client
	log      print.Logger
}

func (self Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	if len(hostingdomain.ReadProposalOverride()) > 0 {
		return Some(self.findProposalViaOverride)
	}
	if self.APIToken.IsSome() {
		return Some(self.findProposalViaAPI)
	}
	return None[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)]()
}

func (self Connector) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _ gitdomain.ProposalTitle, _ gitdomain.ProposalBody) (string, error) {
	toCompare := parentBranch.String() + "..." + branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	if self.APIToken.IsSome() {
		return Some(self.searchProposal)
	}
	return None[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)]()
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error] {
	if self.APIToken.IsSome() {
		return Some(self.squashMergeProposal)
	}
	return None[func(number int, message gitdomain.CommitMessage) error]()
}

func (self Connector) UpdateProposalSourceFn() Option[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	if self.APIToken.IsSome() {
		return Some(self.updateProposalTarget)
	}
	return None[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error]()
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	pullRequests := FilterPullRequests(openPullRequests, branch, target)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	case 1:
		proposal := parsePullRequest(pullRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(proposal), nil
	default:
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	self.log.Ok()
	proposalURLOverride := hostingdomain.ReadProposalOverride()
	if proposalURLOverride == hostingdomain.OverrideNoProposal {
		return None[hostingdomain.Proposal](), nil
	}
	return Some(hostingdomain.Proposal{
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          proposalURLOverride,
	}), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	pullRequests := FilterPullRequests2(openPullRequests, branch)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	case 1:
		pullRequest := pullRequests[0]
		proposal := parsePullRequest(pullRequest)
		self.log.Success(proposal.Target.String())
		return Some(proposal), nil
	default:
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.HostingGithubMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
	_, _, err := self.client.MergePullRequest(self.Organization, self.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   commitMessageParts.Subject,
		Message: commitMessageParts.Text,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	self.log.Start(messages.APIProposalLookupStart)
	_, _, err = self.client.GetPullRequest(self.Organization, self.Repository, int64(number))
	self.log.Ok()
	return err
}

func (self Connector) updateProposalTarget(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(number), gitea.EditPullRequestOption{
		Base: targetName,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func FilterPullRequests(pullRequests []*gitea.PullRequest, branch, target gitdomain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() && pullRequest.Base.Name == target.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func FilterPullRequests2(pullRequests []*gitea.PullRequest, branch gitdomain.LocalBranchName) []*gitea.PullRequest {
	result := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() {
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

func parsePullRequest(pullRequest *gitea.PullRequest) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
		URL:          pullRequest.HTMLURL,
	}
}
