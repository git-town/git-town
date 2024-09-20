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
	result, err := self.client.Repositories.PullRequests.Get(&bitbucket.PullRequestsOptions{
		Owner:             "git-town-qa",
		RepoSlug:          "test-repo",
		SourceBranch:      branch.String(),
		DestinationBranch: target.String(),
		// States: []string{},
	})
	if err != nil {
		self.log.Failed(err)
		return None[hostingdomain.Proposal](), nil
	}
	if result == nil {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	keyValues, ok := result.(map[string]interface{})
	if !ok {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	sizeRaw, has := keyValues["size"]
	if !has {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	sizeFloat, ok := sizeRaw.(float64)
	if !ok {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	sizeInt := int(sizeFloat)
	if sizeInt == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if sizeInt > 1 {
		self.log.Failed(fmt.Errorf(messages.ProposalMultipleFromToFound, sizeInt, branch, target))
		return None[hostingdomain.Proposal](), nil
	}
	valuesRaw1, has := keyValues["values"]
	if !has {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	valuesRaw2, ok := valuesRaw1.([]interface{})
	if !ok {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	if len(valuesRaw2) == 0 {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	valuesRaw3, ok := valuesRaw2[0].(map[string]interface{})
	if !ok {
		self.log.Failed(errors.New("unexpected result data structure"))
		return None[hostingdomain.Proposal](), nil
	}
	proposal, err := parsePullRequest(valuesRaw3)
	if err != nil {
		self.log.Failed(err)
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

func (self Connector) SearchProposals(_ gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	return None[hostingdomain.Proposal](), nil
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
