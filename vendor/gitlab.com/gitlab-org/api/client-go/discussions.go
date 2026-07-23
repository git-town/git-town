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
	"net/http"
	"time"
)

type (
	// DiscussionsServiceInterface defines all the API methods for the DiscussionsService
	DiscussionsServiceInterface interface {
		// ListIssueDiscussions gets a list of all discussions for a single issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#list-project-issue-discussion-items
		ListIssueDiscussions(pid any, issue int64, opt *ListIssueDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error)

		// GetIssueDiscussion returns a single discussion for a specific project issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-issue-discussion-item
		GetIssueDiscussion(pid any, issue int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// CreateIssueDiscussion creates a new discussion to a single project issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-issue-thread
		CreateIssueDiscussion(pid any, issue int64, opt *CreateIssueDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// AddIssueDiscussionNote creates a new note in an existing discussion of an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-issue-thread
		AddIssueDiscussionNote(pid any, issue int64, discussion string, opt *AddIssueDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// UpdateIssueDiscussionNote modifies an existing note in a discussion of an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#modify-existing-issue-thread-note
		UpdateIssueDiscussionNote(pid any, issue int64, discussion string, note int64, opt *UpdateIssueDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// DeleteIssueDiscussionNote deletes a note from a discussion of an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#delete-an-issue-thread-note
		DeleteIssueDiscussionNote(pid any, issue int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error)

		// ListSnippetDiscussions gets all discussions for a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#list-project-snippet-discussion-items
		ListSnippetDiscussions(pid any, snippet int64, opt *ListSnippetDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error)

		// GetSnippetDiscussion returns a single discussion for a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-snippet-discussion-item
		GetSnippetDiscussion(pid any, snippet int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// CreateSnippetDiscussion creates a new discussion for a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-snippet-thread
		CreateSnippetDiscussion(pid any, snippet int64, opt *CreateSnippetDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// AddSnippetDiscussionNote adds a new note to a snippet discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-snippet-thread
		AddSnippetDiscussionNote(pid any, snippet int64, discussion string, opt *AddSnippetDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// UpdateSnippetDiscussionNote modifies an existing note in a snippet discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#modify-existing-snippet-thread-note
		UpdateSnippetDiscussionNote(pid any, snippet int64, discussion string, note int64, opt *UpdateSnippetDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// DeleteSnippetDiscussionNote deletes a note from a snippet discussion.
		//
		// GitLab API docs: https://docs.gitlab.com/api/discussions/#delete-a-snippet-thread-note
		DeleteSnippetDiscussionNote(pid any, snippet int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error)

		// ListGroupEpicDiscussions gets all discussions for a group epic.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#list-group-epic-discussion-items
		ListGroupEpicDiscussions(gid any, epic int64, opt *ListGroupEpicDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error)

		// GetEpicDiscussion returns a single discussion for a group epic.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-epic-discussion-item
		GetEpicDiscussion(gid any, epic int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// CreateEpicDiscussion creates a new discussion for a group epic.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-epic-thread
		CreateEpicDiscussion(gid any, epic int64, opt *CreateEpicDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// AddEpicDiscussionNote adds a new note to an epic discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-epic-thread
		AddEpicDiscussionNote(gid any, epic int64, discussion string, opt *AddEpicDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// UpdateEpicDiscussionNote modifies an existing note in an epic discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#modify-existing-epic-thread-note
		UpdateEpicDiscussionNote(gid any, epic int64, discussion string, note int64, opt *UpdateEpicDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// DeleteEpicDiscussionNote deletes a note from an epic discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#delete-an-epic-thread-note
		DeleteEpicDiscussionNote(gid any, epic int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error)

		// ListMergeRequestDiscussions gets all discussions for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#list-project-merge-request-discussion-items
		ListMergeRequestDiscussions(pid any, mergeRequest int64, opt *ListMergeRequestDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error)

		// GetMergeRequestDiscussion returns a single discussion for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-merge-request-discussion-item
		GetMergeRequestDiscussion(pid any, mergeRequest int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// CreateMergeRequestDiscussion creates a new discussion for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-merge-request-thread
		CreateMergeRequestDiscussion(pid any, mergeRequest int64, opt *CreateMergeRequestDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// ResolveMergeRequestDiscussion resolves or unresolves a merge request discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#resolve-a-merge-request-thread
		ResolveMergeRequestDiscussion(pid any, mergeRequest int64, discussion string, opt *ResolveMergeRequestDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// AddMergeRequestDiscussionNote adds a new note to a merge request discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-merge-request-thread
		AddMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, opt *AddMergeRequestDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// UpdateMergeRequestDiscussionNote modifies an existing note in a merge request discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#modify-an-existing-merge-request-thread-note
		UpdateMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, note int64, opt *UpdateMergeRequestDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// DeleteMergeRequestDiscussionNote deletes a note from a merge request discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#delete-a-merge-request-thread-note
		DeleteMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error)

		// ListCommitDiscussions gets all discussions for a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#list-project-commit-discussion-items
		ListCommitDiscussions(pid any, commit string, opt *ListCommitDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error)

		// GetCommitDiscussion returns a single discussion for a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-commit-discussion-item

		// GetCommitDiscussion returns a single discussion for a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#get-single-commit-discussion-item
		GetCommitDiscussion(pid any, commit string, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// CreateCommitDiscussion creates a new discussion for a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-commit-thread

		// CreateCommitDiscussion creates a new discussion for a commit.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#create-new-commit-thread
		CreateCommitDiscussion(pid any, commit string, opt *CreateCommitDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error)

		// AddCommitDiscussionNote adds a new note to a commit discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-commit-thread

		// AddCommitDiscussionNote adds a new note to a commit discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#add-note-to-existing-commit-thread
		AddCommitDiscussionNote(pid any, commit string, discussion string, opt *AddCommitDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// UpdateCommitDiscussionNote modifies an existing note in a commit discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#modify-an-existing-commit-thread-note
		UpdateCommitDiscussionNote(pid any, commit string, discussion string, note int64, opt *UpdateCommitDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error)

		// DeleteCommitDiscussionNote deletes a note from a commit discussion.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/discussions/#delete-a-commit-thread-note
		DeleteCommitDiscussionNote(pid any, commit string, discussion string, note int64, options ...RequestOptionFunc) (*Response, error)
	}

	// DiscussionsService handles communication with the discussions related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/discussions/
	DiscussionsService struct {
		client *Client
	}
)

var _ DiscussionsServiceInterface = (*DiscussionsService)(nil)

// Discussion represents a GitLab discussion.
//
// GitLab API docs: https://docs.gitlab.com/api/discussions/
type Discussion struct {
	ID             string  `json:"id"`
	IndividualNote bool    `json:"individual_note"`
	Notes          []*Note `json:"notes"`
}

func (d Discussion) String() string {
	return Stringify(d)
}

// ListIssueDiscussionsOptions represents the available ListIssueDiscussions()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#list-project-issue-discussion-items
type ListIssueDiscussionsOptions struct {
	ListOptions
}

func (s *DiscussionsService) ListIssueDiscussions(pid any, issue int64, opt *ListIssueDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error) {
	return do[[]*Discussion](s.client,
		withPath("projects/%s/issues/%d/discussions", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) GetIssueDiscussion(pid any, issue int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withPath("projects/%s/issues/%d/discussions/%s", ProjectID{pid}, issue, discussion),
		withRequestOpts(options...),
	)
}

// CreateIssueDiscussionOptions represents the available CreateIssueDiscussion()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#create-new-issue-thread
type CreateIssueDiscussionOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) CreateIssueDiscussion(pid any, issue int64, opt *CreateIssueDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/issues/%d/discussions", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddIssueDiscussionNoteOptions represents the available AddIssueDiscussionNote()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#add-note-to-existing-issue-thread
type AddIssueDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) AddIssueDiscussionNote(pid any, issue int64, discussion string, opt *AddIssueDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/issues/%d/discussions/%s/notes", ProjectID{pid}, issue, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateIssueDiscussionNoteOptions represents the available
// UpdateIssueDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#modify-existing-issue-thread-note
type UpdateIssueDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) UpdateIssueDiscussionNote(pid any, issue int64, discussion string, note int64, opt *UpdateIssueDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/issues/%d/discussions/%s/notes/%d", ProjectID{pid}, issue, discussion, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) DeleteIssueDiscussionNote(pid any, issue int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/issues/%d/discussions/%s/notes/%d", ProjectID{pid}, issue, discussion, note),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListSnippetDiscussionsOptions represents the available ListSnippetDiscussions()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#list-project-snippet-discussion-items
type ListSnippetDiscussionsOptions struct {
	ListOptions
}

func (s *DiscussionsService) ListSnippetDiscussions(pid any, snippet int64, opt *ListSnippetDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error) {
	return do[[]*Discussion](s.client,
		withPath("projects/%s/snippets/%d/discussions", ProjectID{pid}, snippet),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) GetSnippetDiscussion(pid any, snippet int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withPath("projects/%s/snippets/%d/discussions/%s", ProjectID{pid}, snippet, discussion),
		withRequestOpts(options...),
	)
}

// CreateSnippetDiscussionOptions represents the available
// CreateSnippetDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#create-new-snippet-thread
type CreateSnippetDiscussionOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) CreateSnippetDiscussion(pid any, snippet int64, opt *CreateSnippetDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/snippets/%d/discussions", ProjectID{pid}, snippet),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddSnippetDiscussionNoteOptions represents the available
// AddSnippetDiscussionNote() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#add-note-to-existing-snippet-thread
type AddSnippetDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) AddSnippetDiscussionNote(pid any, snippet int64, discussion string, opt *AddSnippetDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/snippets/%d/discussions/%s/notes", ProjectID{pid}, snippet, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateSnippetDiscussionNoteOptions represents the available
// UpdateSnippetDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#modify-existing-snippet-thread-note
type UpdateSnippetDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) UpdateSnippetDiscussionNote(pid any, snippet int64, discussion string, note int64, opt *UpdateSnippetDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/snippets/%d/discussions/%s/notes/%d", ProjectID{pid}, snippet, discussion, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) DeleteSnippetDiscussionNote(pid any, snippet int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/snippets/%d/discussions/%s/notes/%d", ProjectID{pid}, snippet, discussion, note),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListGroupEpicDiscussionsOptions represents the available
// ListEpicDiscussions() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#list-group-epic-discussion-items
type ListGroupEpicDiscussionsOptions struct {
	ListOptions
}

func (s *DiscussionsService) ListGroupEpicDiscussions(gid any, epic int64, opt *ListGroupEpicDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error) {
	return do[[]*Discussion](s.client,
		withPath("groups/%s/epics/%d/discussions", GroupID{gid}, epic),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) GetEpicDiscussion(gid any, epic int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withPath("groups/%s/epics/%d/discussions/%s", GroupID{gid}, epic, discussion),
		withRequestOpts(options...),
	)
}

// CreateEpicDiscussionOptions represents the available CreateEpicDiscussion()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#create-new-epic-thread
type CreateEpicDiscussionOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) CreateEpicDiscussion(gid any, epic int64, opt *CreateEpicDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/epics/%d/discussions", GroupID{gid}, epic),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddEpicDiscussionNoteOptions represents the available
// AddEpicDiscussionNote() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#add-note-to-existing-epic-thread
type AddEpicDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) AddEpicDiscussionNote(gid any, epic int64, discussion string, opt *AddEpicDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/epics/%d/discussions/%s/notes", GroupID{gid}, epic, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateEpicDiscussionNoteOptions represents the available UpdateEpicDiscussion()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#modify-existing-epic-thread-note
type UpdateEpicDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) UpdateEpicDiscussionNote(gid any, epic int64, discussion string, note int64, opt *UpdateEpicDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/epics/%d/discussions/%s/notes/%d", GroupID{gid}, epic, discussion, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) DeleteEpicDiscussionNote(gid any, epic int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/epics/%d/discussions/%s/notes/%d", GroupID{gid}, epic, discussion, note),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListMergeRequestDiscussionsOptions represents the available
// ListMergeRequestDiscussions() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#list-project-merge-request-discussion-items
type ListMergeRequestDiscussionsOptions struct {
	ListOptions
}

func (s *DiscussionsService) ListMergeRequestDiscussions(pid any, mergeRequest int64, opt *ListMergeRequestDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error) {
	return do[[]*Discussion](s.client,
		withPath("projects/%s/merge_requests/%d/discussions", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) GetMergeRequestDiscussion(pid any, mergeRequest int64, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withPath("projects/%s/merge_requests/%d/discussions/%s", ProjectID{pid}, mergeRequest, discussion),
		withRequestOpts(options...),
	)
}

// CreateMergeRequestDiscussionOptions represents the available
// CreateMergeRequestDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#create-new-merge-request-thread
type CreateMergeRequestDiscussionOptions struct {
	Body      *string          `url:"body,omitempty" json:"body,omitempty"`
	CommitID  *string          `url:"commit_id,omitempty" json:"commit_id,omitempty"`
	CreatedAt *time.Time       `url:"created_at,omitempty" json:"created_at,omitempty"`
	Position  *PositionOptions `url:"position,omitempty" json:"position,omitempty"`
}

// PositionOptions represents the position option of a discussion.
type PositionOptions struct {
	BaseSHA      *string           `url:"base_sha,omitempty" json:"base_sha,omitempty"`
	HeadSHA      *string           `url:"head_sha,omitempty" json:"head_sha,omitempty"`
	StartSHA     *string           `url:"start_sha,omitempty" json:"start_sha,omitempty"`
	NewPath      *string           `url:"new_path,omitempty" json:"new_path,omitempty"`
	OldPath      *string           `url:"old_path,omitempty" json:"old_path,omitempty"`
	PositionType *string           `url:"position_type,omitempty" json:"position_type"`
	NewLine      *int64            `url:"new_line,omitempty" json:"new_line,omitempty"`
	OldLine      *int64            `url:"old_line,omitempty" json:"old_line,omitempty"`
	LineRange    *LineRangeOptions `url:"line_range,omitempty" json:"line_range,omitempty"`
	Width        *int64            `url:"width,omitempty" json:"width,omitempty"`
	Height       *int64            `url:"height,omitempty" json:"height,omitempty"`
	X            *float64          `url:"x,omitempty" json:"x,omitempty"`
	Y            *float64          `url:"y,omitempty" json:"y,omitempty"`
}

// LineRangeOptions represents the line range option of a discussion.
type LineRangeOptions struct {
	Start *LinePositionOptions `url:"start,omitempty" json:"start,omitempty"`
	End   *LinePositionOptions `url:"end,omitempty" json:"end,omitempty"`
}

// LinePositionOptions represents the line position option of a discussion.
type LinePositionOptions struct {
	LineCode *string `url:"line_code,omitempty" json:"line_code,omitempty"`
	Type     *string `url:"type,omitempty" json:"type,omitempty"`
	OldLine  *int64  `url:"old_line,omitempty" json:"old_line,omitempty"`
	NewLine  *int64  `url:"new_line,omitempty" json:"new_line,omitempty"`
}

func (s *DiscussionsService) CreateMergeRequestDiscussion(pid any, mergeRequest int64, opt *CreateMergeRequestDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/discussions", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ResolveMergeRequestDiscussionOptions represents the available
// ResolveMergeRequestDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#resolve-a-merge-request-thread
type ResolveMergeRequestDiscussionOptions struct {
	Resolved *bool `url:"resolved,omitempty" json:"resolved,omitempty"`
}

func (s *DiscussionsService) ResolveMergeRequestDiscussion(pid any, mergeRequest int64, discussion string, opt *ResolveMergeRequestDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/discussions/%s", ProjectID{pid}, mergeRequest, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddMergeRequestDiscussionNoteOptions represents the available
// AddMergeRequestDiscussionNote() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#add-note-to-existing-merge-request-thread
type AddMergeRequestDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) AddMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, opt *AddMergeRequestDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/discussions/%s/notes", ProjectID{pid}, mergeRequest, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateMergeRequestDiscussionNoteOptions represents the available
// UpdateMergeRequestDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#modify-an-existing-merge-request-thread-note
type UpdateMergeRequestDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
	Resolved  *bool      `url:"resolved,omitempty" json:"resolved,omitempty"`
}

func (s *DiscussionsService) UpdateMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, note int64, opt *UpdateMergeRequestDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/discussions/%s/notes/%d", ProjectID{pid}, mergeRequest, discussion, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) DeleteMergeRequestDiscussionNote(pid any, mergeRequest int64, discussion string, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/merge_requests/%d/discussions/%s/notes/%d", ProjectID{pid}, mergeRequest, discussion, note),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListCommitDiscussionsOptions represents the available
// ListCommitDiscussions() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#list-project-commit-discussion-items
type ListCommitDiscussionsOptions struct {
	ListOptions
}

func (s *DiscussionsService) ListCommitDiscussions(pid any, commit string, opt *ListCommitDiscussionsOptions, options ...RequestOptionFunc) ([]*Discussion, *Response, error) {
	return do[[]*Discussion](s.client,
		withPath("projects/%s/repository/commits/%s/discussions", ProjectID{pid}, commit),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) GetCommitDiscussion(pid any, commit string, discussion string, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withPath("projects/%s/repository/commits/%s/discussions/%s", ProjectID{pid}, commit, discussion),
		withRequestOpts(options...),
	)
}

// CreateCommitDiscussionOptions represents the available
// CreateCommitDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#create-new-commit-thread
type CreateCommitDiscussionOptions struct {
	Body      *string       `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time    `url:"created_at,omitempty" json:"created_at,omitempty"`
	Position  *NotePosition `url:"position,omitempty" json:"position,omitempty"`
}

func (s *DiscussionsService) CreateCommitDiscussion(pid any, commit string, opt *CreateCommitDiscussionOptions, options ...RequestOptionFunc) (*Discussion, *Response, error) {
	return do[*Discussion](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits/%s/discussions", ProjectID{pid}, commit),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddCommitDiscussionNoteOptions represents the available
// AddCommitDiscussionNote() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#add-note-to-existing-commit-thread
type AddCommitDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) AddCommitDiscussionNote(pid any, commit string, discussion string, opt *AddCommitDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/commits/%s/discussions/%s/notes", ProjectID{pid}, commit, discussion),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateCommitDiscussionNoteOptions represents the available
// UpdateCommitDiscussion() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/discussions/#modify-an-existing-commit-thread-note
type UpdateCommitDiscussionNoteOptions struct {
	Body      *string    `url:"body,omitempty" json:"body,omitempty"`
	CreatedAt *time.Time `url:"created_at,omitempty" json:"created_at,omitempty"`
}

func (s *DiscussionsService) UpdateCommitDiscussionNote(pid any, commit string, discussion string, note int64, opt *UpdateCommitDiscussionNoteOptions, options ...RequestOptionFunc) (*Note, *Response, error) {
	return do[*Note](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/repository/commits/%s/discussions/%s/notes/%d", ProjectID{pid}, commit, discussion, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DiscussionsService) DeleteCommitDiscussionNote(pid any, commit string, discussion string, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/repository/commits/%s/discussions/%s/notes/%d", ProjectID{pid}, commit, discussion, note),
		withRequestOpts(options...),
	)
	return resp, err
}
