package gitlab

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type Connector struct {
	client *gitlab.Client
	Data
	log print.Logger
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if len(forgedomain.ReadProposalOverride()) > 0 {
		return Some(self.findProposalViaOverride)
	}
	if self.Data.APIToken.IsSome() {
		return Some(self.findProposalViaAPI)
	}
	return None[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
}

func (self Connector) SearchProposalFn() Option[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if self.APIToken.IsNone() {
		return None[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
	}
	return Some(self.searchProposal)
}

func (self Connector) SquashMergeProposalFn() Option[func(int, gitdomain.CommitMessage) error] {
	if self.APIToken.IsNone() {
		return None[func(int, gitdomain.CommitMessage) error]()
	}
	return Some(self.squashMergeProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	if self.APIToken.IsNone() {
		return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error]()
	}
	return Some(self.updateProposalTarget)
}

func (self Connector) VerifyConnection() (string, error) {
	user, _, err := self.client.Users.CurrentUser()
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func (self Connector) VerifyReadProposalPermission() error {
	_, _, err := self.client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
	})
	return err
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
		TargetBranch: gitlab.Ptr(target.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parseMergeRequest(mergeRequests[0])
		self.log.Success(strconv.Itoa(proposal.Number))
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitLab}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(mergeRequests), branch, target)
	}
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := forgedomain.ReadProposalOverride()
	self.log.Ok()
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
		ForgeType: forgedomain.ForgeTypeGitLab,
	}), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parseMergeRequest(mergeRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitLab}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(mergeRequests), branch)
	}
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGitlabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	_, _, err := self.client.MergeRequests.AcceptMergeRequest(self.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.Ptr(message.String()),
		Squash:              gitlab.Ptr(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Ptr(false),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self Connector) updateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	data := proposalData.Data()
	self.log.Start(messages.ForgeGitlabUpdateMRViaAPI, data.Number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), data.Number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.Ptr(target.String()),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

// NewGitlabConfig provides GitLab configuration data if the current repo is hosted on GitLab,
// otherwise nil.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	gitlabData := Data{
		APIToken: args.APIToken,
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	client, err := gitlab.NewClient(gitlabData.APIToken.String(), gitlab.WithBaseURL(gitlabData.baseURL()))
	if err != nil {
		return Connector{}, err
	}
	connector := Connector{
		Data:   gitlabData,
		client: client,
		log:    args.Log,
	}
	return connector, nil
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.GitLabToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

func parseMergeRequest(mergeRequest *gitlab.BasicMergeRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: true,
		Number:       mergeRequest.IID,
		Source:       gitdomain.NewLocalBranchName(mergeRequest.SourceBranch),
		Target:       gitdomain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		Body:         NewOption(mergeRequest.Description),
		URL:          mergeRequest.WebURL,
	}
}
