package codeberg

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v19/internal/cli/colors"
	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/forge/forgedomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/git/giturl"
	"github.com/git-town/git-town/v19/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// Connector provides standardized connectivity for the given repository (codeberg.org/owner/repo)
// via the Codeberg API.
type Connector struct {
	forgedomain.Data
	APIToken Option[configdomain.CodebergToken]
	client   *forgejo.Client
	log      print.Logger
}

func (self Connector) DefaultProposalMessage(proposal forgedomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
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

func (self Connector) UpdateProposalSourceFn() Option[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	if self.APIToken.IsSome() {
		return Some(self.updateProposalTarget)
	}
	return None[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error]()
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 50,
		},
		State: forgejo.StateOpen,
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
		return Some(proposal), nil
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
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          proposalURLOverride,
	}), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 50,
		},
		State: forgejo.StateOpen,
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
		return Some(proposal), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.ForgeCodebergMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
	_, _, err := self.client.MergePullRequest(self.Organization, self.Repository, int64(number), forgejo.MergePullRequestOption{
		Style:   forgejo.MergeStyleSquash,
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
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(number), forgejo.EditPullRequestOption{
		Base: targetName,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func FilterPullRequests(pullRequests []*forgejo.PullRequest, branch, target gitdomain.LocalBranchName) []*forgejo.PullRequest {
	result := []*forgejo.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() && pullRequest.Base.Name == target.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func FilterPullRequests2(pullRequests []*forgejo.PullRequest, branch gitdomain.LocalBranchName) []*forgejo.PullRequest {
	result := []*forgejo.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

// NewConnector provides a new connector instance.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	codebergClient, err := forgejo.NewClient("https://"+args.RemoteURL.Host, forgejo.SetToken(args.APIToken.String()))
	return Connector{
		APIToken: args.APIToken,
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: codebergClient,
		log:    args.Log,
	}, err
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.CodebergToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

func parsePullRequest(pullRequest *forgejo.PullRequest) forgedomain.Proposal {
	return forgedomain.Proposal{
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
		URL:          pullRequest.HTMLURL,
	}
}
