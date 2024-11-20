package bitbucketdatacenter

import (
	"context"
	"errors"
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"net/url"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	hostingdomain.Data
	log    print.Logger
	client *bitbucketv1.APIClient
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	client := bitbucketv1.NewAPIClient(
		context.WithValue(context.TODO(), bitbucketv1.ContextBasicAuth, bitbucketv1.BasicAuth{UserName: args.UserName.String(), Password: args.AppPassword.String()}),
		bitbucketv1.NewConfiguration(args.RemoteURL.Host),
	)
	return Connector{
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: client,
		log:    args.Log,
	}
}

type NewConnectorArgs struct {
	AppPassword     Option[configdomain.BitbucketAppPassword]
	HostingPlatform Option[configdomain.HostingPlatform]
	Log             print.Logger
	RemoteURL       giturl.Parts
	UserName        Option[configdomain.BitbucketUsername]
}

func (self Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	proposalURLOverride := hostingdomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		return Some(self.findProposalViaOverride)
	}
	return Some(self.findProposalViaAPI)
}

func (self Connector) NewProposalURL(branch, parentBranch, _ gitdomain.LocalBranchName, _ gitdomain.ProposalTitle, _ gitdomain.ProposalBody) (string, error) {
	return fmt.Sprintf("%s/pull-requests?create&sourceBranch=%s&targetBranch=%s",
			self.RepositoryURL(),
			url.QueryEscape(branch.String()),
			url.QueryEscape(parentBranch.String())),
		nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	return Some(self.searchProposal)
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error] {
	return Some(self.squashMergeProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	return Some(self.updateProposalSource)
}

func (self Connector) UpdateProposalTargetFn() Option[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	return Some(self.updateProposalTarget)
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)

	// TODO filter for source/target branch
	resp, err := self.client.DefaultApi.GetPullRequestsPage(
		self.Organization, self.Repository, nil,
	)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	prs, err := bitbucketv1.GetPullRequestsResponse(resp)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	size := len(prs)
	switch {
	case size == 0:
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	case size > 1:
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromToFound, size, branch, target))
		return None[hostingdomain.Proposal](), nil
	}

	proposal, err := parsePullRequest(prs[0], self.RepositoryURL())
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}

	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(proposal), nil
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := hostingdomain.ReadProposalOverride()
	self.log.Ok()
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
	self.log.Start(messages.APIProposalLookupStart)

	// TODO filter for source branch
	resp, err := self.client.DefaultApi.GetPullRequestsPage(
		self.Organization, self.Repository, nil,
	)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	prs, err := bitbucketv1.GetPullRequestsResponse(resp)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	size := len(prs)
	switch {
	case size == 0:
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	case size > 1:
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromFound, size, branch))
		return None[hostingdomain.Proposal](), nil
	}

	proposal, err := parsePullRequest(prs[0], self.RepositoryURL())
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}

	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(proposal), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingBitbucketMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	_, err := self.client.DefaultApi.Merge(
		self.Organization,
		self.Repository,
		number,
		nil,
		nil, // TODO include message
		nil,
	)
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self Connector) updateProposalSource(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error {
	// TODO
	sourceName := source.String()
	self.log.Start(messages.APIUpdateProposalSource, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(sourceName))
	self.log.Failed("unsupported operation")
	return nil
}

func (self Connector) updateProposalTarget(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	// TODO
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(targetName))
	self.log.Failed("unsupported operation")
	return nil
}

func parsePullRequest(pullRequest bitbucketv1.PullRequest, repoURL string) (result hostingdomain.Proposal, err error) {
	return hostingdomain.Proposal{
		MergeWithAPI: false,
		Number:       pullRequest.ID,
		Source:       gitdomain.NewLocalBranchName(pullRequest.FromRef.DisplayID),
		Target:       gitdomain.NewLocalBranchName(pullRequest.ToRef.DisplayID),
		Title:        pullRequest.Title,
		URL:          fmt.Sprintf("%s/pull-requests/%v/overview", repoURL, pullRequest),
	}, nil
}
