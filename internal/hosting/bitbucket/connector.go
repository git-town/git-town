package bitbucket

import (
	"errors"
	"fmt"
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
	client := bitbucket.NewBasicAuth("", args.AppPassword.String())
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
}

func (self Connector) CanMakeAPICalls() bool {
	return true
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
	result1, err := self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    "git-town-qa", // TODO
		RepoSlug: "test-repo",   // TODO
		Query:    fmt.Sprintf(`source.branch.name = "%s"`, branch),
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}
	if result1 == nil {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	result2, ok := result1.(map[string]interface{})
	if !ok {
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	size1, has := result2["size"]
	if !has {
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed("unexpected result data structure")
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
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	proposal2, ok := proposal1.([]interface{})
	if !ok {
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	if len(proposal2) == 0 {
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	proposal3, ok := proposal2[0].(map[string]interface{})
	if !ok {
		self.log.Failed("unexpected result data structure")
		return None[hostingdomain.Proposal](), nil
	}
	proposal, err := parsePullRequest(proposal3)
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), nil
	}
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.Number)), proposal.Title))
	return Some(proposal), nil
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

func (self Connector) SearchProposals(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	response1, err := self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    "git-town-qa",
		RepoSlug: "test-repo",
		Query:    fmt.Sprintf(`source.branch.name = "%s"`, branch),
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	response2, ok := response1.(map[string]interface{})
	if !ok {
		self.log.Failed("bitbucket API response has unknown structure")
		return None[hostingdomain.Proposal](), nil
	}
	size1, has := response2["size"]
	if !has {
		self.log.Failed("bitbucket API response has no size")
		return None[hostingdomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed("unknown size data type")
		return None[hostingdomain.Proposal](), nil
	}
	size := int(size2)
	if size == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if size > 1 {
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromFound, size, branch))
		return None[hostingdomain.Proposal](), nil
	}
	values1, has := response2["values"]
	if !has {
		self.log.Failed("bitbucket API has no values")
		return None[hostingdomain.Proposal](), nil
	}
	values2, ok := values1.([]interface{})
	if !ok {
		self.log.Failed("unknown data structure for values")
		return None[hostingdomain.Proposal](), nil
	}
	values3 := values2[0].(map[string]interface{})
	title1, has := values3["title"]
	if !has {
		self.log.Failed("no title field")
		return None[hostingdomain.Proposal](), nil
	}
	title2 := title1.(string)
	number1, has := values3["id"]
	if !has {
		self.log.Failed("no id field")
		return None[hostingdomain.Proposal](), nil
	}
	number2 := number1.(float64)
	number3 := int(number2)
	dest1, has := values3["destination"]
	if !has {
		self.log.Failed("no source field")
		return None[hostingdomain.Proposal](), nil
	}
	dest2, ok := dest1.(map[string]interface{})
	if !ok {
		self.log.Failed("unknown data type for source")
		return None[hostingdomain.Proposal](), nil
	}
	dest3, has := dest2["branch"]
	if !has {
		self.log.Failed("has no branch field")
		return None[hostingdomain.Proposal](), nil
	}
	dest4, ok := dest3.(map[string]interface{})
	if !ok {
		self.log.Failed("unknown data structure for branch field")
		return None[hostingdomain.Proposal](), nil
	}
	dest5, has := dest4["name"]
	if !has {
		self.log.Failed("has no name field")
		return None[hostingdomain.Proposal](), nil
	}
	dest6, ok := dest5.(string)
	if !ok {
		self.log.Failed("name is not a string")
		return None[hostingdomain.Proposal](), nil
	}
	dest7 := gitdomain.NewLocalBranchName(dest6)
	link1, has := values3["links"]
	if !has {
		self.log.Failed("no links attribute")
		return None[hostingdomain.Proposal](), nil
	}
	link2, ok := link1.(map[string]interface{})
	if !ok {
		self.log.Failed("unknown links structure")
		return None[hostingdomain.Proposal](), nil
	}
	link3, has := link2["html"]
	if !has {
		self.log.Failed("unknown html links")
		return None[hostingdomain.Proposal](), nil
	}
	link4, ok := link3.(map[string]interface{})
	if !ok {
		self.log.Failed("unknown html links structure")
		return None[hostingdomain.Proposal](), nil
	}
	link5, has := link4["href"]
	if !has {
		self.log.Failed("no href attribute")
		return None[hostingdomain.Proposal](), nil
	}
	url, ok := link5.(string)
	if !ok {
		self.log.Failed("href is not string")
		return None[hostingdomain.Proposal](), nil
	}
	proposal := hostingdomain.Proposal{
		MergeWithAPI: false,
		Number:       number3,
		Target:       dest7,
		Title:        title2,
		URL:          url,
	}
	self.log.Success(proposal.Target.String())
	return Some(proposal), nil
}

func (self Connector) SquashMergeProposal(_ int, _ gitdomain.CommitMessage) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}

func (self Connector) UpdateProposalBase(_ int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error {
	finalMessages.Add("The BitBucket driver does not support updating proposals yet.")
	return nil
}

func (self Connector) UpdateProposalHead(_ int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error {
	finalMessages.Add("The BitBucket driver does not support updating proposals yet.")
	return nil
}

func parsePullRequest(pullRequest map[string]interface{}) (result hostingdomain.Proposal, err error) {
	id1, has := pullRequest["id"]
	if !has {
		return result, errors.New("missing id attribute in proposal")
	}
	id2, ok := id1.(float64)
	if !ok {
		return result, errors.New("unknown data type for pull request title")
	}
	number := int(id2)
	titleRaw, has := pullRequest["title"]
	if !has {
		return result, errors.New("missing title attribute in proposal")
	}
	title, ok := titleRaw.(string)
	if !ok {
		return result, errors.New("unknown data type for pull request title")
	}
	destination1, has := pullRequest["destination"]
	if !has {
		return result, errors.New("missing destination attribute in proposal")
	}
	destination2, ok := destination1.(map[string]interface{})
	if !ok {
		return result, errors.New("unknown data structure for destination")
	}
	destination3, has := destination2["branch"]
	if !has {
		return result, errors.New("no branch attribute")
	}
	destination4, ok := destination3.(map[string]interface{})
	if !ok {
		return result, errors.New("unknown data structure for destination")
	}
	destination5, has := destination4["name"]
	if !has {
		return result, errors.New("no branch attribute")
	}
	destination, ok := destination5.(string)
	if !ok {
		return result, errors.New("unknown data structure for destination")
	}
	link1, has := pullRequest["links"]
	if !has {
		return result, errors.New("no links attribute")
	}
	link2, ok := link1.(map[string]interface{})
	if !ok {
		return result, errors.New("unknown links structure")
	}
	link3, has := link2["html"]
	if !has {
		return result, errors.New("unknown html links")
	}
	link4, ok := link3.(map[string]interface{})
	if !ok {
		return result, errors.New("unknown html links structure")
	}
	link5, has := link4["href"]
	if !has {
		return result, errors.New("no href attribute")
	}
	url, ok := link5.(string)
	if !ok {
		return result, errors.New("href is not string")
	}
	return hostingdomain.Proposal{
		MergeWithAPI: false,
		Number:       number,
		Target:       gitdomain.NewLocalBranchName(destination),
		Title:        title,
		URL:          url,
	}, nil
}
