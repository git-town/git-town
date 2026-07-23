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
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// StateID identifies the state of an issue or merge request.
//
// There are no GitLab API docs on the subject, but the mappings can be found in
// GitLab's codebase:
// https://gitlab.com/gitlab-org/gitlab-foss/-/blob/ba5be4989e/app/models/concerns/issuable.rb#L39-42
type StateID int64

const (
	StateIDNone   StateID = 0
	StateIDOpen   StateID = 1
	StateIDClosed StateID = 2
	StateIDMerged StateID = 3
	StateIDLocked StateID = 4
)

// BuildEvent represents a build event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#job-events
type BuildEvent struct {
	ObjectKind        string           `json:"object_kind"`
	Ref               string           `json:"ref"`
	Tag               bool             `json:"tag"`
	BeforeSHA         string           `json:"before_sha"`
	SHA               string           `json:"sha"`
	BuildID           int64            `json:"build_id"`
	BuildName         string           `json:"build_name"`
	BuildStage        string           `json:"build_stage"`
	BuildStatus       string           `json:"build_status"`
	BuildCreatedAt    string           `json:"build_created_at"`
	BuildStartedAt    string           `json:"build_started_at"`
	BuildFinishedAt   string           `json:"build_finished_at"`
	BuildDuration     float64          `json:"build_duration"`
	BuildAllowFailure bool             `json:"build_allow_failure"`
	ProjectID         int64            `json:"project_id"`
	ProjectName       string           `json:"project_name"`
	User              *EventUser       `json:"user"`
	Commit            BuildEventCommit `json:"commit"`
	Repository        *Repository      `json:"repository"`
}

type BuildEventCommit struct {
	ID          int64  `json:"id"`
	SHA         string `json:"sha"`
	Message     string `json:"message"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	Status      string `json:"status"`
	Duration    int64  `json:"duration"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
}

// CommitCommentEvent represents a comment on a commit event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#comment-on-a-commit
type CommitCommentEvent struct {
	ObjectKind       string                             `json:"object_kind"`
	EventType        string                             `json:"event_type"`
	User             *User                              `json:"user"`
	ProjectID        int64                              `json:"project_id"`
	Project          CommitCommentEventProject          `json:"project"`
	Repository       *Repository                        `json:"repository"`
	ObjectAttributes CommitCommentEventObjectAttributes `json:"object_attributes"`
	Commit           *CommitCommentEventCommit          `json:"commit"`
}

type CommitCommentEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type CommitCommentEventObjectAttributes struct {
	ID           int64              `json:"id"`
	Note         string             `json:"note"`
	NoteableType string             `json:"noteable_type"`
	AuthorID     int64              `json:"author_id"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
	ProjectID    int64              `json:"project_id"`
	Attachment   string             `json:"attachment"`
	LineCode     string             `json:"line_code"`
	CommitID     string             `json:"commit_id"`
	NoteableID   int64              `json:"noteable_id"`
	System       bool               `json:"system"`
	StDiff       *Diff              `json:"st_diff"`
	Description  string             `json:"description"`
	Action       CommentEventAction `json:"action"`
	URL          string             `json:"url"`
}

type CommitCommentEventCommit struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
}

type EventCommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// DeploymentEvent represents a deployment event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#deployment-events
type DeploymentEvent struct {
	ObjectKind             string                 `json:"object_kind"`
	Status                 string                 `json:"status"`
	StatusChangedAt        string                 `json:"status_changed_at"`
	DeploymentID           int64                  `json:"deployment_id"`
	DeployableID           int64                  `json:"deployable_id"`
	DeployableURL          string                 `json:"deployable_url"`
	Environment            string                 `json:"environment"`
	EnvironmentSlug        string                 `json:"environment_slug"`
	EnvironmentExternalURL string                 `json:"environment_external_url"`
	Project                DeploymentEventProject `json:"project"`
	Ref                    string                 `json:"ref"`
	ShortSHA               string                 `json:"short_sha"`
	User                   *EventUser             `json:"user"`
	UserURL                string                 `json:"user_url"`
	CommitURL              string                 `json:"commit_url"`
	CommitTitle            string                 `json:"commit_title"`
}

type DeploymentEventProject struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	AvatarURL         *string `json:"avatar_url"`
	GitSSHURL         string  `json:"git_ssh_url"`
	GitHTTPURL        string  `json:"git_http_url"`
	Namespace         string  `json:"namespace"`
	PathWithNamespace string  `json:"path_with_namespace"`
	DefaultBranch     string  `json:"default_branch"`
	Homepage          string  `json:"homepage"`
	URL               string  `json:"url"`
	SSHURL            string  `json:"ssh_url"`
	HTTPURL           string  `json:"http_url"`
	WebURL            string  `json:"web_url"`
	VisibilityLevel   int64   `json:"visibility_level"`
	CIConfigPath      string  `json:"ci_config_path"`
}

// FeatureFlagEvent represents a feature flag event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#feature-flag-events
type FeatureFlagEvent struct {
	ObjectKind       string                           `json:"object_kind"`
	Project          FeatureFlagEventProject          `json:"project"`
	User             *EventUser                       `json:"user"`
	UserURL          string                           `json:"user_url"`
	ObjectAttributes FeatureFlagEventObjectAttributes `json:"object_attributes"`
}

type FeatureFlagEventProject struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	AvatarURL         *string `json:"avatar_url"`
	GitSSHURL         string  `json:"git_ssh_url"`
	GitHTTPURL        string  `json:"git_http_url"`
	Namespace         string  `json:"namespace"`
	PathWithNamespace string  `json:"path_with_namespace"`
	DefaultBranch     string  `json:"default_branch"`
	Homepage          string  `json:"homepage"`
	URL               string  `json:"url"`
	SSHURL            string  `json:"ssh_url"`
	HTTPURL           string  `json:"http_url"`
	WebURL            string  `json:"web_url"`
	VisibilityLevel   int64   `json:"visibility_level"`
	CIConfigPath      string  `json:"ci_config_path"`
}

type FeatureFlagEventObjectAttributes struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// GroupResourceAccessTokenEvent represents a resource access token event for a
// group.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#project-and-group-access-token-events
type GroupResourceAccessTokenEvent struct {
	EventName        string                                        `json:"event_name"`
	ObjectKind       string                                        `json:"object_kind"`
	Group            GroupResourceAccessTokenEventGroup            `json:"group"`
	ObjectAttributes GroupResourceAccessTokenEventObjectAttributes `json:"object_attributes"`
}

// GroupResourceAccessTokenEventGroup represents a group in a resource access
// token event.
type GroupResourceAccessTokenEventGroup struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	GroupPath string `json:"group_path"`
	FullPath  string `json:"full_path"`
}

type GroupResourceAccessTokenEventObjectAttributes struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	ExpiresAt *ISOTime `json:"expires_at"`
}

// IssueCommentEvent represents a comment on an issue event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#comment-on-an-issue
type IssueCommentEvent struct {
	ObjectKind       string                            `json:"object_kind"`
	EventType        string                            `json:"event_type"`
	User             *User                             `json:"user"`
	ProjectID        int64                             `json:"project_id"`
	Project          IssueCommentEventProject          `json:"project"`
	Repository       *Repository                       `json:"repository"`
	ObjectAttributes IssueCommentEventObjectAttributes `json:"object_attributes"`
	Issue            IssueCommentEventIssue            `json:"issue"`
}

type IssueCommentEventProject struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type IssueCommentEventObjectAttributes struct {
	ID           int64              `json:"id"`
	Note         string             `json:"note"`
	NoteableType string             `json:"noteable_type"`
	AuthorID     int64              `json:"author_id"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
	ProjectID    int64              `json:"project_id"`
	Attachment   string             `json:"attachment"`
	LineCode     string             `json:"line_code"`
	CommitID     string             `json:"commit_id"`
	DiscussionID string             `json:"discussion_id"`
	NoteableID   int64              `json:"noteable_id"`
	System       bool               `json:"system"`
	StDiff       []*Diff            `json:"st_diff"`
	Description  string             `json:"description"`
	Action       CommentEventAction `json:"action"`
	URL          string             `json:"url"`
}

type IssueCommentEventIssue struct {
	ID                  int64         `json:"id"`
	IID                 int64         `json:"iid"`
	ProjectID           int64         `json:"project_id"`
	MilestoneID         int64         `json:"milestone_id"`
	AuthorID            int64         `json:"author_id"`
	Position            int64         `json:"position"`
	BranchName          string        `json:"branch_name"`
	Description         string        `json:"description"`
	State               string        `json:"state"`
	Title               string        `json:"title"`
	Labels              []*EventLabel `json:"labels"`
	LastEditedAt        string        `json:"last_edit_at"`
	LastEditedByID      int64         `json:"last_edited_by_id"`
	UpdatedAt           string        `json:"updated_at"`
	UpdatedByID         int64         `json:"updated_by_id"`
	CreatedAt           string        `json:"created_at"`
	ClosedAt            string        `json:"closed_at"`
	DueDate             *ISOTime      `json:"due_date"`
	URL                 string        `json:"url"`
	TimeEstimate        int64         `json:"time_estimate"`
	Confidential        bool          `json:"confidential"`
	TotalTimeSpent      int64         `json:"total_time_spent"`
	HumanTotalTimeSpent string        `json:"human_total_time_spent"`
	HumanTimeEstimate   string        `json:"human_time_estimate"`
	AssigneeIDs         []int64       `json:"assignee_ids"`
	AssigneeID          int64         `json:"assignee_id"`
}

// IssueEvent represents a issue event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#work-item-events
type IssueEvent struct {
	ObjectKind       string                     `json:"object_kind"`
	EventType        string                     `json:"event_type"`
	User             *EventUser                 `json:"user"`
	Project          IssueEventProject          `json:"project"`
	Repository       *Repository                `json:"repository"`
	ObjectAttributes IssueEventObjectAttributes `json:"object_attributes"`
	Assignee         *EventUser                 `json:"assignee"`
	Assignees        *[]EventUser               `json:"assignees"`
	Labels           []*EventLabel              `json:"labels"`
	Changes          IssueEventChanges          `json:"changes"`
}

type IssueEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type IssueEventObjectAttributes struct {
	ID                  int64                                      `json:"id"`
	Title               string                                     `json:"title"`
	AssigneeIDs         []int64                                    `json:"assignee_ids"`
	AssigneeID          int64                                      `json:"assignee_id"`
	AuthorID            int64                                      `json:"author_id"`
	ProjectID           int64                                      `json:"project_id"`
	CreatedAt           string                                     `json:"created_at"` // Should be *time.Time (see Gitlab issue #21468)
	UpdatedAt           string                                     `json:"updated_at"` // Should be *time.Time (see Gitlab issue #21468)
	UpdatedByID         int64                                      `json:"updated_by_id"`
	LastEditedAt        string                                     `json:"last_edited_at"`
	LastEditedByID      int64                                      `json:"last_edited_by_id"`
	RelativePosition    int64                                      `json:"relative_position"`
	BranchName          string                                     `json:"branch_name"`
	Description         string                                     `json:"description"`
	MilestoneID         int64                                      `json:"milestone_id"`
	StateID             StateID                                    `json:"state_id"`
	Confidential        bool                                       `json:"confidential"`
	DiscussionLocked    bool                                       `json:"discussion_locked"`
	DueDate             *ISOTime                                   `json:"due_date"`
	MovedToID           int64                                      `json:"moved_to_id"`
	DuplicatedToID      int64                                      `json:"duplicated_to_id"`
	TimeEstimate        int64                                      `json:"time_estimate"`
	TotalTimeSpent      int64                                      `json:"total_time_spent"`
	TimeChange          int64                                      `json:"time_change"`
	HumanTotalTimeSpent string                                     `json:"human_total_time_spent"`
	HumanTimeEstimate   string                                     `json:"human_time_estimate"`
	HumanTimeChange     string                                     `json:"human_time_change"`
	Weight              int64                                      `json:"weight"`
	IID                 int64                                      `json:"iid"`
	URL                 string                                     `json:"url"`
	State               string                                     `json:"state"`
	Action              string                                     `json:"action"`
	Severity            string                                     `json:"severity"`
	EscalationStatus    string                                     `json:"escalation_status"`
	EscalationPolicy    IssueEventObjectAttributesEscalationPolicy `json:"escalation_policy"`
	Labels              []*EventLabel                              `json:"labels"`
}

type IssueEventObjectAttributesEscalationPolicy struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type IssueEventChanges struct {
	Assignees      EventChangesAssignees           `json:"assignees"`
	Description    EventChangesDescription         `json:"description"`
	Labels         EventChangesLabels              `json:"labels"`
	Title          EventChangesTitle               `json:"title"`
	ClosedAt       IssueEventChangesClosedAt       `json:"closed_at"`
	StateID        EventChangesStateID             `json:"state_id"`
	UpdatedAt      EventChangesUpdatedAt           `json:"updated_at"`
	UpdatedByID    EventChangesUpdatedByID         `json:"updated_by_id"`
	TotalTimeSpent IssueEventChangesTotalTimeSpent `json:"total_time_spent"`
}

type EventChangesAssignees struct {
	Previous []*EventUser `json:"previous"`
	Current  []*EventUser `json:"current"`
}

type EventChangesDescription struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type EventChangesLabels struct {
	Previous []*EventLabel `json:"previous"`
	Current  []*EventLabel `json:"current"`
}

type EventChangesTitle struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type IssueEventChangesClosedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type EventChangesStateID struct {
	Previous StateID `json:"previous"`
	Current  StateID `json:"current"`
}

type EventChangesUpdatedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type EventChangesUpdatedByID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type IssueEventChangesTotalTimeSpent struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

// JobEvent represents a job event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#job-events
type JobEvent struct {
	ObjectKind          string              `json:"object_kind"`
	Ref                 string              `json:"ref"`
	Tag                 bool                `json:"tag"`
	BeforeSHA           string              `json:"before_sha"`
	SHA                 string              `json:"sha"`
	BuildID             int64               `json:"build_id"`
	BuildName           string              `json:"build_name"`
	BuildStage          string              `json:"build_stage"`
	BuildStatus         string              `json:"build_status"`
	BuildCreatedAt      string              `json:"build_created_at"`
	BuildStartedAt      string              `json:"build_started_at"`
	BuildFinishedAt     string              `json:"build_finished_at"`
	BuildDuration       float64             `json:"build_duration"`
	BuildQueuedDuration float64             `json:"build_queued_duration"`
	BuildAllowFailure   bool                `json:"build_allow_failure"`
	BuildFailureReason  string              `json:"build_failure_reason"`
	RetriesCount        int64               `json:"retries_count"`
	PipelineID          int64               `json:"pipeline_id"`
	ProjectID           int64               `json:"project_id"`
	ProjectName         string              `json:"project_name"`
	User                *EventUser          `json:"user"`
	Commit              JobEventCommit      `json:"commit"`
	Repository          *Repository         `json:"repository"`
	Runner              JobEventRunner      `json:"runner"`
	Environment         EventEnvironment    `json:"environment"`
	SourcePipeline      EventSourcePipeline `json:"source_pipeline"`
}

type JobEventCommit struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	Message     string `json:"message"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorURL   string `json:"author_url"`
	Status      string `json:"status"`
	Duration    int64  `json:"duration"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
}

type JobEventRunner struct {
	ID          int64    `json:"id"`
	Active      bool     `json:"active"`
	RunnerType  string   `json:"runner_type"`
	IsShared    bool     `json:"is_shared"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type EventEnvironment struct {
	Name           string `json:"name"`
	Action         string `json:"action"`
	DeploymentTier string `json:"deployment_tier"`
}

type EventSourcePipeline struct {
	Project    EventSourcePipelineProject `json:"project"`
	PipelineID int64                      `json:"pipeline_id"`
	JobID      int64                      `json:"job_id"`
}

type EventSourcePipelineProject struct {
	ID                int64  `json:"id"`
	WebURL            string `json:"web_url"`
	PathWithNamespace string `json:"path_with_namespace"`
}

// MemberEvent represents a member event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#group-member-events
type MemberEvent struct {
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	GroupName    string     `json:"group_name"`
	GroupPath    string     `json:"group_path"`
	GroupID      int64      `json:"group_id"`
	UserUsername string     `json:"user_username"`
	UserName     string     `json:"user_name"`
	UserEmail    string     `json:"user_email"`
	UserID       int64      `json:"user_id"`
	GroupAccess  string     `json:"group_access"`
	GroupPlan    string     `json:"group_plan"`
	ExpiresAt    *time.Time `json:"expires_at"`
	EventName    string     `json:"event_name"`
}

// MergeCommentEvent represents a comment on a merge event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#comment-on-a-merge-request
type MergeCommentEvent struct {
	ObjectKind       string                            `json:"object_kind"`
	EventType        string                            `json:"event_type"`
	User             *EventUser                        `json:"user"`
	ProjectID        int64                             `json:"project_id"`
	Project          MergeCommentEventProject          `json:"project"`
	ObjectAttributes MergeCommentEventObjectAttributes `json:"object_attributes"`
	Repository       *Repository                       `json:"repository"`
	MergeRequest     MergeCommentEventMergeRequest     `json:"merge_request"`
}

type MergeCommentEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type MergeCommentEventObjectAttributes struct {
	Attachment       string             `json:"attachment"`
	AuthorID         int64              `json:"author_id"`
	ChangePosition   *NotePosition      `json:"change_position"`
	CommitID         string             `json:"commit_id"`
	CreatedAt        string             `json:"created_at"`
	DiscussionID     string             `json:"discussion_id"`
	ID               int64              `json:"id"`
	LineCode         string             `json:"line_code"`
	Note             string             `json:"note"`
	NoteableID       int64              `json:"noteable_id"`
	NoteableType     string             `json:"noteable_type"`
	OriginalPosition *NotePosition      `json:"original_position"`
	Position         *NotePosition      `json:"position"`
	ProjectID        int64              `json:"project_id"`
	ResolvedAt       string             `json:"resolved_at"`
	ResolvedByID     int64              `json:"resolved_by_id"`
	ResolvedByPush   bool               `json:"resolved_by_push"`
	StDiff           *Diff              `json:"st_diff"`
	System           bool               `json:"system"`
	Type             string             `json:"type"`
	UpdatedAt        string             `json:"updated_at"`
	UpdatedByID      int64              `json:"updated_by_id"`
	Description      string             `json:"description"`
	Action           CommentEventAction `json:"action"`
	URL              string             `json:"url"`
}

type MergeCommentEventMergeRequest struct {
	ID                        int64                       `json:"id"`
	TargetBranch              string                      `json:"target_branch"`
	SourceBranch              string                      `json:"source_branch"`
	SourceProjectID           int64                       `json:"source_project_id"`
	AuthorID                  int64                       `json:"author_id"`
	AssigneeID                int64                       `json:"assignee_id"`
	AssigneeIDs               []int64                     `json:"assignee_ids"`
	ReviewerIDs               []int64                     `json:"reviewer_ids"`
	Title                     string                      `json:"title"`
	CreatedAt                 string                      `json:"created_at"`
	UpdatedAt                 string                      `json:"updated_at"`
	MilestoneID               int64                       `json:"milestone_id"`
	State                     string                      `json:"state"`
	MergeStatus               string                      `json:"merge_status"`
	TargetProjectID           int64                       `json:"target_project_id"`
	IID                       int64                       `json:"iid"`
	Description               string                      `json:"description"`
	Position                  int64                       `json:"position"`
	Labels                    []*EventLabel               `json:"labels"`
	LockedAt                  string                      `json:"locked_at"`
	UpdatedByID               int64                       `json:"updated_by_id"`
	MergeError                string                      `json:"merge_error"`
	MergeParams               *MergeParams                `json:"merge_params"`
	MergeWhenPipelineSucceeds bool                        `json:"merge_when_pipeline_succeeds"`
	MergeUserID               int64                       `json:"merge_user_id"`
	MergeCommitSHA            string                      `json:"merge_commit_sha"`
	DeletedAt                 string                      `json:"deleted_at"`
	InProgressMergeCommitSHA  string                      `json:"in_progress_merge_commit_sha"`
	LockVersion               int64                       `json:"lock_version"`
	ApprovalsBeforeMerge      string                      `json:"approvals_before_merge"`
	RebaseCommitSHA           string                      `json:"rebase_commit_sha"`
	TimeEstimate              int64                       `json:"time_estimate"`
	Squash                    bool                        `json:"squash"`
	LastEditedAt              string                      `json:"last_edited_at"`
	LastEditedByID            int64                       `json:"last_edited_by_id"`
	Source                    *Repository                 `json:"source"`
	Target                    *Repository                 `json:"target"`
	LastCommit                EventMergeRequestLastCommit `json:"last_commit"`
	WorkInProgress            bool                        `json:"work_in_progress"`
	TotalTimeSpent            int64                       `json:"total_time_spent"`
	HeadPipelineID            int64                       `json:"head_pipeline_id"`
	Assignee                  *EventUser                  `json:"assignee"`
	DetailedMergeStatus       string                      `json:"detailed_merge_status"`
	URL                       string                      `json:"url"`
}

type EventMergeRequestLastCommit struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
}

// MergeEvent represents a merge event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#merge-request-events
type MergeEvent struct {
	ObjectKind       string                     `json:"object_kind"`
	EventType        string                     `json:"event_type"`
	User             *EventUser                 `json:"user"`
	Project          MergeEventProject          `json:"project"`
	ObjectAttributes MergeEventObjectAttributes `json:"object_attributes"`
	Repository       *Repository                `json:"repository"`
	Labels           []*EventLabel              `json:"labels"`
	Changes          MergeEventChanges          `json:"changes"`
	Assignees        []*EventUser               `json:"assignees"`
	Reviewers        []*EventUser               `json:"reviewers"`
}

type MergeEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	CIConfigPath      string          `json:"ci_config_path"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type MergeEventObjectAttributes struct {
	ID                          int64                       `json:"id"`
	TargetBranch                string                      `json:"target_branch"`
	SourceBranch                string                      `json:"source_branch"`
	SourceProjectID             int64                       `json:"source_project_id"`
	AuthorID                    int64                       `json:"author_id"`
	AssigneeID                  int64                       `json:"assignee_id"`
	AssigneeIDs                 []int64                     `json:"assignee_ids"`
	ReviewerIDs                 []int64                     `json:"reviewer_ids"`
	Title                       string                      `json:"title"`
	CreatedAt                   string                      `json:"created_at"` // Should be *time.Time (see Gitlab issue #21468)
	UpdatedAt                   string                      `json:"updated_at"` // Should be *time.Time (see Gitlab issue #21468)
	StCommits                   []*Commit                   `json:"st_commits"`
	StDiffs                     []*Diff                     `json:"st_diffs"`
	LastEditedAt                string                      `json:"last_edited_at"`
	LastEditedByID              int64                       `json:"last_edited_by_id"`
	MilestoneID                 int64                       `json:"milestone_id"`
	StateID                     StateID                     `json:"state_id"`
	State                       string                      `json:"state"`
	MergeStatus                 string                      `json:"merge_status"`
	TargetProjectID             int64                       `json:"target_project_id"`
	IID                         int64                       `json:"iid"`
	Description                 string                      `json:"description"`
	Position                    int64                       `json:"position"`
	LockedAt                    string                      `json:"locked_at"`
	UpdatedByID                 int64                       `json:"updated_by_id"`
	MergeError                  string                      `json:"merge_error"`
	MergeParams                 *MergeParams                `json:"merge_params"`
	MergeWhenBuildSucceeds      bool                        `json:"merge_when_build_succeeds"`
	MergeUserID                 int64                       `json:"merge_user_id"`
	MergeCommitSHA              string                      `json:"merge_commit_sha"`
	DeletedAt                   string                      `json:"deleted_at"`
	ApprovalsBeforeMerge        string                      `json:"approvals_before_merge"`
	RebaseCommitSHA             string                      `json:"rebase_commit_sha"`
	InProgressMergeCommitSHA    string                      `json:"in_progress_merge_commit_sha"`
	LockVersion                 int64                       `json:"lock_version"`
	TimeEstimate                int64                       `json:"time_estimate"`
	Source                      *Repository                 `json:"source"`
	Target                      *Repository                 `json:"target"`
	HeadPipelineID              *int64                      `json:"head_pipeline_id"`
	LastCommit                  EventMergeRequestLastCommit `json:"last_commit"`
	BlockingDiscussionsResolved bool                        `json:"blocking_discussions_resolved"`
	WorkInProgress              bool                        `json:"work_in_progress"`
	Draft                       bool                        `json:"draft"`
	TotalTimeSpent              int64                       `json:"total_time_spent"`
	TimeChange                  int64                       `json:"time_change"`
	HumanTotalTimeSpent         string                      `json:"human_total_time_spent"`
	HumanTimeChange             string                      `json:"human_time_change"`
	HumanTimeEstimate           string                      `json:"human_time_estimate"`
	FirstContribution           bool                        `json:"first_contribution"`
	URL                         string                      `json:"url"`
	Labels                      []*EventLabel               `json:"labels"`
	Action                      string                      `json:"action"`
	DetailedMergeStatus         string                      `json:"detailed_merge_status"`
	OldRev                      string                      `json:"oldrev"`
	System                      bool                        `json:"system"`
	SystemAction                string                      `json:"system_action"`
}

type MergeEventChanges struct {
	Assignees       EventChangesAssignees            `json:"assignees"`
	Reviewers       MergeEventChangesReviewers       `json:"reviewers"`
	Description     EventChangesDescription          `json:"description"`
	Draft           MergeEventChangesDraft           `json:"draft"`
	Labels          EventChangesLabels               `json:"labels"`
	LastEditedAt    MergeEventChangesLastEditedAt    `json:"last_edited_at"`
	LastEditedByID  MergeEventChangesLastEditedByID  `json:"last_edited_by_id"`
	MergeStatus     MergeEventChangesMergeStatus     `json:"merge_status"`
	MilestoneID     MergeEventChangesMilestoneID     `json:"milestone_id"`
	SourceBranch    MergeEventChangesSourceBranch    `json:"source_branch"`
	SourceProjectID MergeEventChangesSourceProjectID `json:"source_project_id"`
	StateID         EventChangesStateID              `json:"state_id"`
	TargetBranch    MergeEventChangesTargetBranch    `json:"target_branch"`
	TargetProjectID MergeEventChangesTargetProjectID `json:"target_project_id"`
	Title           EventChangesTitle                `json:"title"`
	UpdatedAt       EventChangesUpdatedAt            `json:"updated_at"`
	UpdatedByID     EventChangesUpdatedByID          `json:"updated_by_id"`
}

type MergeEventChangesReviewers struct {
	Previous []*EventUser `json:"previous"`
	Current  []*EventUser `json:"current"`
}

type MergeEventChangesDraft struct {
	Previous bool `json:"previous"`
	Current  bool `json:"current"`
}

type MergeEventChangesLastEditedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type MergeEventChangesLastEditedByID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type MergeEventChangesMergeStatus struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type MergeEventChangesMilestoneID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type MergeEventChangesSourceBranch struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type MergeEventChangesSourceProjectID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type MergeEventChangesTargetBranch struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type MergeEventChangesTargetProjectID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

// EventUser represents a user record in an event and is used as an event
// initiator or a merge assignee.
type EventUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// MergeParams represents the merge params.
type MergeParams struct {
	ForceRemoveSourceBranch bool `json:"force_remove_source_branch"`
}

// UnmarshalJSON decodes the merge parameters
//
// This allows support of ForceRemoveSourceBranch for both type
// bool (>11.9) and string (<11.9)
func (p *MergeParams) UnmarshalJSON(b []byte) error {
	type Alias MergeParams
	raw := struct {
		*Alias
		ForceRemoveSourceBranch any `json:"force_remove_source_branch"`
	}{
		Alias: (*Alias)(p),
	}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	switch v := raw.ForceRemoveSourceBranch.(type) {
	case nil:
		// No action needed.
	case bool:
		p.ForceRemoveSourceBranch = v
	case string:
		p.ForceRemoveSourceBranch, err = strconv.ParseBool(v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to unmarshal ForceRemoveSourceBranch of type: %T", v)
	}

	return nil
}

// PipelineEvent represents a pipeline event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#pipeline-events
type PipelineEvent struct {
	ObjectKind       string                        `json:"object_kind"`
	ObjectAttributes PipelineEventObjectAttributes `json:"object_attributes"`
	MergeRequest     PipelineEventMergeRequest     `json:"merge_request"`
	User             *EventUser                    `json:"user"`
	Project          PipelineEventProject          `json:"project"`
	Commit           PipelineEventCommit           `json:"commit"`
	SourcePipeline   EventSourcePipeline           `json:"source_pipeline"`
	Builds           []PipelineEventBuild          `json:"builds"`
}

type PipelineEventObjectAttributes struct {
	ID             int64                                   `json:"id"`
	IID            int64                                   `json:"iid"`
	Name           string                                  `json:"name"`
	Ref            string                                  `json:"ref"`
	Tag            bool                                    `json:"tag"`
	SHA            string                                  `json:"sha"`
	BeforeSHA      string                                  `json:"before_sha"`
	Source         string                                  `json:"source"`
	Status         string                                  `json:"status"`
	DetailedStatus string                                  `json:"detailed_status"`
	Stages         []string                                `json:"stages"`
	CreatedAt      string                                  `json:"created_at"`
	FinishedAt     string                                  `json:"finished_at"`
	Duration       int64                                   `json:"duration"`
	QueuedDuration int64                                   `json:"queued_duration"`
	URL            string                                  `json:"url"`
	Variables      []PipelineEventObjectAttributesVariable `json:"variables"`
}

type PipelineEventObjectAttributesVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PipelineEventMergeRequest struct {
	ID                  int64  `json:"id"`
	IID                 int64  `json:"iid"`
	Title               string `json:"title"`
	SourceBranch        string `json:"source_branch"`
	SourceProjectID     int64  `json:"source_project_id"`
	TargetBranch        string `json:"target_branch"`
	TargetProjectID     int64  `json:"target_project_id"`
	State               string `json:"state"`
	MergeRequestStatus  string `json:"merge_status"`
	DetailedMergeStatus string `json:"detailed_merge_status"`
	URL                 string `json:"url"`
}

type PipelineEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type PipelineEventCommit struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Title     string            `json:"title"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
}

type PipelineEventBuild struct {
	ID             int64                           `json:"id"`
	Stage          string                          `json:"stage"`
	Name           string                          `json:"name"`
	Status         string                          `json:"status"`
	CreatedAt      string                          `json:"created_at"`
	StartedAt      string                          `json:"started_at"`
	FinishedAt     string                          `json:"finished_at"`
	Duration       float64                         `json:"duration"`
	QueuedDuration float64                         `json:"queued_duration"`
	FailureReason  string                          `json:"failure_reason"`
	When           string                          `json:"when"`
	Manual         bool                            `json:"manual"`
	AllowFailure   bool                            `json:"allow_failure"`
	User           *EventUser                      `json:"user"`
	Runner         PipelineEventBuildRunner        `json:"runner"`
	ArtifactsFile  PipelineEventBuildArtifactsFile `json:"artifacts_file"`
	Environment    EventEnvironment                `json:"environment"`
}

type PipelineEventBuildRunner struct {
	ID          int64    `json:"id"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	IsShared    bool     `json:"is_shared"`
	RunnerType  string   `json:"runner_type"`
	Tags        []string `json:"tags"`
}

type PipelineEventBuildArtifactsFile struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// ProjectResourceAccessTokenEvent represents a resource access token event for
// a project.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#project-and-group-access-token-events
type ProjectResourceAccessTokenEvent struct {
	EventName        string                                          `json:"event_name"`
	ObjectKind       string                                          `json:"object_kind"`
	Project          ProjectResourceAccessTokenEventProject          `json:"project"`
	ObjectAttributes ProjectResourceAccessTokenEventObjectAttributes `json:"object_attributes"`
}

type ProjectResourceAccessTokenEventProject struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CIConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

type ProjectResourceAccessTokenEventObjectAttributes struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	ExpiresAt *ISOTime `json:"expires_at"`
}

// PushEvent represents a push event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#push-events
type PushEvent struct {
	ObjectKind        string             `json:"object_kind"`
	EventName         string             `json:"event_name"`
	Before            string             `json:"before"`
	After             string             `json:"after"`
	Ref               string             `json:"ref"`
	RefProtected      bool               `json:"ref_protected"`
	CheckoutSHA       string             `json:"checkout_sha"`
	UserID            int64              `json:"user_id"`
	UserName          string             `json:"user_name"`
	UserUsername      string             `json:"user_username"`
	UserEmail         string             `json:"user_email"`
	UserAvatar        string             `json:"user_avatar"`
	ProjectID         int64              `json:"project_id"`
	Project           PushEventProject   `json:"project"`
	Repository        *Repository        `json:"repository"`
	Commits           []*PushEventCommit `json:"commits"`
	TotalCommitsCount int64              `json:"total_commits_count"`
}

type PushEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type PushEventCommit struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Title     string            `json:"title"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
	Added     []string          `json:"added"`
	Modified  []string          `json:"modified"`
	Removed   []string          `json:"removed"`
}

// ReleaseEvent represents a release event
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#release-events
type ReleaseEvent struct {
	ID          int64               `json:"id"`
	CreatedAt   string              `json:"created_at"` // Should be *time.Time (see Gitlab issue #21468)
	Description string              `json:"description"`
	Name        string              `json:"name"`
	Tag         string              `json:"tag"`
	ReleasedAt  string              `json:"released_at"` // Should be *time.Time (see Gitlab issue #21468)
	ObjectKind  string              `json:"object_kind"`
	Project     ReleaseEventProject `json:"project"`
	URL         string              `json:"url"`
	Action      string              `json:"action"`
	Assets      ReleaseEventAssets  `json:"assets"`
	Commit      ReleaseEventCommit  `json:"commit"`
}

type ReleaseEventProject struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	WebURL            string  `json:"web_url"`
	AvatarURL         *string `json:"avatar_url"`
	GitSSHURL         string  `json:"git_ssh_url"`
	GitHTTPURL        string  `json:"git_http_url"`
	Namespace         string  `json:"namespace"`
	VisibilityLevel   int64   `json:"visibility_level"`
	PathWithNamespace string  `json:"path_with_namespace"`
	DefaultBranch     string  `json:"default_branch"`
	CIConfigPath      string  `json:"ci_config_path"`
	Homepage          string  `json:"homepage"`
	URL               string  `json:"url"`
	SSHURL            string  `json:"ssh_url"`
	HTTPURL           string  `json:"http_url"`
}

type ReleaseEventAssets struct {
	Count   int64                      `json:"count"`
	Links   []ReleaseEventAssetsLink   `json:"links"`
	Sources []ReleaseEventAssetsSource `json:"sources"`
}

type ReleaseEventAssetsLink struct {
	ID       int64  `json:"id"`
	External bool   `json:"external"`
	LinkType string `json:"link_type"`
	Name     string `json:"name"`
	URL      string `json:"url"`
}

type ReleaseEventAssetsSource struct {
	Format string `json:"format"`
	URL    string `json:"url"`
}

type ReleaseEventCommit struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Title     string            `json:"title"`
	Timestamp string            `json:"timestamp"` // Should be *time.Time (see Gitlab issue #21468)
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
}

// SnippetCommentEvent represents a comment on a snippet event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#comment-on-a-code-snippet
type SnippetCommentEvent struct {
	ObjectKind       string                              `json:"object_kind"`
	EventType        string                              `json:"event_type"`
	User             *EventUser                          `json:"user"`
	ProjectID        int64                               `json:"project_id"`
	Project          SnippetCommentEventProject          `json:"project"`
	Repository       *Repository                         `json:"repository"`
	ObjectAttributes SnippetCommentEventObjectAttributes `json:"object_attributes"`
	Snippet          *SnippetCommentEventSnippet         `json:"snippet"`
}

type SnippetCommentEventProject struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type SnippetCommentEventObjectAttributes struct {
	ID           int64              `json:"id"`
	Note         string             `json:"note"`
	NoteableType string             `json:"noteable_type"`
	AuthorID     int64              `json:"author_id"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
	ProjectID    int64              `json:"project_id"`
	Attachment   string             `json:"attachment"`
	LineCode     string             `json:"line_code"`
	CommitID     string             `json:"commit_id"`
	NoteableID   int64              `json:"noteable_id"`
	System       bool               `json:"system"`
	StDiff       *Diff              `json:"st_diff"`
	Description  string             `json:"description"`
	Action       CommentEventAction `json:"action"`
	URL          string             `json:"url"`
}

type SnippetCommentEventSnippet struct {
	ID                 int64  `json:"id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	AuthorID           int64  `json:"author_id"`
	ProjectID          int64  `json:"project_id"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	Filename           string `json:"file_name"`
	ExpiresAt          string `json:"expires_at"`
	Type               string `json:"type"`
	VisibilityLevel    int64  `json:"visibility_level"`
	Description        string `json:"description"`
	Secret             bool   `json:"secret"`
	RepositoryReadOnly bool   `json:"repository_read_only"`
}

// SubGroupEvent represents a subgroup event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#subgroup-events
type SubGroupEvent struct {
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	EventName      string     `json:"event_name"`
	Name           string     `json:"name"`
	Path           string     `json:"path"`
	FullPath       string     `json:"full_path"`
	GroupID        int64      `json:"group_id"`
	ParentGroupID  int64      `json:"parent_group_id"`
	ParentName     string     `json:"parent_name"`
	ParentPath     string     `json:"parent_path"`
	ParentFullPath string     `json:"parent_full_path"`
}

// TagEvent represents a tag event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#tag-events
type TagEvent struct {
	ObjectKind        string            `json:"object_kind"`
	EventName         string            `json:"event_name"`
	Before            string            `json:"before"`
	After             string            `json:"after"`
	Ref               string            `json:"ref"`
	CheckoutSHA       string            `json:"checkout_sha"`
	UserID            int64             `json:"user_id"`
	UserName          string            `json:"user_name"`
	UserUsername      string            `json:"user_username"`
	UserAvatar        string            `json:"user_avatar"`
	UserEmail         string            `json:"user_email"`
	ProjectID         int64             `json:"project_id"`
	Message           string            `json:"message"`
	Project           TagEventProject   `json:"project"`
	Repository        *Repository       `json:"repository"`
	Commits           []*TagEventCommit `json:"commits"`
	TotalCommitsCount int64             `json:"total_commits_count"`
}

type TagEventProject struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type TagEventCommit struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Title     string            `json:"title"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
	Added     []string          `json:"added"`
	Modified  []string          `json:"modified"`
	Removed   []string          `json:"removed"`
}

// WikiPageEvent represents a wiki page event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#wiki-page-events
type WikiPageEvent struct {
	ObjectKind       string                        `json:"object_kind"`
	User             *EventUser                    `json:"user"`
	Project          WikiPageEventProject          `json:"project"`
	Wiki             WikiPageEventWiki             `json:"wiki"`
	ObjectAttributes WikiPageEventObjectAttributes `json:"object_attributes"`
}

type WikiPageEventProject struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
}

type WikiPageEventWiki struct {
	WebURL            string `json:"web_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
}

type WikiPageEventObjectAttributes struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Format  string `json:"format"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	URL     string `json:"url"`
	Action  string `json:"action"`
	DiffURL string `json:"diff_url"`
}

// EmojiEvent represents an emoji event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#emoji-events
type EmojiEvent struct {
	ObjectKind       string                     `json:"object_kind"`
	EventType        string                     `json:"event_type"`
	User             EventUser                  `json:"user"`
	ProjectID        int64                      `json:"project_id"`
	Project          EmojiEventProject          `json:"project"`
	ObjectAttributes EmojiEventObjectAttributes `json:"object_attributes"`
	Note             *EmojiEventNote            `json:"note,omitempty"`
	Issue            *EmojiEventIssue           `json:"issue,omitempty"`
	MergeRequest     *EmojiEventMergeRequest    `json:"merge_request,omitempty"`
	ProjectSnippet   *EmojiEventSnippet         `json:"project_snippet,omitempty"`
	Commit           *EmojiEventCommit          `json:"commit,omitempty"`
}

// EmojiEventProject represents a project in an emoji event.
type EmojiEventProject struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CIConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

type EmojiEventObjectAttributes struct {
	UserID        int64  `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	AwardableType string `json:"awardable_type"`
	AwardableID   int64  `json:"awardable_id"`
	UpdatedAt     string `json:"updated_at"`
	Action        string `json:"action"`
	AwardedOnURL  string `json:"awarded_on_url"`
}

type EmojiEventNote struct {
	Attachment       *string       `json:"attachment"`
	AuthorID         int64         `json:"author_id"`
	ChangePosition   *NotePosition `json:"change_position"`
	CommitID         *string       `json:"commit_id"`
	CreatedAt        string        `json:"created_at"`
	DiscussionID     string        `json:"discussion_id"`
	ID               int64         `json:"id"`
	LineCode         *string       `json:"line_code"`
	Note             string        `json:"note"`
	NoteableID       int64         `json:"noteable_id"`
	NoteableType     string        `json:"noteable_type"`
	OriginalPosition *NotePosition `json:"original_position"`
	Position         *NotePosition `json:"position"`
	ProjectID        int64         `json:"project_id"`
	ResolvedAt       *string       `json:"resolved_at"`
	ResolvedByID     *int64        `json:"resolved_by_id"`
	ResolvedByPush   *bool         `json:"resolved_by_push"`
	StDiff           *Diff         `json:"st_diff"`
	System           bool          `json:"system"`
	Type             *string       `json:"type"`
	UpdatedAt        string        `json:"updated_at"`
	UpdatedByID      *int64        `json:"updated_by_id"`
	Description      string        `json:"description"`
	URL              string        `json:"url"`
}

type EmojiEventIssue struct {
	ID                  int64         `json:"id"`
	IID                 int64         `json:"iid"`
	ProjectID           int64         `json:"project_id"`
	AuthorID            int64         `json:"author_id"`
	ClosedAt            *string       `json:"closed_at"`
	Confidential        bool          `json:"confidential"`
	CreatedAt           string        `json:"created_at"`
	Description         string        `json:"description"`
	DiscussionLocked    *bool         `json:"discussion_locked"`
	DueDate             *ISOTime      `json:"due_date"`
	LastEditedAt        *string       `json:"last_edited_at"`
	LastEditedByID      *int64        `json:"last_edited_by_id"`
	MilestoneID         *int64        `json:"milestone_id"`
	MovedToID           *int64        `json:"moved_to_id"`
	DuplicatedToID      *int64        `json:"duplicated_to_id"`
	RelativePosition    int64         `json:"relative_position"`
	StateID             StateID       `json:"state_id"`
	TimeEstimate        int64         `json:"time_estimate"`
	Title               string        `json:"title"`
	UpdatedAt           string        `json:"updated_at"`
	UpdatedByID         *int64        `json:"updated_by_id"`
	Weight              *int64        `json:"weight"`
	HealthStatus        *string       `json:"health_status"`
	URL                 string        `json:"url"`
	TotalTimeSpent      int64         `json:"total_time_spent"`
	TimeChange          int64         `json:"time_change"`
	HumanTotalTimeSpent *string       `json:"human_total_time_spent"`
	HumanTimeChange     *string       `json:"human_time_change"`
	HumanTimeEstimate   *string       `json:"human_time_estimate"`
	AssigneeIDs         []int64       `json:"assignee_ids"`
	AssigneeID          *int64        `json:"assignee_id"`
	Labels              []*EventLabel `json:"labels"`
	State               string        `json:"state"`
	Severity            string        `json:"severity"`
}

// EmojiEventMergeRequest represents a merge request in an emoji event.
type EmojiEventMergeRequest struct {
	ID                        int64                       `json:"id"`
	TargetBranch              string                      `json:"target_branch"`
	SourceBranch              string                      `json:"source_branch"`
	SourceProjectID           int64                       `json:"source_project_id"`
	AuthorID                  int64                       `json:"author_id"`
	AssigneeID                int64                       `json:"assignee_id"`
	AssigneeIDs               []int64                     `json:"assignee_ids"`
	ReviewerIDs               []int64                     `json:"reviewer_ids"`
	Title                     string                      `json:"title"`
	CreatedAt                 string                      `json:"created_at"`
	UpdatedAt                 string                      `json:"updated_at"`
	MilestoneID               int64                       `json:"milestone_id"`
	State                     string                      `json:"state"`
	MergeStatus               string                      `json:"merge_status"`
	TargetProjectID           int64                       `json:"target_project_id"`
	IID                       int64                       `json:"iid"`
	Description               string                      `json:"description"`
	Position                  int64                       `json:"position"`
	Labels                    []*EventLabel               `json:"labels"`
	LockedAt                  string                      `json:"locked_at"`
	UpdatedByID               int64                       `json:"updated_by_id"`
	MergeError                string                      `json:"merge_error"`
	MergeParams               *MergeParams                `json:"merge_params"`
	MergeWhenPipelineSucceeds bool                        `json:"merge_when_pipeline_succeeds"`
	MergeUserID               int64                       `json:"merge_user_id"`
	MergeCommitSHA            string                      `json:"merge_commit_sha"`
	DeletedAt                 string                      `json:"deleted_at"`
	InProgressMergeCommitSHA  string                      `json:"in_progress_merge_commit_sha"`
	LockVersion               int64                       `json:"lock_version"`
	ApprovalsBeforeMerge      string                      `json:"approvals_before_merge"`
	RebaseCommitSHA           string                      `json:"rebase_commit_sha"`
	TimeEstimate              int64                       `json:"time_estimate"`
	Squash                    bool                        `json:"squash"`
	LastEditedAt              string                      `json:"last_edited_at"`
	LastEditedByID            int64                       `json:"last_edited_by_id"`
	Source                    *Repository                 `json:"source"`
	Target                    *Repository                 `json:"target"`
	LastCommit                EventMergeRequestLastCommit `json:"last_commit"`
	WorkInProgress            bool                        `json:"work_in_progress"`
	TotalTimeSpent            int64                       `json:"total_time_spent"`
	HeadPipelineID            int64                       `json:"head_pipeline_id"`
	Assignee                  *EventUser                  `json:"assignee"`
	DetailedMergeStatus       string                      `json:"detailed_merge_status"`
	URL                       string                      `json:"url"`
}

// EmojiEventSnippet represents a snippet in an emoji event.
type EmojiEventSnippet struct {
	ID                 int64  `json:"id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	AuthorID           int64  `json:"author_id"`
	ProjectID          int64  `json:"project_id"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	Filename           string `json:"file_name"`
	ExpiresAt          string `json:"expires_at"`
	Type               string `json:"type"`
	VisibilityLevel    int64  `json:"visibility_level"`
	Description        string `json:"description"`
	Secret             bool   `json:"secret"`
	RepositoryReadOnly bool   `json:"repository_read_only"`
}

// EmojiEventCommit represents a commit in an emoji event.
type EmojiEventCommit struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Timestamp *time.Time        `json:"timestamp"`
	URL       string            `json:"url"`
	Author    EventCommitAuthor `json:"author"`
}

// MilestoneWebhookEvent represents a milestone webhook event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#milestone-events
type MilestoneWebhookEvent struct {
	ObjectKind       string                         `json:"object_kind"`
	EventType        string                         `json:"event_type"`
	Project          MilestoneEventProject          `json:"project"`
	Group            *MilestoneEventGroup           `json:"group,omitempty"`
	ObjectAttributes MilestoneEventObjectAttributes `json:"object_attributes"`
	Action           string                         `json:"action"`
}

type MilestoneEventObjectAttributes struct {
	ID          int64    `json:"id"`
	IID         int64    `json:"iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	DueDate     *ISOTime `json:"due_date"`
	StartDate   *ISOTime `json:"start_date"`
	GroupID     *int64   `json:"group_id"`
	ProjectID   int64    `json:"project_id"`
}

// MilestoneEventGroup represents a group in a milestone event.
type MilestoneEventGroup struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	GroupPath string `json:"group_path"`
	FullPath  string `json:"full_path"`
}

// MilestoneEventProject represents a project in a milestone event.
type MilestoneEventProject struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CIConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

// ProjectWebhookEvent represents a project webhook event for group webhooks.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#project-events
type ProjectWebhookEvent struct {
	EventName            string              `json:"event_name"`
	CreatedAt            string              `json:"created_at"`
	UpdatedAt            string              `json:"updated_at"`
	Name                 string              `json:"name"`
	Path                 string              `json:"path"`
	PathWithNamespace    string              `json:"path_with_namespace"`
	ProjectID            int64               `json:"project_id"`
	ProjectNamespaceID   int64               `json:"project_namespace_id"`
	Owners               []ProjectEventOwner `json:"owners"`
	ProjectVisibility    string              `json:"project_visibility"`
	OldPathWithNamespace string              `json:"old_path_with_namespace,omitempty"`
}

type ProjectEventOwner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// VulnerabilityEvent represents a vulnerability event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#vulnerability-events
type VulnerabilityEvent struct {
	ObjectKind       string                             `json:"object_kind"`
	ObjectAttributes VulnerabilityEventObjectAttributes `json:"object_attributes"`
}

type VulnerabilityEventObjectAttributes struct {
	ID                      int64                          `json:"id"`
	URL                     string                         `json:"url"`
	Title                   string                         `json:"title"`
	State                   string                         `json:"state"`
	ProjectID               int64                          `json:"project_id"`
	Location                VulnerabilityEventLocation     `json:"location"`
	CVSS                    []VulnerabilityEventCVSS       `json:"cvss"`
	Severity                string                         `json:"severity"`
	SeverityOverridden      bool                           `json:"severity_overridden"`
	Identifiers             []VulnerabilityEventIdentifier `json:"identifiers"`
	Issues                  []VulnerabilityEventIssue      `json:"issues"`
	ReportType              string                         `json:"report_type"`
	Confidence              string                         `json:"confidence"`
	ConfidenceOverridden    bool                           `json:"confidence_overridden"`
	ConfirmedAt             string                         `json:"confirmed_at"`
	ConfirmedByID           int64                          `json:"confirmed_by_id"`
	DismissedAt             string                         `json:"dismissed_at"`
	DismissedByID           int64                          `json:"dismissed_by_id"`
	ResolvedAt              string                         `json:"resolved_at"`
	ResolvedByID            int64                          `json:"resolved_by_id"`
	AutoResolved            bool                           `json:"auto_resolved"`
	ResolvedOnDefaultBranch bool                           `json:"resolved_on_default_branch"`
	CreatedAt               string                         `json:"created_at"`
	UpdatedAt               string                         `json:"updated_at"`
}

type VulnerabilityEventLocation struct {
	File       string                               `json:"file"`
	Dependency VulnerabilityEventLocationDependency `json:"dependency"`
}

type VulnerabilityEventLocationDependency struct {
	Package VulnerabilityEventLocationDependencyPackage `json:"package"`
	Version string                                      `json:"version"`
}

type VulnerabilityEventLocationDependencyPackage struct {
	Name string `json:"name"`
}

type VulnerabilityEventCVSS struct {
	Vector string `json:"vector"`
	Vendor string `json:"vendor"`
}

type VulnerabilityEventIdentifier struct {
	Name         string `json:"name"`
	ExternalID   string `json:"external_id"`
	ExternalType string `json:"external_type"`
	URL          string `json:"url"`
}

type VulnerabilityEventIssue struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// EventLabel represents a label inside a webhook event.
//
// GitLab API docs:
// https://docs.gitlab.com/user/project/integrations/webhook_events/#work-item-events
type EventLabel struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	ProjectID   int64  `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int64  `json:"group_id"`
}
