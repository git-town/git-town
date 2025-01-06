package bitbucketcloud

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/git-town/git-town/v17/internal/cli/colors"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/git/giturl"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

// Connector provides access to the API of Bitbucket installations.
type Connector struct {
	hostingdomain.Data
	client *bitbucket.Client
	log    print.Logger
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) Connector {
	client := bitbucket.NewBasicAuth(args.UserName.String(), args.AppPassword.String())
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
	return fmt.Sprintf("%s/pull-requests/new?source=%s&dest=%s%%2F%s%%3A%s",
			self.RepositoryURL(),
			url.QueryEscape(branch.String()),
			url.QueryEscape(self.Organization),
			url.QueryEscape(self.Repository),
			url.QueryEscape(parentBranch.String())),
		nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
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
	query := fmt.Sprintf("source.branch.name = %q AND destination.branch.name = %q", branch, target)
	result1, err := self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    query,
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	if result1 == nil {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	result2, ok := result1.(map[string]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size1, has := result2["size"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size := int(size2)
	if size == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if size > 1 {
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromToFound, size, branch, target))
		return None[hostingdomain.Proposal](), nil
	}
	proposal1, has := result2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	proposal2, ok := proposal1.([]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	if len(proposal2) == 0 {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	proposal3, ok := proposal2[0].(map[string]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	proposal4, err := parsePullRequest(proposal3)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}
	self.log.Success(fmt.Sprintf("#%d", proposal4.Number))
	return Some(proposal4), nil
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
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	response1, err := self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    fmt.Sprintf("source.branch.name = %q", branch),
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	response2, ok := response1.(map[string]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size1, has := response2["size"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	size3 := int(size2)
	if size3 == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if size3 > 1 {
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromFound, size3, branch))
		return None[hostingdomain.Proposal](), nil
	}
	values1, has := response2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	values2, ok := values1.([]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[hostingdomain.Proposal](), nil
	}
	proposal1 := values2[0].(map[string]interface{})
	proposal2, err := parsePullRequest(proposal1)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}
	self.log.Success(proposal2.Target.String())
	return Some(proposal2), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingBitbucketMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	_, err := self.client.Repositories.PullRequests.Merge(&bitbucket.PullRequestsOptions{
		ID:       strconv.Itoa(number),
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Message:  message.String(),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self Connector) updateProposalSource(number int, source gitdomain.LocalBranchName, _ stringslice.Collector) error {
	self.log.Start(messages.APIUpdateProposalSource, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(source.String()))
	_, err := self.client.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:           strconv.Itoa(number),
		Owner:        self.Organization,
		RepoSlug:     self.Repository,
		SourceBranch: source.String(),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self Connector) updateProposalTarget(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(targetName))
	_, err := self.client.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:                strconv.Itoa(number),
		Owner:             self.Organization,
		RepoSlug:          self.Repository,
		DestinationBranch: target.String(),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func parsePullRequest(pullRequest map[string]interface{}) (result hostingdomain.Proposal, err error) {
	id1, has := pullRequest["id"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	id2, ok := id1.(float64)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	number := int(id2)
	title1, has := pullRequest["title"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	title2, ok := title1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination1, has := pullRequest["destination"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination2, ok := destination1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination3, has := destination2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination4, ok := destination3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination5, has := destination4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination6, ok := destination5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source1 := pullRequest["source"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source2, ok := source1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source3, has := source2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source4, ok := source3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source5, has := source4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source6, ok := source5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url1, has := pullRequest["links"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url2, ok := url1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url3, has := url2["html"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url4, ok := url3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url5, has := url4["href"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url6, ok := url5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	return hostingdomain.Proposal{
		MergeWithAPI: false,
		Number:       number,
		Source:       gitdomain.NewLocalBranchName(source6),
		Target:       gitdomain.NewLocalBranchName(destination6),
		Title:        title2,
		URL:          url6,
	}, nil
}
