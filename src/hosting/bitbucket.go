package hosting

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/messages"
)

// BitbucketConnector provides access to the API of Bitbucket installations.
type BitbucketConnector struct {
	CommonConfig
	getSHAForBranch SHAForBranchFunc
}

// NewBitbucketConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewBitbucketConnector(args NewBitbucketConnectorArgs) (*BitbucketConnector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "bitbucket.org" && args.HostingService != config.HostingBitbucket) {
		return nil, nil //nolint:nilnil
	}
	return &BitbucketConnector{
		CommonConfig: CommonConfig{
			APIToken:     "",
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		getSHAForBranch: args.GetSHAForBranch,
	}, nil
}

type NewBitbucketConnectorArgs struct {
	OriginURL       *giturl.Parts
	HostingService  config.Hosting
	GetSHAForBranch SHAForBranchFunc
}

func (bc *BitbucketConnector) FindProposal(_, _ domain.LocalBranchName) (*Proposal, error) {
	return nil, fmt.Errorf(messages.HostingBitBucketNotImplemented)
}

func (bc *BitbucketConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (bc *BitbucketConnector) HostingServiceName() string {
	return "Bitbucket"
}

func (bc *BitbucketConnector) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	query := url.Values{}
	branchSHA, err := bc.getSHAForBranch(branch.BranchName())
	if err != nil {
		return "", fmt.Errorf(messages.ProposalURLProblem, branch, parentBranch, err)
	}
	query.Add("source", strings.Join([]string{bc.Organization + "/" + bc.Repository, branchSHA.TruncateTo(12).String(), branch.String()}, ":"))
	query.Add("dest", strings.Join([]string{bc.Organization + "/" + bc.Repository, "", parentBranch.String()}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", bc.RepositoryURL(), query.Encode()), nil
}

func (bc *BitbucketConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", bc.Hostname, bc.Organization, bc.Repository)
}

func (bc *BitbucketConnector) SquashMergeProposal(_ int, _ string) (mergeSHA domain.SHA, err error) {
	return domain.SHA{}, errors.New(messages.HostingBitBucketNotImplemented)
}

func (bc *BitbucketConnector) UpdateProposalTarget(_ int, _ domain.LocalBranchName) error {
	return errors.New(messages.HostingBitBucketNotImplemented)
}
