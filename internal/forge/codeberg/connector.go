package codeberg

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides standardized connectivity for the given repository (codeberg.org/owner/repo)
// via the Codeberg API.
type Connector struct {
	forgedomain.Data
	APIToken Option[configdomain.CodebergToken]
	client   *forgejo.Client
	log      print.Logger
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	toCompare := data.ParentBranch.String() + "..." + data.Branch.String()
	url := fmt.Sprintf("%s/compare/%s", self.RepositoryURL(), url.PathEscape(toCompare))
	browser.Open(url, data.FrontendRunner)
	return nil
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return forgedomain.CommitBody(data, fmt.Sprintf("%s (#%d)", data.Title, data.Number))
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

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	browser.Open(self.RepositoryURL(), runner)
	return nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SearchProposalFn() Option[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if self.APIToken.IsSome() {
		return Some(self.searchProposal)
	}
	return None[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
}

func (self Connector) SquashMergeProposalFn() Option[func(int, gitdomain.CommitMessage) error] {
	if self.APIToken.IsSome() {
		return Some(self.squashMergeProposal)
	}
	return None[func(int, gitdomain.CommitMessage) error]()
}

func (self Connector) UpdateProposalSourceFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName) error] {
	return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName) error] {
	if self.APIToken.IsSome() {
		return Some(self.updateProposalTarget)
	}
	return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName) error]()
}

func (self Connector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, _, err := self.client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyConnectionResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 1,
		},
	})
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
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
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeCodeberg}), nil
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
		ForgeType: forgedomain.ForgeTypeCodeberg,
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
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeCodeberg}), nil
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

func (self Connector) updateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), forgejo.EditPullRequestOption{
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

func parsePullRequest(pullRequest *forgejo.PullRequest) forgedomain.ProposalData {
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
