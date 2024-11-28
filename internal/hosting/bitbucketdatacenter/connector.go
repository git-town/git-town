package bitbucketdatacenter

import (
	"context"
	"fmt"
	"net/url"

	"github.com/carlmjohnson/requests"
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
	log      print.Logger
	token    string
	username string
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	return Connector{
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		log:      args.Log,
		token:    args.AppPassword.String(),
		username: args.UserName.String(),
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
	return None[func(number int, message gitdomain.CommitMessage) error]()
}

func (self Connector) UpdateProposalSourceFn() Option[func(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	return None[func(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	return None[func(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error]()
}

func (self Connector) apiBaseURL() string {
	return fmt.Sprintf(
		"https://%s/rest/api/latest/projects/%s/repos/%s/pull-requests",
		self.Hostname,
		self.Organization,
		self.Repository,
	)
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)

	ctx := context.TODO()

	fromRefID := fmt.Sprintf("refs/heads/%v", branch)
	toRefID := fmt.Sprintf("refs/heads/%v", target)

	var resp PullRequestResponse

	err := requests.URL(self.apiBaseURL()).
		BasicAuth(self.username, self.token).
		Param("at", toRefID).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}

	var needle *PullRequest

	for _, pr := range resp.Values {
		if pr.FromRef.ID == fromRefID && pr.ToRef.ID == toRefID {
			needle = &pr
			break
		}
	}

	if needle == nil {
		self.log.Success("no PR found matching source and target branch")
		return None[hostingdomain.Proposal](), nil
	}

	proposal := parsePullRequest(*needle, self.RepositoryURL())

	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(proposal), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)

	ctx := context.TODO()

	fromRefID := fmt.Sprintf("refs/heads/%v", branch)

	var resp PullRequestResponse

	err := requests.URL(self.apiBaseURL()).
		BasicAuth(self.username, self.token).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}

	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}

	var needle *PullRequest

	for _, pr := range resp.Values {
		if pr.FromRef.ID == fromRefID {
			needle = &pr
			break
		}
	}

	if needle == nil {
		self.log.Success("no PR found matching source branch")
		return None[hostingdomain.Proposal](), nil
	}

	proposal := parsePullRequest(*needle, self.RepositoryURL())

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

func parsePullRequest(pullRequest PullRequest, repoURL string) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		MergeWithAPI: false,
		Number:       pullRequest.ID,
		Source:       gitdomain.NewLocalBranchName(pullRequest.FromRef.DisplayID),
		Target:       gitdomain.NewLocalBranchName(pullRequest.ToRef.DisplayID),
		Title:        pullRequest.Title,
		URL:          fmt.Sprintf("%s/pull-requests/%v/overview", repoURL, pullRequest),
	}
}
