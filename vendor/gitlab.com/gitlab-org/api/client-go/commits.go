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
	"errors"
	"net/http"
	"time"
)

// CommitsService handles communication with the commit related methods
// of the GitLab API.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/
type (
	CommitsServiceInterface interface {
		// ListCommits gets a list of repository commits in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#list-repository-commits
		ListCommits(pid any, opt *ListCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)

		// GetCommitRefs gets all references (from branches or tags) a commit is pushed to.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#get-references-a-commit-is-pushed-to
		GetCommitRefs(pid any, sha string, opt *GetCommitRefsOptions, options ...RequestOptionFunc) ([]*CommitRef, *Response, error)

		// GetCommit gets a specific commit identified by the commit hash or name of a
		// branch or tag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#get-a-single-commit
		GetCommit(pid any, sha string, opt *GetCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error)

		// CreateCommit creates a commit with multiple files and actions.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#create-a-commit-with-multiple-files-and-actions
		CreateCommit(pid any, opt *CreateCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error)

		// GetCommitDiff gets the diff of a commit in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#get-the-diff-of-a-commit
		GetCommitDiff(pid any, sha string, opt *GetCommitDiffOptions, options ...RequestOptionFunc) ([]*Diff, *Response, error)

		// GetCommitComments gets the comments of a commit in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#get-the-comments-of-a-commit
		GetCommitComments(pid any, sha string, opt *GetCommitCommentsOptions, options ...RequestOptionFunc) ([]*CommitComment, *Response, error)

		// PostCommitComment adds a comment to a commit. Optionally you can post
		// comments on a specific line of a commit. Therefore both path, line_new and
		// line_old are required.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#post-comment-to-commit
		PostCommitComment(pid any, sha string, opt *PostCommitCommentOptions, options ...RequestOptionFunc) (*CommitComment, *Response, error)

		// GetCommitStatuses gets the statuses of a commit in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#list-the-statuses-of-a-commit
		GetCommitStatuses(pid any, sha string, opt *GetCommitStatusesOptions, options ...RequestOptionFunc) ([]*CommitStatus, *Response, error)

		// SetCommitStatus sets the status of a commit in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#set-the-pipeline-status-of-a-commit
		SetCommitStatus(pid any, sha string, opt *SetCommitStatusOptions, options ...RequestOptionFunc) (*CommitStatus, *Response, error)

		// ListMergeRequestsByCommit gets merge request associated with a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#list-merge-requests-associated-with-a-commit
		ListMergeRequestsByCommit(pid any, sha string, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)

		// CherryPickCommit cherry picks a commit to a given branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#cherry-pick-a-commit
		CherryPickCommit(pid any, sha string, opt *CherryPickCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error)

		// RevertCommit reverts a commit in a given branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#revert-a-commit
		RevertCommit(pid any, sha string, opt *RevertCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error)

		// GetGPGSignature gets a GPG signature of a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/commits/#get-signature-of-a-commit
		GetGPGSignature(pid any, sha string, options ...RequestOptionFunc) (*GPGSignature, *Response, error)
	}

	// CommitsService handles communication with the commit related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/commits/
	CommitsService struct {
		client *Client
	}
)

var _ CommitsServiceInterface = (*CommitsService)(nil)

// Commit represents a GitLab commit.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/
type Commit struct {
	ID               string            `json:"id"`
	ShortID          string            `json:"short_id"`
	Title            string            `json:"title"`
	AuthorName       string            `json:"author_name"`
	AuthorEmail      string            `json:"author_email"`
	AuthoredDate     *time.Time        `json:"authored_date"`
	CommitterName    string            `json:"committer_name"`
	CommitterEmail   string            `json:"committer_email"`
	CommittedDate    *time.Time        `json:"committed_date"`
	CreatedAt        *time.Time        `json:"created_at"`
	Message          string            `json:"message"`
	ParentIDs        []string          `json:"parent_ids"`
	Stats            *CommitStats      `json:"stats"`
	Status           *BuildStateValue  `json:"status"`
	LastPipeline     *PipelineInfo     `json:"last_pipeline"`
	ProjectID        int64             `json:"project_id"`
	Trailers         map[string]string `json:"trailers"`
	ExtendedTrailers map[string]string `json:"extended_trailers"`
	WebURL           string            `json:"web_url"`
}

// CommitStats represents the number of added and deleted files in a commit.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/
type CommitStats struct {
	Additions int64 `json:"additions"`
	Deletions int64 `json:"deletions"`
	Total     int64 `json:"total"`
}

func (c Commit) String() string {
	return Stringify(c)
}

// ListCommitsOptions represents the available ListCommits() options.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#list-repository-commits
type ListCommitsOptions struct {
	ListOptions
	RefName     *string    `url:"ref_name,omitempty" json:"ref_name,omitempty"`
	Since       *time.Time `url:"since,omitempty" json:"since,omitempty"`
	Until       *time.Time `url:"until,omitempty" json:"until,omitempty"`
	Path        *string    `url:"path,omitempty" json:"path,omitempty"`
	Author      *string    `url:"author,omitempty" json:"author,omitempty"`
	All         *bool      `url:"all,omitempty" json:"all,omitempty"`
	WithStats   *bool      `url:"with_stats,omitempty" json:"with_stats,omitempty"`
	FirstParent *bool      `url:"first_parent,omitempty" json:"first_parent,omitempty"`
	Trailers    *bool      `url:"trailers,omitempty" json:"trailers,omitempty"`
}

func (s *CommitsService) ListCommits(pid any, opt *ListCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CommitRef represents the reference of branches/tags in a commit.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-references-a-commit-is-pushed-to
type CommitRef struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// GetCommitRefsOptions represents the available GetCommitRefs() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-references-a-commit-is-pushed-to
type GetCommitRefsOptions struct {
	ListOptions
	Type *string `url:"type,omitempty" json:"type,omitempty"`
}

func (s *CommitsService) GetCommitRefs(pid any, sha string, opt *GetCommitRefsOptions, options ...RequestOptionFunc) ([]*CommitRef, *Response, error) {
	return do[[]*CommitRef](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/refs", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetCommitOptions represents the available GetCommit() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-a-single-commit
type GetCommitOptions struct {
	Stats *bool `url:"stats,omitempty" json:"stats,omitempty"`
}

func (s *CommitsService) GetCommit(pid any, sha string, opt *GetCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error) {
	if sha == "" {
		return nil, nil, errors.New("SHA must be a non-empty string")
	}

	return do[*Commit](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateCommitOptions represents the available options for a new commit.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#create-a-commit-with-multiple-files-and-actions
type CreateCommitOptions struct {
	Branch        *string                `url:"branch,omitempty" json:"branch,omitempty"`
	CommitMessage *string                `url:"commit_message,omitempty" json:"commit_message,omitempty"`
	StartBranch   *string                `url:"start_branch,omitempty" json:"start_branch,omitempty"`
	StartSHA      *string                `url:"start_sha,omitempty" json:"start_sha,omitempty"`
	StartProject  *string                `url:"start_project,omitempty" json:"start_project,omitempty"`
	Actions       []*CommitActionOptions `url:"actions" json:"actions"`
	AuthorEmail   *string                `url:"author_email,omitempty" json:"author_email,omitempty"`
	AuthorName    *string                `url:"author_name,omitempty" json:"author_name,omitempty"`
	Stats         *bool                  `url:"stats,omitempty" json:"stats,omitempty"`
	Force         *bool                  `url:"force,omitempty" json:"force,omitempty"`
}

// CommitActionOptions represents the available options for a new single
// file action.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#create-a-commit-with-multiple-files-and-actions
type CommitActionOptions struct {
	Action          *FileActionValue `url:"action,omitempty" json:"action,omitempty"`
	FilePath        *string          `url:"file_path,omitempty" json:"file_path,omitempty"`
	PreviousPath    *string          `url:"previous_path,omitempty" json:"previous_path,omitempty"`
	Content         *string          `url:"content,omitempty" json:"content,omitempty"`
	Encoding        *string          `url:"encoding,omitempty" json:"encoding,omitempty"`
	LastCommitID    *string          `url:"last_commit_id,omitempty" json:"last_commit_id,omitempty"`
	ExecuteFilemode *bool            `url:"execute_filemode,omitempty" json:"execute_filemode,omitempty"`
}

func (s *CommitsService) CreateCommit(pid any, opt *CreateCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error) {
	return do[*Commit](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// Diff represents a GitLab diff.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/
type Diff struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"new_path"`
	OldPath     string `json:"old_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

func (d Diff) String() string {
	return Stringify(d)
}

// GetCommitDiffOptions represents the available GetCommitDiff() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-the-diff-of-a-commit
type GetCommitDiffOptions struct {
	ListOptions
	Unidiff *bool `url:"unidiff,omitempty" json:"unidiff,omitempty"`
}

func (s *CommitsService) GetCommitDiff(pid any, sha string, opt *GetCommitDiffOptions, options ...RequestOptionFunc) ([]*Diff, *Response, error) {
	return do[[]*Diff](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/diff", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CommitComment represents a GitLab commit comment.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/
type CommitComment struct {
	Note     string `json:"note"`
	Path     string `json:"path"`
	Line     int64  `json:"line"`
	LineType string `json:"line_type"`
	Author   Author `json:"author"`
}

// Author represents a GitLab commit author
type Author struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	Blocked   bool       `json:"blocked"`
	CreatedAt *time.Time `json:"created_at"`
}

func (c CommitComment) String() string {
	return Stringify(c)
}

// GetCommitCommentsOptions represents the available GetCommitComments() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-the-comments-of-a-commit
type GetCommitCommentsOptions struct {
	ListOptions
}

func (s *CommitsService) GetCommitComments(pid any, sha string, opt *GetCommitCommentsOptions, options ...RequestOptionFunc) ([]*CommitComment, *Response, error) {
	return do[[]*CommitComment](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/comments", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// PostCommitCommentOptions represents the available PostCommitComment()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#post-comment-to-commit
type PostCommitCommentOptions struct {
	Note     *string `url:"note,omitempty" json:"note,omitempty"`
	Path     *string `url:"path" json:"path"`
	Line     *int64  `url:"line" json:"line"`
	LineType *string `url:"line_type" json:"line_type"`
}

func (s *CommitsService) PostCommitComment(pid any, sha string, opt *PostCommitCommentOptions, options ...RequestOptionFunc) (*CommitComment, *Response, error) {
	return do[*CommitComment](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits/%s/comments", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetCommitStatusesOptions represents the available GetCommitStatuses() options.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#list-the-statuses-of-a-commit
type GetCommitStatusesOptions struct {
	ListOptions
	Ref        *string `url:"ref,omitempty" json:"ref,omitempty"`
	Stage      *string `url:"stage,omitempty" json:"stage,omitempty"`
	Name       *string `url:"name,omitempty" json:"name,omitempty"`
	PipelineID *int64  `url:"pipeline_id,omitempty" json:"pipeline_id,omitempty"`
	All        *bool   `url:"all,omitempty" json:"all,omitempty"`
}

// CommitStatus represents a GitLab commit status.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#commit-status
type CommitStatus struct {
	ID           int64      `json:"id"`
	SHA          string     `json:"sha"`
	Ref          string     `json:"ref"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	Name         string     `json:"name"`
	AllowFailure bool       `json:"allow_failure"`
	Coverage     float64    `json:"coverage"`
	PipelineID   int64      `json:"pipeline_id"`
	Author       Author     `json:"author"`
	Description  string     `json:"description"`
	TargetURL    string     `json:"target_url"`
}

func (s *CommitsService) GetCommitStatuses(pid any, sha string, opt *GetCommitStatusesOptions, options ...RequestOptionFunc) ([]*CommitStatus, *Response, error) {
	return do[[]*CommitStatus](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/statuses", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// SetCommitStatusOptions represents the available SetCommitStatus() options.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#set-the-pipeline-status-of-a-commit
type SetCommitStatusOptions struct {
	State       BuildStateValue `url:"state" json:"state"`
	Ref         *string         `url:"ref,omitempty" json:"ref,omitempty"`
	Name        *string         `url:"name,omitempty" json:"name,omitempty"`
	Context     *string         `url:"context,omitempty" json:"context,omitempty"`
	TargetURL   *string         `url:"target_url,omitempty" json:"target_url,omitempty"`
	Description *string         `url:"description,omitempty" json:"description,omitempty"`
	Coverage    *float64        `url:"coverage,omitempty" json:"coverage,omitempty"`
	PipelineID  *int64          `url:"pipeline_id,omitempty" json:"pipeline_id,omitempty"`
}

func (s *CommitsService) SetCommitStatus(pid any, sha string, opt *SetCommitStatusOptions, options ...RequestOptionFunc) (*CommitStatus, *Response, error) {
	return do[*CommitStatus](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/statuses/%s", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *CommitsService) ListMergeRequestsByCommit(pid any, sha string, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	return do[[]*BasicMergeRequest](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/merge_requests", ProjectID{pid}, sha),
		withRequestOpts(options...),
	)
}

// CherryPickCommitOptions represents the available CherryPickCommit() options.
//
// GitLab API docs: https://docs.gitlab.com/api/commits/#cherry-pick-a-commit
type CherryPickCommitOptions struct {
	Branch  *string `url:"branch,omitempty" json:"branch,omitempty"`
	DryRun  *bool   `url:"dry_run,omitempty" json:"dry_run,omitempty"`
	Message *string `url:"message,omitempty" json:"message,omitempty"`
}

// RevertCommitOptions represents the available RevertCommit() options.
// GitLab API docs: https://docs.gitlab.com/api/commits/#revert-a-commit
type RevertCommitOptions struct {
	Branch *string `url:"branch,omitempty" json:"branch,omitempty"`
}

func (s *CommitsService) CherryPickCommit(pid any, sha string, opt *CherryPickCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error) {
	return do[*Commit](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits/%s/cherry_pick", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *CommitsService) RevertCommit(pid any, sha string, opt *RevertCommitOptions, options ...RequestOptionFunc) (*Commit, *Response, error) {
	return do[*Commit](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits/%s/revert", ProjectID{pid}, sha),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GPGSignature represents a GitLab commit's GPG Signature.
//
// GitLab API docs:
// https://docs.gitlab.com/api/commits/#get-signature-of-a-commit
type GPGSignature struct {
	KeyID              int64  `json:"gpg_key_id"`
	KeyPrimaryKeyID    string `json:"gpg_key_primary_keyid"`
	KeyUserName        string `json:"gpg_key_user_name"`
	KeyUserEmail       string `json:"gpg_key_user_email"`
	VerificationStatus string `json:"verification_status"`
	KeySubkeyID        int64  `json:"gpg_key_subkey_id"`
}

func (s *CommitsService) GetGPGSignature(pid any, sha string, options ...RequestOptionFunc) (*GPGSignature, *Response, error) {
	return do[*GPGSignature](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/commits/%s/signature", ProjectID{pid}, sha),
		withRequestOpts(options...),
	)
}
