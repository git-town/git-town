package gitlab

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/xanzy/go-gitlab"
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type Connector struct {
	client *gitlab.Client
	Data
	log print.Logger
}

func (self Connector) CanMakeAPICalls() bool {
	return self.Data.APIToken.IsSome() || len(hostingdomain.ReadProposalOverride()) > 0
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	return Some(func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
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
				Source:       branch,
				Target:       target,
				Title:        "title",
				URL:          proposalURLOverride,
			}), nil
		}
		opts := &gitlab.ListProjectMergeRequestsOptions{
			State:        gitlab.Ptr("opened"),
			SourceBranch: gitlab.Ptr(branch.String()),
			TargetBranch: gitlab.Ptr(target.String()),
		}
		mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
		if err != nil {
			self.log.Failed(err.Error())
			return None[hostingdomain.Proposal](), err
		}
		switch len(mergeRequests) {
		case 0:
			self.log.Success("none")
			return None[hostingdomain.Proposal](), nil
		case 1:
			proposal := parseMergeRequest(mergeRequests[0])
			self.log.Success(strconv.Itoa(proposal.Number))
			return Some(proposal), nil
		default:
			return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(mergeRequests), branch, target)
		}
	})
}

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	return Some(func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
		self.log.Start(messages.APIParentBranchLookupStart, branch.String())
		opts := &gitlab.ListProjectMergeRequestsOptions{
			State:        gitlab.Ptr("opened"),
			SourceBranch: gitlab.Ptr(branch.String()),
		}
		mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
		if err != nil {
			self.log.Failed(err.Error())
			return None[hostingdomain.Proposal](), err
		}
		switch len(mergeRequests) {
		case 0:
			self.log.Success("none")
			return None[hostingdomain.Proposal](), nil
		case 1:
			proposal := parseMergeRequest(mergeRequests[0])
			self.log.Success(proposal.Target.String())
			return Some(proposal), nil
		default:
			return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(mergeRequests), branch)
		}
	})
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error] {
	return Some(func(number int, message gitdomain.CommitMessage) error {
		if number <= 0 {
			return errors.New(messages.ProposalNoNumberGiven)
		}
		self.log.Start(messages.HostingGitlabMergingViaAPI, number)
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
	})
}

func (self Connector) UpdateProposalSourceFn() Option[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTarget(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	self.log.Start(messages.HostingGitlabUpdateMRViaAPI, number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), number, &gitlab.UpdateMergeRequestOptions{
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
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	clientOptFunc := gitlab.WithBaseURL(gitlabData.baseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{}) //exhaustruct:ignore
	client, err := gitlab.NewOAuthClient(gitlabData.APIToken.String(), httpClient, clientOptFunc)
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

func parseMergeRequest(mergeRequest *gitlab.MergeRequest) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		MergeWithAPI: true,
		Number:       mergeRequest.IID,
		Source:       gitdomain.NewLocalBranchName(mergeRequest.SourceBranch),
		Target:       gitdomain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		URL:          mergeRequest.WebURL,
	}
}
