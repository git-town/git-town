package gitea

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v20/internal/cli/colors"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/git/giturl"
	"github.com/git-town/git-town/v20/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"golang.org/x/oauth2"
)

type Connector struct {
	forgedomain.Data
	APIToken Option[configdomain.GiteaToken]
	client   *gitea.Client
	log      print.Logger
}

func (self Connector) DefaultProposalMessage(proposal forgedomain.Proposal) string {
	return forgedomain.CommitBody(proposal.Data, fmt.Sprintf("%s (#%d)", proposal.Data.GetTitle(), proposal.Data.GetNumber()))
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if len(forgedomain.ReadProposalOverride()) > 0 {
		return Some(self.findProposalViaOverride)
	}
	if self.APIToken.IsSome() {
		return Some(self.findProposalViaAPI)
	}
	return None[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
}

func (self Connector) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _ gitdomain.ProposalTitle, _ gitdomain.ProposalBody) (string, error) {
	toCompare := parentBranch.String() + "..." + branch.String()
	return fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if self.APIToken.IsSome() {
		return Some(self.searchProposal)
	}
	return None[func(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error] {
	if self.APIToken.IsSome() {
		return Some(self.squashMergeProposal)
	}
	return None[func(number int, message gitdomain.CommitMessage) error]()
}

func (self Connector) UpdateProposalSourceFn() Option[func(proposal forgedomain.Proposal, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(proposal forgedomain.Proposal, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(proposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	if self.APIToken.IsSome() {
		return Some(self.updateProposalTarget)
	}
	return None[func(proposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error]()
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	pullRequests := FilterPullRequests(openPullRequests, branch, target)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parsePullRequest(pullRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitea}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	self.log.Ok()
	proposalURLOverride := forgedomain.ReadProposalOverride()
	if proposalURLOverride == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         None[string](),
			MergeWithAPI: true,
			Number:       123,
			Source:       branch,
			Target:       target,
			Title:        "title",
			URL:          proposalURLOverride,
		},
		ForgeType: forgedomain.ForgeTypeGitea,
	}), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	pullRequests := FilterPullRequests2(openPullRequests, branch)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		pullRequest := pullRequests[0]
		proposal := parsePullRequest(pullRequest)
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitea}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.ForgeGithubMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
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

func (self Connector) updateProposalTarget(proposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.Data.GetNumber())), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(proposal.Data.GetNumber()), gitea.EditPullRequestOption{
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
		Data: forgedomain.Data{
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

func parsePullRequest(pullRequest *gitea.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
		Body:         NewOption(pullRequest.Body),
		URL:          pullRequest.HTMLURL,
	}
}
