package bitbucketdatacenter

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to enforce conformance to the ProposalFinder interface
var _ forgedomain.ProposalFinder = bbdcAPIConnector

func (self AuthConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	proposalURLOverride := forgedomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		return self.findProposalViaOverride(branch, target)
	}
	return self.findProposalViaAPI(branch, target)
}

func (self AuthConnector) apiBaseURL() string {
	return fmt.Sprintf(
		"https://%s/rest/api/latest/projects/%s/repos/%s/pull-requests",
		self.Hostname,
		self.Organization,
		self.Repository,
	)
}

func (self AuthConnector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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

func (self AuthConnector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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

type PullRequestResponse struct {
	IsLastPage    bool          `json:"isLastPage"`
	Limit         int           `json:"limit"`
	NextPageStart int           `json:"nextPageStart"`
	Size          int           `json:"size"`
	Start         int           `json:"start"`
	Values        []PullRequest `json:"values"`
}

type Participant struct {
	Approved           bool   `json:"approved"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Status             string `json:"status"`
	User               User   `json:"user"`
}

type User struct {
	Active       bool   `json:"active"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type PullRequest struct {
	Closed       bool          `json:"closed"`
	ClosedDate   int64         `json:"closedDate"`
	CreatedDate  int64         `json:"createdDate"`
	Description  string        `json:"description"`
	Draft        bool          `json:"draft"`
	FromRef      Ref           `json:"fromRef"`
	ID           int           `json:"id"`
	Locked       bool          `json:"locked"`
	Open         bool          `json:"open"`
	Participants []Participant `json:"participants"`
	Reviewers    []Participant `json:"reviewers"`
	State        string        `json:"state"`
	Title        string        `json:"title"`
	ToRef        Ref           `json:"toRef"`
	UpdatedDate  int64         `json:"updatedDate"`
	Version      int           `json:"version"`
}

type Project struct {
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Public      bool   `json:"public"`
	Scope       string `json:"scope"`
	Type        string `json:"type"`
}

type Ref struct {
	DisplayID    string `json:"displayId"`
	ID           string `json:"id"`
	LatestCommit string `json:"latestCommit"`
	Repository   struct {
		Repository
		Origin Repository `json:"origin"`
	} `json:"repository"`
	Type string `json:"type"`
}

type Repository struct {
	Archived      bool     `json:"archived"`
	DefaultBranch string   `json:"defaultBranch"`
	Description   string   `json:"description"`
	Forkable      bool     `json:"forkable"`
	HierarchyID   string   `json:"hierarchyId"`
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Partition     int      `json:"partition"`
	Project       Project  `json:"project"`
	Public        bool     `json:"public"`
	RelatedLinks  struct{} `json:"relatedLinks"`
	ScmID         string   `json:"scmId"`
	Scope         string   `json:"scope"`
	Slug          string   `json:"slug"`
	State         string   `json:"state"`
	StatusMessage string   `json:"statusMessage"`
}
