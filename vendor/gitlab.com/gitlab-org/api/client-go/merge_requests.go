//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type (
	MergeRequestsServiceInterface interface {
		ListMergeRequests(opt *ListMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)
		ListProjectMergeRequests(pid any, opt *ListProjectMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)
		ListGroupMergeRequests(gid any, opt *ListGroupMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)
		GetMergeRequest(pid any, mergeRequest int64, opt *GetMergeRequestsOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		GetMergeRequestApprovals(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error)
		GetMergeRequestCommits(pid any, mergeRequest int64, opt *GetMergeRequestCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		GetMergeRequestChanges(pid any, mergeRequest int64, opt *GetMergeRequestChangesOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		ListMergeRequestDiffs(pid any, mergeRequest int64, opt *ListMergeRequestDiffsOptions, options ...RequestOptionFunc) ([]*MergeRequestDiff, *Response, error)
		ShowMergeRequestRawDiffs(pid any, mergeRequest int64, opt *ShowMergeRequestRawDiffsOptions, options ...RequestOptionFunc) ([]byte, *Response, error)
		GetMergeRequestParticipants(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*BasicUser, *Response, error)
		GetMergeRequestReviewers(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*MergeRequestReviewer, *Response, error)
		ListMergeRequestPipelines(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*PipelineInfo, *Response, error)
		CreateMergeRequestPipeline(pid any, mergeRequest int64, options ...RequestOptionFunc) (*PipelineInfo, *Response, error)
		GetIssuesClosedOnMerge(pid any, mergeRequest int64, opt *GetIssuesClosedOnMergeOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		ListRelatedIssues(pid any, mergeRequest int64, opt *ListRelatedIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		CreateMergeRequest(pid any, opt *CreateMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		UpdateMergeRequest(pid any, mergeRequest int64, opt *UpdateMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		DeleteMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Response, error)
		AcceptMergeRequest(pid any, mergeRequest int64, opt *AcceptMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		CancelMergeWhenPipelineSucceeds(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		RebaseMergeRequest(pid any, mergeRequest int64, opt *RebaseMergeRequestOptions, options ...RequestOptionFunc) (*Response, error)
		GetMergeRequestDiffVersions(pid any, mergeRequest int64, opt *GetMergeRequestDiffVersionsOptions, options ...RequestOptionFunc) ([]*MergeRequestDiffVersion, *Response, error)
		GetSingleMergeRequestDiffVersion(pid any, mergeRequest, version int64, opt *GetSingleMergeRequestDiffVersionOptions, options ...RequestOptionFunc) (*MergeRequestDiffVersion, *Response, error)
		SubscribeToMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		UnsubscribeFromMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		CreateTodo(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Todo, *Response, error)
		SetTimeEstimate(pid any, mergeRequest int64, opt *SetTimeEstimateOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error)
		ResetTimeEstimate(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error)
		AddSpentTime(pid any, mergeRequest int64, opt *AddSpentTimeOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error)
		ResetSpentTime(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error)
		GetTimeSpent(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error)
		CreateMergeRequestDependency(pid any, mergeRequest int64, opts CreateMergeRequestDependencyOptions, options ...RequestOptionFunc) (*MergeRequestDependency, *Response, error)
		DeleteMergeRequestDependency(pid any, mergeRequest int64, blockingMergeRequest int64, options ...RequestOptionFunc) (*Response, error)
		GetMergeRequestDependencies(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]MergeRequestDependency, *Response, error)
	}

	// MergeRequestsService handles communication with the merge requests related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/merge_requests/
	MergeRequestsService struct {
		client    *Client
		timeStats *timeStatsService
	}
)

var _ MergeRequestsServiceInterface = (*MergeRequestsService)(nil)

// BasicMergeRequest represents a basic GitLab merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/
type BasicMergeRequest struct {
	ID                          int64                  `json:"id"`
	IID                         int64                  `json:"iid"`
	TargetBranch                string                 `json:"target_branch"`
	SourceBranch                string                 `json:"source_branch"`
	ProjectID                   int64                  `json:"project_id"`
	Title                       string                 `json:"title"`
	State                       string                 `json:"state"`
	Imported                    bool                   `json:"imported"`
	ImportedFrom                string                 `json:"imported_from"`
	CreatedAt                   *time.Time             `json:"created_at"`
	UpdatedAt                   *time.Time             `json:"updated_at"`
	Upvotes                     int64                  `json:"upvotes"`
	Downvotes                   int64                  `json:"downvotes"`
	Author                      *BasicUser             `json:"author"`
	Assignee                    *BasicUser             `json:"assignee"`
	Assignees                   []*BasicUser           `json:"assignees"`
	Reviewers                   []*BasicUser           `json:"reviewers"`
	SourceProjectID             int64                  `json:"source_project_id"`
	TargetProjectID             int64                  `json:"target_project_id"`
	Labels                      Labels                 `json:"labels"`
	LabelDetails                []*LabelDetails        `json:"label_details"`
	Description                 string                 `json:"description"`
	Draft                       bool                   `json:"draft"`
	Milestone                   *Milestone             `json:"milestone"`
	MergeWhenPipelineSucceeds   bool                   `json:"merge_when_pipeline_succeeds"`
	DetailedMergeStatus         string                 `json:"detailed_merge_status"`
	MergeUser                   *BasicUser             `json:"merge_user"`
	MergedAt                    *time.Time             `json:"merged_at"`
	MergeAfter                  *time.Time             `json:"merge_after"`
	PreparedAt                  *time.Time             `json:"prepared_at"`
	ClosedBy                    *BasicUser             `json:"closed_by"`
	ClosedAt                    *time.Time             `json:"closed_at"`
	SHA                         string                 `json:"sha"`
	MergeCommitSHA              string                 `json:"merge_commit_sha"`
	SquashCommitSHA             string                 `json:"squash_commit_sha"`
	UserNotesCount              int64                  `json:"user_notes_count"`
	ShouldRemoveSourceBranch    bool                   `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch     bool                   `json:"force_remove_source_branch"`
	AllowCollaboration          bool                   `json:"allow_collaboration"`
	AllowMaintainerToPush       bool                   `json:"allow_maintainer_to_push"`
	WebURL                      string                 `json:"web_url"`
	References                  *IssueReferences       `json:"references"`
	DiscussionLocked            bool                   `json:"discussion_locked"`
	TimeStats                   *TimeStats             `json:"time_stats"`
	Squash                      bool                   `json:"squash"`
	SquashOnMerge               bool                   `json:"squash_on_merge"`
	TaskCompletionStatus        *TasksCompletionStatus `json:"task_completion_status"`
	HasConflicts                bool                   `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                   `json:"blocking_discussions_resolved"`

	// Deprecated: will be removed in v5 of the API, use MergeUser instead
	MergedBy *BasicUser `json:"merged_by"`
}

func (m BasicMergeRequest) String() string {
	return Stringify(m)
}

// MergeRequest represents a GitLab merge request.
//
// GitLab API docs: https://docs.gitlab.com/api/merge_requests/
type MergeRequest struct {
	BasicMergeRequest
	MergeError                  string               `json:"merge_error"`
	Subscribed                  bool                 `json:"subscribed"`
	ChangesCount                string               `json:"changes_count"`
	User                        MergeRequestUser     `json:"user"`
	LatestBuildStartedAt        *time.Time           `json:"latest_build_started_at"`
	LatestBuildFinishedAt       *time.Time           `json:"latest_build_finished_at"`
	FirstDeployedToProductionAt *time.Time           `json:"first_deployed_to_production_at"`
	Pipeline                    *PipelineInfo        `json:"pipeline"`
	HeadPipeline                *Pipeline            `json:"head_pipeline"`
	DiffRefs                    MergeRequestDiffRefs `json:"diff_refs"`
	RebaseInProgress            bool                 `json:"rebase_in_progress"`
	DivergedCommitsCount        int64                `json:"diverged_commits_count"`
	FirstContribution           bool                 `json:"first_contribution"`

	// Deprecated: use Draft instead
	WorkInProgress bool `json:"work_in_progress"`
}

func (m MergeRequest) String() string {
	return Stringify(m)
}

// MergeRequestUser represents a GitLab merge request user.
//
// GitLab API docs: https://docs.gitlab.com/api/merge_requests/
type MergeRequestUser struct {
	CanMerge bool `json:"can_merge"`
}

func (u MergeRequestUser) String() string {
	return Stringify(u)
}

// MergeRequestDiffRefs represents a GitLab merge request diff refs.
//
// GitLab API docs: https://docs.gitlab.com/api/merge_requests/
type MergeRequestDiffRefs struct {
	BaseSha  string `json:"base_sha"`
	HeadSha  string `json:"head_sha"`
	StartSha string `json:"start_sha"`
}

func (d MergeRequestDiffRefs) String() string {
	return Stringify(d)
}

func (m *MergeRequest) UnmarshalJSON(data []byte) error {
	type alias MergeRequest

	raw := make(map[string]any)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	labelDetails, ok := raw["labels"].([]any)
	if ok && len(labelDetails) > 0 {
		// We only want to change anything if we got label details.
		if _, ok := labelDetails[0].(map[string]any); !ok {
			return json.Unmarshal(data, (*alias)(m))
		}

		labels := make([]any, len(labelDetails))
		for i, details := range labelDetails {
			labels[i] = details.(map[string]any)["name"]
		}

		// Set the correct values
		raw["labels"] = labels
		raw["label_details"] = labelDetails

		data, err = json.Marshal(raw)
		if err != nil {
			return err
		}
	}

	return json.Unmarshal(data, (*alias)(m))
}

// MergeRequestDiff represents GitLab merge request diff.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-request-diffs
type MergeRequestDiff struct {
	OldPath       string `json:"old_path"`
	NewPath       string `json:"new_path"`
	AMode         string `json:"a_mode"`
	BMode         string `json:"b_mode"`
	Diff          string `json:"diff"`
	NewFile       bool   `json:"new_file"`
	RenamedFile   bool   `json:"renamed_file"`
	DeletedFile   bool   `json:"deleted_file"`
	GeneratedFile bool   `json:"generated_file"`
	Collapsed     bool   `json:"collapsed"`
	TooLarge      bool   `json:"too_large"`
}

// MergeRequestDiffVersion represents GitLab merge request version.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-merge-request-diff-versions
type MergeRequestDiffVersion struct {
	ID             int64      `json:"id"`
	HeadCommitSHA  string     `json:"head_commit_sha,omitempty"`
	BaseCommitSHA  string     `json:"base_commit_sha,omitempty"`
	StartCommitSHA string     `json:"start_commit_sha,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	MergeRequestID int64      `json:"merge_request_id,omitempty"`
	State          string     `json:"state,omitempty"`
	RealSize       string     `json:"real_size,omitempty"`
	Commits        []*Commit  `json:"commits,omitempty"`
	Diffs          []*Diff    `json:"diffs,omitempty"`
}

func (m MergeRequestDiffVersion) String() string {
	return Stringify(m)
}

// MergeRequestReviewer represents a reviewer entry returned by the reviewers API.
// Matches the JSON shape used in tests: {"user": {...}, "state": "...", "created_at": "..."}
// Placed here because it's used by MergeRequestsService.GetMergeRequestReviewers and tests/mock.
type MergeRequestReviewer struct {
	User      *BasicUser `json:"user"`
	State     string     `json:"state"`
	CreatedAt *time.Time `json:"created_at"`
}

// ListMergeRequestsOptions represents the available ListMergeRequests()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-requests
type ListMergeRequestsOptions struct {
	ListOptions
	Approved               *string           `url:"approved,omitempty" json:"approved,omitempty"`
	State                  *string           `url:"state,omitempty" json:"state,omitempty"`
	OrderBy                *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                   *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Milestone              *string           `url:"milestone,omitempty" json:"milestone,omitempty"`
	View                   *string           `url:"view,omitempty" json:"view,omitempty"`
	Labels                 *LabelOptions     `url:"labels,comma,omitempty" json:"labels,omitempty"`
	NotLabels              *LabelOptions     `url:"not[labels],comma,omitempty" json:"not[labels],omitempty"`
	WithLabelsDetails      *bool             `url:"with_labels_details,omitempty" json:"with_labels_details,omitempty"`
	WithMergeStatusRecheck *bool             `url:"with_merge_status_recheck,omitempty" json:"with_merge_status_recheck,omitempty"`
	CreatedAfter           *time.Time        `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore          *time.Time        `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter           *time.Time        `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore          *time.Time        `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Scope                  *string           `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID               *int64            `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername         *string           `url:"author_username,omitempty" json:"author_username,omitempty"`
	NotAuthorUsername      *string           `url:"not[author_username],omitempty" json:"not[author_username],omitempty"`
	AssigneeID             *AssigneeIDValue  `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	ApproverIDs            *ApproverIDsValue `url:"approver_ids,omitempty" json:"approver_ids,omitempty"`
	ApprovedByIDs          *ApproverIDsValue `url:"approved_by_ids,omitempty" json:"approved_by_ids,omitempty"`
	ReviewerID             *ReviewerIDValue  `url:"reviewer_id,omitempty" json:"reviewer_id,omitempty"`
	ReviewerUsername       *string           `url:"reviewer_username,omitempty" json:"reviewer_username,omitempty"`
	MyReactionEmoji        *string           `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	SourceBranch           *string           `url:"source_branch,omitempty" json:"source_branch,omitempty"`
	TargetBranch           *string           `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	Search                 *string           `url:"search,omitempty" json:"search,omitempty"`
	In                     *string           `url:"in,omitempty" json:"in,omitempty"`
	Draft                  *bool             `url:"draft,omitempty" json:"draft,omitempty"`
	WIP                    *string           `url:"wip,omitempty" json:"wip,omitempty"`
}

// ListMergeRequests gets all merge requests. The state parameter can be used
// to get only merge requests with a given state (opened, closed, or merged)
// or all of them (all). The pagination parameters page and per_page can be
// used to restrict the list of merge requests.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-requests
func (s *MergeRequestsService) ListMergeRequests(opt *ListMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	mrs, resp, err := do[[]*MergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("merge_requests"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}

	return toBasic(mrs), resp, nil
}

// ListProjectMergeRequestsOptions represents the available ListMergeRequests()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-project-merge-requests
type ListProjectMergeRequestsOptions struct {
	ListOptions
	IIDs                   *[]int64          `url:"iids[],omitempty" json:"iids,omitempty"`
	State                  *string           `url:"state,omitempty" json:"state,omitempty"`
	OrderBy                *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                   *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Milestone              *string           `url:"milestone,omitempty" json:"milestone,omitempty"`
	View                   *string           `url:"view,omitempty" json:"view,omitempty"`
	Environment            *string           `url:"environment,omitempty" json:"environment,omitempty"`
	Labels                 *LabelOptions     `url:"labels,comma,omitempty" json:"labels,omitempty"`
	NotLabels              *LabelOptions     `url:"not[labels],comma,omitempty" json:"not[labels],omitempty"`
	WithLabelsDetails      *bool             `url:"with_labels_details,omitempty" json:"with_labels_details,omitempty"`
	WithMergeStatusRecheck *bool             `url:"with_merge_status_recheck,omitempty" json:"with_merge_status_recheck,omitempty"`
	CreatedAfter           *time.Time        `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore          *time.Time        `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter           *time.Time        `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore          *time.Time        `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	DeployedBefore         *time.Time        `url:"deployed_before,omitempty" json:"deployed_before,omitempty"`
	DeployedAfter          *time.Time        `url:"deployed_after,omitempty" json:"deployed_after,omitempty"`
	Scope                  *string           `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID               *int64            `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername         *string           `url:"author_username,omitempty" json:"author_username,omitempty"`
	NotAuthorUsername      *string           `url:"not[author_username],omitempty" json:"not[author_username],omitempty"`
	AssigneeID             *AssigneeIDValue  `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	ApproverIDs            *ApproverIDsValue `url:"approver_ids,omitempty" json:"approver_ids,omitempty"`
	ApprovedByIDs          *ApproverIDsValue `url:"approved_by_ids,omitempty" json:"approved_by_ids,omitempty"`
	ReviewerID             *ReviewerIDValue  `url:"reviewer_id,omitempty" json:"reviewer_id,omitempty"`
	ReviewerUsername       *string           `url:"reviewer_username,omitempty" json:"reviewer_username,omitempty"`
	MyReactionEmoji        *string           `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	SourceBranch           *string           `url:"source_branch,omitempty" json:"source_branch,omitempty"`
	TargetBranch           *string           `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	Search                 *string           `url:"search,omitempty" json:"search,omitempty"`
	Draft                  *bool             `url:"draft,omitempty" json:"draft,omitempty"`
	WIP                    *string           `url:"wip,omitempty" json:"wip,omitempty"`
}

// ListProjectMergeRequests gets all merge requests for this project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-project-merge-requests
func (s *MergeRequestsService) ListProjectMergeRequests(pid any, opt *ListProjectMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	mrs, resp, err := do[[]*MergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}

	return toBasic(mrs), resp, nil
}

// ListGroupMergeRequestsOptions represents the available ListGroupMergeRequests()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-group-merge-requests
type ListGroupMergeRequestsOptions struct {
	ListOptions
	State                  *string           `url:"state,omitempty" json:"state,omitempty"`
	OrderBy                *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                   *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Milestone              *string           `url:"milestone,omitempty" json:"milestone,omitempty"`
	View                   *string           `url:"view,omitempty" json:"view,omitempty"`
	Labels                 *LabelOptions     `url:"labels,comma,omitempty" json:"labels,omitempty"`
	NotLabels              *LabelOptions     `url:"not[labels],comma,omitempty" json:"not[labels],omitempty"`
	WithLabelsDetails      *bool             `url:"with_labels_details,omitempty" json:"with_labels_details,omitempty"`
	WithMergeStatusRecheck *bool             `url:"with_merge_status_recheck,omitempty" json:"with_merge_status_recheck,omitempty"`
	CreatedAfter           *time.Time        `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore          *time.Time        `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter           *time.Time        `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore          *time.Time        `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Scope                  *string           `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID               *int64            `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername         *string           `url:"author_username,omitempty" json:"author_username,omitempty"`
	NotAuthorUsername      *string           `url:"not[author_username],omitempty" json:"not[author_username],omitempty"`
	AssigneeID             *AssigneeIDValue  `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	ApproverIDs            *ApproverIDsValue `url:"approver_ids,omitempty" json:"approver_ids,omitempty"`
	ApprovedByIDs          *ApproverIDsValue `url:"approved_by_ids,omitempty" json:"approved_by_ids,omitempty"`
	ReviewerID             *ReviewerIDValue  `url:"reviewer_id,omitempty" json:"reviewer_id,omitempty"`
	ReviewerUsername       *string           `url:"reviewer_username,omitempty" json:"reviewer_username,omitempty"`
	MyReactionEmoji        *string           `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	SourceBranch           *string           `url:"source_branch,omitempty" json:"source_branch,omitempty"`
	TargetBranch           *string           `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	Search                 *string           `url:"search,omitempty" json:"search,omitempty"`
	In                     *string           `url:"in,omitempty" json:"in,omitempty"`
	Draft                  *bool             `url:"draft,omitempty" json:"draft,omitempty"`
	WIP                    *string           `url:"wip,omitempty" json:"wip,omitempty"`
}

// ListGroupMergeRequests gets all merge requests for this group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-group-merge-requests
func (s *MergeRequestsService) ListGroupMergeRequests(gid any, opt *ListGroupMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	mrs, resp, err := do[[]*MergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/merge_requests", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}

	return toBasic(mrs), resp, nil
}

// GetMergeRequestsOptions represents the available GetMergeRequests()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-mr
type GetMergeRequestsOptions struct {
	RenderHTML                  *bool `url:"render_html,omitempty" json:"render_html,omitempty"`
	IncludeDivergedCommitsCount *bool `url:"include_diverged_commits_count,omitempty" json:"include_diverged_commits_count,omitempty"`
	IncludeRebaseInProgress     *bool `url:"include_rebase_in_progress,omitempty" json:"include_rebase_in_progress,omitempty"`
}

// GetMergeRequest shows information about a single merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-mr
func (s *MergeRequestsService) GetMergeRequest(pid any, mergeRequest int64, opt *GetMergeRequestsOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetMergeRequestApprovals gets information about a merge requests approvals
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#single-merge-request-approval
func (s *MergeRequestsService) GetMergeRequestApprovals(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error) {
	return do[*MergeRequestApprovals](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/approvals", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// GetMergeRequestCommitsOptions represents the available GetMergeRequestCommits()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-commits
type GetMergeRequestCommitsOptions struct {
	ListOptions
}

// GetMergeRequestCommits gets a list of merge request commits.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-commits
func (s *MergeRequestsService) GetMergeRequestCommits(pid any, mergeRequest int64, opt *GetMergeRequestCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/commits", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetMergeRequestChangesOptions represents the available GetMergeRequestChanges()
// options.
// Deprecated: This endpoint has been replaced by
// MergeRequestsService.ListMergeRequestDiffs()
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-changes
type GetMergeRequestChangesOptions struct {
	AccessRawDiffs *bool `url:"access_raw_diffs,omitempty" json:"access_raw_diffs,omitempty"`
	Unidiff        *bool `url:"unidiff,omitempty" json:"unidiff,omitempty"`
}

// GetMergeRequestChanges shows information about the merge request including
// its files and changes.
//
// Deprecated: This endpoint has been replaced by
// MergeRequestsService.ListMergeRequestDiffs()
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-changes
func (s *MergeRequestsService) GetMergeRequestChanges(pid any, mergeRequest int64, opt *GetMergeRequestChangesOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/changes", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListMergeRequestDiffsOptions represents the available ListMergeRequestDiffs()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-request-diffs
type ListMergeRequestDiffsOptions struct {
	ListOptions
	Unidiff *bool `url:"unidiff,omitempty" json:"unidiff,omitempty"`
}

// ListMergeRequestDiffs List diffs of the files changed in a merge request
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-request-diffs
func (s *MergeRequestsService) ListMergeRequestDiffs(pid any, mergeRequest int64, opt *ListMergeRequestDiffsOptions, options ...RequestOptionFunc) ([]*MergeRequestDiff, *Response, error) {
	return do[[]*MergeRequestDiff](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/diffs", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ShowMergeRequestRawDiffsOptions represents the available ShowMergeRequestRawDiffs()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#show-merge-request-raw-diffs
type ShowMergeRequestRawDiffsOptions struct{}

// ShowMergeRequestRawDiffs Show raw diffs of the files changed in a merge request
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#show-merge-request-raw-diffs
func (s *MergeRequestsService) ShowMergeRequestRawDiffs(pid any, mergeRequest int64, opt *ShowMergeRequestRawDiffsOptions, options ...RequestOptionFunc) ([]byte, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/merge_requests/%d/raw_diffs", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)

	return b.Bytes(), resp, err
}

// GetMergeRequestParticipants gets a list of merge request participants.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-participants
func (s *MergeRequestsService) GetMergeRequestParticipants(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*BasicUser, *Response, error) {
	return do[[]*BasicUser](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/participants", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// GetMergeRequestReviewers gets a list of merge request reviewers.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-reviewers
func (s *MergeRequestsService) GetMergeRequestReviewers(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*MergeRequestReviewer, *Response, error) {
	return do[[]*MergeRequestReviewer](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/reviewers", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// ListMergeRequestPipelines gets all pipelines for the provided merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-merge-request-pipelines
func (s *MergeRequestsService) ListMergeRequestPipelines(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*PipelineInfo, *Response, error) {
	return do[[]*PipelineInfo](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/pipelines", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// CreateMergeRequestPipeline creates a new pipeline for a merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-merge-request-pipeline
func (s *MergeRequestsService) CreateMergeRequestPipeline(pid any, mergeRequest int64, options ...RequestOptionFunc) (*PipelineInfo, *Response, error) {
	return do[*PipelineInfo](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/pipelines", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// GetIssuesClosedOnMergeOptions represents the available GetIssuesClosedOnMerge()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-issues-that-close-on-merge
type GetIssuesClosedOnMergeOptions struct {
	ListOptions
}

// GetIssuesClosedOnMerge gets all the issues that would be closed by merging the
// provided merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-issues-that-close-on-merge
func (s *MergeRequestsService) GetIssuesClosedOnMerge(pid any, mergeRequest int64, opt *GetIssuesClosedOnMergeOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	return do[[]*Issue](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/closes_issues", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListRelatedIssuesOptions represents the available ListRelatedIssues() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-issues-related-to-the-merge-request
type ListRelatedIssuesOptions struct {
	ListOptions
}

// ListRelatedIssues gets all the issues related to provided merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#list-issues-related-to-the-merge-request
func (s *MergeRequestsService) ListRelatedIssues(pid any, mergeRequest int64, opt *ListRelatedIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	return do[[]*Issue](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/related_issues", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateMergeRequestOptions represents the available CreateMergeRequest()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-mr
type CreateMergeRequestOptions struct {
	Title              *string       `url:"title,omitempty" json:"title,omitempty"`
	Description        *string       `url:"description,omitempty" json:"description,omitempty"`
	SourceBranch       *string       `url:"source_branch,omitempty" json:"source_branch,omitempty"`
	TargetBranch       *string       `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	Labels             *LabelOptions `url:"labels,comma,omitempty" json:"labels,omitempty"`
	AssigneeID         *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	AssigneeIDs        *[]int64      `url:"assignee_ids,omitempty" json:"assignee_ids,omitempty"`
	ReviewerIDs        *[]int64      `url:"reviewer_ids,omitempty" json:"reviewer_ids,omitempty"`
	TargetProjectID    *int64        `url:"target_project_id,omitempty" json:"target_project_id,omitempty"`
	MilestoneID        *int64        `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	RemoveSourceBranch *bool         `url:"remove_source_branch,omitempty" json:"remove_source_branch,omitempty"`
	Squash             *bool         `url:"squash,omitempty" json:"squash,omitempty"`
	AllowCollaboration *bool         `url:"allow_collaboration,omitempty" json:"allow_collaboration,omitempty"`

	// Deprecated: will be removed in v5 of the API, use the Merge Request Approvals API instead
	ApprovalsBeforeMerge *int64 `url:"approvals_before_merge,omitempty" json:"approvals_before_merge,omitempty"`
}

// CreateMergeRequest creates a new merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-mr
func (s *MergeRequestsService) CreateMergeRequest(pid any, opt *CreateMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateMergeRequestOptions represents the available UpdateMergeRequest()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#update-mr
type UpdateMergeRequestOptions struct {
	Title              *string       `url:"title,omitempty" json:"title,omitempty"`
	Description        *string       `url:"description,omitempty" json:"description,omitempty"`
	TargetBranch       *string       `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	AssigneeID         *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	AssigneeIDs        *[]int64      `url:"assignee_ids,omitempty" json:"assignee_ids,omitempty"`
	ReviewerIDs        *[]int64      `url:"reviewer_ids,omitempty" json:"reviewer_ids,omitempty"`
	Labels             *LabelOptions `url:"labels,comma,omitempty" json:"labels,omitempty"`
	AddLabels          *LabelOptions `url:"add_labels,comma,omitempty" json:"add_labels,omitempty"`
	RemoveLabels       *LabelOptions `url:"remove_labels,comma,omitempty" json:"remove_labels,omitempty"`
	MilestoneID        *int64        `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	StateEvent         *string       `url:"state_event,omitempty" json:"state_event,omitempty"`
	RemoveSourceBranch *bool         `url:"remove_source_branch,omitempty" json:"remove_source_branch,omitempty"`
	Squash             *bool         `url:"squash,omitempty" json:"squash,omitempty"`
	DiscussionLocked   *bool         `url:"discussion_locked,omitempty" json:"discussion_locked,omitempty"`
	AllowCollaboration *bool         `url:"allow_collaboration,omitempty" json:"allow_collaboration,omitempty"`
}

// UpdateMergeRequest updates an existing project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#update-mr
func (s *MergeRequestsService) UpdateMergeRequest(pid any, mergeRequest int64, opt *UpdateMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteMergeRequest deletes a merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#delete-a-merge-request
func (s *MergeRequestsService) DeleteMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/merge_requests/%d", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// AcceptMergeRequestOptions represents the available AcceptMergeRequest()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#merge-a-merge-request
type AcceptMergeRequestOptions struct {
	AutoMerge                *bool   `url:"auto_merge,omitempty" json:"auto_merge,omitempty"`
	MergeCommitMessage       *string `url:"merge_commit_message,omitempty" json:"merge_commit_message,omitempty"`
	SquashCommitMessage      *string `url:"squash_commit_message,omitempty" json:"squash_commit_message,omitempty"`
	Squash                   *bool   `url:"squash,omitempty" json:"squash,omitempty"`
	ShouldRemoveSourceBranch *bool   `url:"should_remove_source_branch,omitempty" json:"should_remove_source_branch,omitempty"`
	SHA                      *string `url:"sha,omitempty" json:"sha,omitempty"`

	// Deprecated: use AutoMerge instead
	MergeWhenPipelineSucceeds *bool `url:"merge_when_pipeline_succeeds,omitempty" json:"merge_when_pipeline_succeeds,omitempty"`
}

// AcceptMergeRequest merges changes submitted with MR using this API. If merge
// success you get 200 OK. If it has some conflicts and can not be merged - you
// get 405 and error message 'Branch cannot be merged'. If merge request is
// already merged or closed - you get 405 and error message 'Method Not Allowed'
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#merge-a-merge-request
func (s *MergeRequestsService) AcceptMergeRequest(pid any, mergeRequest int64, opt *AcceptMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/merge", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CancelMergeWhenPipelineSucceeds cancels a merge when pipeline succeeds. If
// you don't have permissions to accept this merge request - you'll get a 401.
// If the merge request is already merged or closed - you get 405 and error
// message 'Method Not Allowed'. In case the merge request is not set to be
// merged when the pipeline succeeds, you'll also get a 406 error.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#cancel-merge-when-pipeline-succeeds
func (s *MergeRequestsService) CancelMergeWhenPipelineSucceeds(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/cancel_merge_when_pipeline_succeeds", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// RebaseMergeRequestOptions represents the available RebaseMergeRequest()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#rebase-a-merge-request
type RebaseMergeRequestOptions struct {
	SkipCI *bool `url:"skip_ci,omitempty" json:"skip_ci,omitempty"`
}

// RebaseMergeRequest automatically rebases the source_branch of the merge
// request against its target_branch. If you don’t have permissions to push
// to the merge request’s source branch, you’ll get a 403 Forbidden response.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#rebase-a-merge-request
func (s *MergeRequestsService) RebaseMergeRequest(pid any, mergeRequest int64, opt *RebaseMergeRequestOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/rebase", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// GetMergeRequestDiffVersionsOptions represents the available
// GetMergeRequestDiffVersions() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-merge-request-diff-versions
type GetMergeRequestDiffVersionsOptions struct {
	ListOptions
}

// GetMergeRequestDiffVersions get a list of merge request diff versions.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-merge-request-diff-versions
func (s *MergeRequestsService) GetMergeRequestDiffVersions(pid any, mergeRequest int64, opt *GetMergeRequestDiffVersionsOptions, options ...RequestOptionFunc) ([]*MergeRequestDiffVersion, *Response, error) {
	return do[[]*MergeRequestDiffVersion](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/versions", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetSingleMergeRequestDiffVersionOptions represents the available
// GetSingleMergeRequestDiffVersion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-a-single-merge-request-diff-version
type GetSingleMergeRequestDiffVersionOptions struct {
	Unidiff *bool `url:"unidiff,omitempty" json:"unidiff,omitempty"`
}

// GetSingleMergeRequestDiffVersion get a single MR diff version
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-a-single-merge-request-diff-version
func (s *MergeRequestsService) GetSingleMergeRequestDiffVersion(pid any, mergeRequest, version int64, opt *GetSingleMergeRequestDiffVersionOptions, options ...RequestOptionFunc) (*MergeRequestDiffVersion, *Response, error) {
	return do[*MergeRequestDiffVersion](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/versions/%d", ProjectID{pid}, mergeRequest, version),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// SubscribeToMergeRequest subscribes the authenticated user to the given merge
// request to receive notifications. If the user is already subscribed to the
// merge request, the status code 304 is returned.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#subscribe-to-a-merge-request
func (s *MergeRequestsService) SubscribeToMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/subscribe", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// UnsubscribeFromMergeRequest unsubscribes the authenticated user from the
// given merge request to not receive notifications from that merge request.
// If the user is not subscribed to the merge request, status code 304 is
// returned.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#unsubscribe-from-a-merge-request
func (s *MergeRequestsService) UnsubscribeFromMergeRequest(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	return do[*MergeRequest](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/unsubscribe", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// CreateTodo manually creates a todo for the current user on a merge request.
// If there already exists a todo for the user on that merge request,
// status code 304 is returned.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-a-to-do-item
func (s *MergeRequestsService) CreateTodo(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Todo, *Response, error) {
	return do[*Todo](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/todo", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// SetTimeEstimate sets the time estimate for a single project merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#set-a-time-estimate-for-a-merge-request
func (s *MergeRequestsService) SetTimeEstimate(pid any, mergeRequest int64, opt *SetTimeEstimateOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	return s.timeStats.setTimeEstimate(pid, "merge_requests", mergeRequest, opt, options...)
}

// ResetTimeEstimate resets the time estimate for a single project merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#reset-the-time-estimate-for-a-merge-request
func (s *MergeRequestsService) ResetTimeEstimate(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	return s.timeStats.resetTimeEstimate(pid, "merge_requests", mergeRequest, options...)
}

// AddSpentTime adds spent time for a single project merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#add-spent-time-for-a-merge-request
func (s *MergeRequestsService) AddSpentTime(pid any, mergeRequest int64, opt *AddSpentTimeOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	return s.timeStats.addSpentTime(pid, "merge_requests", mergeRequest, opt, options...)
}

// ResetSpentTime resets the spent time for a single project merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#reset-spent-time-for-a-merge-request
func (s *MergeRequestsService) ResetSpentTime(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	return s.timeStats.resetSpentTime(pid, "merge_requests", mergeRequest, options...)
}

// GetTimeSpent gets the spent time for a single project merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-time-tracking-stats
func (s *MergeRequestsService) GetTimeSpent(pid any, mergeRequest int64, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	return s.timeStats.getTimeSpent(pid, "merge_requests", mergeRequest, options...)
}

// MergeRequestDependency represents a GitLab merge request dependency.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-a-merge-request-dependency
type MergeRequestDependency struct {
	ID                   int64                `json:"id"`
	BlockingMergeRequest BlockingMergeRequest `json:"blocking_merge_request"`
	ProjectID            int64                `json:"project_id"`
}

// BlockingMergeRequest represents a GitLab merge request dependency.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-a-merge-request-dependency
type BlockingMergeRequest struct {
	ID                          int64                  `json:"id"`
	Iid                         int64                  `json:"iid"`
	TargetBranch                string                 `json:"target_branch"`
	SourceBranch                string                 `json:"source_branch"`
	ProjectID                   int64                  `json:"project_id"`
	Title                       string                 `json:"title"`
	State                       string                 `json:"state"`
	CreatedAt                   time.Time              `json:"created_at"`
	UpdatedAt                   time.Time              `json:"updated_at"`
	Upvotes                     int64                  `json:"upvotes"`
	Downvotes                   int64                  `json:"downvotes"`
	Author                      *BasicUser             `json:"author"`
	Assignee                    *BasicUser             `json:"assignee"`
	Assignees                   []*BasicUser           `json:"assignees"`
	Reviewers                   []*BasicUser           `json:"reviewers"`
	SourceProjectID             int64                  `json:"source_project_id"`
	TargetProjectID             int64                  `json:"target_project_id"`
	Labels                      *LabelOptions          `json:"labels"`
	Description                 string                 `json:"description"`
	Draft                       bool                   `json:"draft"`
	Milestone                   *string                `json:"milestone"`
	AutoMerge                   bool                   `json:"auto_merge"`
	DetailedMergeStatus         string                 `json:"detailed_merge_status"`
	MergedAt                    *time.Time             `json:"merged_at"`
	ClosedBy                    *BasicUser             `json:"closed_by"`
	ClosedAt                    *time.Time             `json:"closed_at"`
	Sha                         string                 `json:"sha"`
	MergeCommitSha              string                 `json:"merge_commit_sha"`
	SquashCommitSha             string                 `json:"squash_commit_sha"`
	UserNotesCount              int64                  `json:"user_notes_count"`
	ShouldRemoveSourceBranch    *bool                  `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch     bool                   `json:"force_remove_source_branch"`
	WebURL                      string                 `json:"web_url"`
	References                  *IssueReferences       `json:"references"`
	DiscussionLocked            *bool                  `json:"discussion_locked"`
	TimeStats                   *TimeStats             `json:"time_stats"`
	Squash                      bool                   `json:"squash"`
	TaskCompletionStatus        *TasksCompletionStatus `json:"task_completion_status"`
	HasConflicts                bool                   `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                   `json:"blocking_discussions_resolved"`
	MergeUser                   *BasicUser             `json:"merge_user"`
	MergeAfter                  time.Time              `json:"merge_after"`
	Imported                    bool                   `json:"imported"`
	ImportedFrom                string                 `json:"imported_from"`
	PreparedAt                  *time.Time             `json:"prepared_at"`
	SquashOnMerge               bool                   `json:"squash_on_merge"`

	// Deprecated: use Draft instead
	WorkInProgress bool `json:"work_in_progress"`
	// Deprecated: will be removed in v5 of the API, use AutoMerge instead
	MergeWhenPipelineSucceeds bool `json:"merge_when_pipeline_succeeds"`
	// Deprecated: will be removed in v5 of the API, use MergeUser instead
	MergedBy *BasicUser `json:"merged_by"`
	// Deprecated: will be removed in v5 of the API, use the Merge Request Approvals API instead
	ApprovalsBeforeMerge *int64 `json:"approvals_before_merge"`
	// Deprecated: will be removed in v5 of the API, use References instead
	Reference string `json:"reference"`
	// Deprecated: in 15.6, use DetailedMergeStatus instead
	MergeStatus string `json:"merge_status"`
}

func (m MergeRequestDependency) String() string {
	return Stringify(m)
}

// CreateMergeRequestDependencyOptions represents the available CreateMergeRequestDependency()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-a-merge-request-dependency
type CreateMergeRequestDependencyOptions struct {
	BlockingMergeRequestID *int64 `url:"blocking_merge_request_id,omitempty" json:"blocking_merge_request_id,omitempty"`
}

// CreateMergeRequestDependency creates a new merge request dependency for a given
// merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#create-a-merge-request-dependency
func (s *MergeRequestsService) CreateMergeRequestDependency(pid any, mergeRequest int64, opts CreateMergeRequestDependencyOptions, options ...RequestOptionFunc) (*MergeRequestDependency, *Response, error) {
	return do[*MergeRequestDependency](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/blocks", ProjectID{pid}, mergeRequest),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// DeleteMergeRequestDependency deletes a merge request dependency for a given
// merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#delete-a-merge-request-dependency
func (s *MergeRequestsService) DeleteMergeRequestDependency(pid any, mergeRequest int64, blockingMergeRequest int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/merge_requests/%d/blocks/%d", ProjectID{pid}, mergeRequest, blockingMergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// GetMergeRequestDependencies gets a list of merge request dependencies.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_requests/#get-merge-request-dependencies
func (s *MergeRequestsService) GetMergeRequestDependencies(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]MergeRequestDependency, *Response, error) {
	return do[[]MergeRequestDependency](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_requests/%d/blocks", ProjectID{pid}, mergeRequest),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

func toBasic(mrs []*MergeRequest) []*BasicMergeRequest {
	ret := make([]*BasicMergeRequest, len(mrs))

	for i, mr := range mrs {
		ret[i] = &mr.BasicMergeRequest
	}

	return ret
}
