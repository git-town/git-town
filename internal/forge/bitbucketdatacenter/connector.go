package bitbucketdatacenter

import (
	"context"
	"fmt"
	"net/url"

	"github.com/carlmjohnson/requests"
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	forgedomain.Data
	log      print.Logger
	token    string
	username string
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	return Connector{
		Data: forgedomain.Data{
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
	AppPassword Option[forgedomain.BitbucketAppPassword]
	ForgeType   Option[forgedomain.ForgeType]
	Log         print.Logger
	RemoteURL   giturl.Parts
	UserName    Option[forgedomain.BitbucketUsername]
}

// Connector implementation
var bbdcConnector Connector                 // type-check
var _ forgedomain.Connector = bbdcConnector // type-check

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	browser.Open(self.NewProposalURL(data), data.FrontendRunner)
	return nil
}

func (self Connector) NewProposalURL(data forgedomain.CreateProposalArgs) string {
	return fmt.Sprintf("%s/pull-requests?create&sourceBranch=%s&targetBranch=%s",
		self.RepositoryURL(),
		url.QueryEscape(data.Branch.String()),
		url.QueryEscape(data.ParentBranch.String()))
}

func (self Connector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	data := proposalData.Data()
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/projects/%s/repos/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

// ProposalFinder implementation
var _ forgedomain.ProposalFinder = bbdcConnector // type-check

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	proposalURLOverride := forgedomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		return self.findProposalViaOverride(branch, target)
	}
	return self.findProposalViaAPI(branch, target)
}

func (self Connector) apiBaseURL() string {
	return fmt.Sprintf(
		"https://%s/rest/api/latest/projects/%s/repos/%s/pull-requests",
		self.Hostname,
		self.Organization,
		self.Repository,
	)
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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
		return None[forgedomain.Proposal](), err
	}
	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
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
		return None[forgedomain.Proposal](), nil
	}
	proposal := parsePullRequest(*needle, self.RepositoryURL())
	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := forgedomain.ReadProposalOverride()
	self.log.Ok()
	if proposalURLOverride == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	data := forgedomain.ProposalData{
		Body:         None[string](),
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          proposalURLOverride,
	}
	return Some(forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}

// ProposalSearcher implementation
var _ forgedomain.ProposalSearcher = bbdcConnector

func (self Connector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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
		return None[forgedomain.Proposal](), err
	}
	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
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
		return None[forgedomain.Proposal](), nil
	}
	proposal := parsePullRequest(*needle, self.RepositoryURL())
	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}

func parsePullRequest(pullRequest PullRequest, repoURL string) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: false,
		Number:       pullRequest.ID,
		Source:       gitdomain.NewLocalBranchName(pullRequest.FromRef.DisplayID),
		Target:       gitdomain.NewLocalBranchName(pullRequest.ToRef.DisplayID),
		Title:        pullRequest.Title,
		Body:         NewOption(pullRequest.Description),
		URL:          fmt.Sprintf("%s/pull-requests/%v/overview", repoURL, pullRequest.ID),
	}
}
