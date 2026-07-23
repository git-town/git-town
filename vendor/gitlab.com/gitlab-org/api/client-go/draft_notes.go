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
)

type (
	// DraftNotesServiceInterface defines all the API methods for the DraftNotesService
	DraftNotesServiceInterface interface {
		// ListDraftNotes gets a list of all draft notes for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#list-all-merge-request-draft-notes
		ListDraftNotes(pid any, mergeRequest int64, opt *ListDraftNotesOptions, options ...RequestOptionFunc) ([]*DraftNote, *Response, error)

		// GetDraftNote gets a single draft note for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#get-a-single-draft-note
		GetDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*DraftNote, *Response, error)

		// CreateDraftNote creates a draft note for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#create-a-draft-note
		CreateDraftNote(pid any, mergeRequest int64, opt *CreateDraftNoteOptions, options ...RequestOptionFunc) (*DraftNote, *Response, error)

		// UpdateDraftNote updates a draft note for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#update-a-draft-note
		UpdateDraftNote(pid any, mergeRequest int64, note int64, opt *UpdateDraftNoteOptions, options ...RequestOptionFunc) (*DraftNote, *Response, error)

		// DeleteDraftNote deletes a single draft note for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#delete-a-draft-note
		DeleteDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*Response, error)

		// PublishDraftNote publishes a single draft note for a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#publish-a-draft-note
		PublishDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*Response, error)

		// PublishAllDraftNotes publishes all draft notes for a merge request that belong to the user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/draft_notes/#publish-a-draft-note
		PublishAllDraftNotes(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Response, error)
	}

	// DraftNotesService handles communication with the draft notes related methods
	// of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/draft_notes/#list-all-merge-request-draft-notes
	DraftNotesService struct {
		client *Client
	}
)

var _ DraftNotesServiceInterface = (*DraftNotesService)(nil)

type DraftNote struct {
	ID                int64         `json:"id"`
	AuthorID          int64         `json:"author_id"`
	MergeRequestID    int64         `json:"merge_request_id"`
	ResolveDiscussion bool          `json:"resolve_discussion"`
	DiscussionID      string        `json:"discussion_id"`
	Note              string        `json:"note"`
	CommitID          string        `json:"commit_id"`
	LineCode          string        `json:"line_code"`
	Position          *NotePosition `json:"position"`
}

// ListDraftNotesOptions represents the available ListDraftNotes()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/draft_notes/#list-all-merge-request-draft-notes
type ListDraftNotesOptions struct {
	ListOptions
	OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
}

func (s *DraftNotesService) ListDraftNotes(pid any, mergeRequest int64, opt *ListDraftNotesOptions, options ...RequestOptionFunc) ([]*DraftNote, *Response, error) {
	return do[[]*DraftNote](s.client,
		withPath("projects/%s/merge_requests/%d/draft_notes", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DraftNotesService) GetDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*DraftNote, *Response, error) {
	return do[*DraftNote](s.client,
		withPath("projects/%s/merge_requests/%d/draft_notes/%d", ProjectID{pid}, mergeRequest, note),
		withRequestOpts(options...),
	)
}

// CreateDraftNoteOptions represents the available CreateDraftNote()
// options.
//
// GitLab API docs:
// GitLab API docs:
// https://docs.gitlab.com/api/draft_notes/#create-a-draft-note
type CreateDraftNoteOptions struct {
	Note                  *string          `url:"note" json:"note"`
	CommitID              *string          `url:"commit_id,omitempty" json:"commit_id,omitempty"`
	InReplyToDiscussionID *string          `url:"in_reply_to_discussion_id,omitempty" json:"in_reply_to_discussion_id,omitempty"`
	ResolveDiscussion     *bool            `url:"resolve_discussion,omitempty" json:"resolve_discussion,omitempty"`
	Position              *PositionOptions `url:"position,omitempty" json:"position,omitempty"`
}

func (s *DraftNotesService) CreateDraftNote(pid any, mergeRequest int64, opt *CreateDraftNoteOptions, options ...RequestOptionFunc) (*DraftNote, *Response, error) {
	return do[*DraftNote](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/draft_notes", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateDraftNoteOptions represents the available UpdateDraftNote()
// options.
//
// GitLab API docs:
// GitLab API docs:
// https://docs.gitlab.com/api/draft_notes/#create-a-draft-note
type UpdateDraftNoteOptions struct {
	Note     *string          `url:"note,omitempty" json:"note,omitempty"`
	Position *PositionOptions `url:"position,omitempty" json:"position,omitempty"`
}

func (s *DraftNotesService) UpdateDraftNote(pid any, mergeRequest int64, note int64, opt *UpdateDraftNoteOptions, options ...RequestOptionFunc) (*DraftNote, *Response, error) {
	return do[*DraftNote](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/draft_notes/%d", ProjectID{pid}, mergeRequest, note),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DraftNotesService) DeleteDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/merge_requests/%d/draft_notes/%d", ProjectID{pid}, mergeRequest, note),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *DraftNotesService) PublishDraftNote(pid any, mergeRequest int64, note int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/merge_requests/%d/draft_notes/%d/publish", ProjectID{pid}, mergeRequest, note),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *DraftNotesService) PublishAllDraftNotes(pid any, mergeRequest int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/draft_notes/bulk_publish", ProjectID{pid}, mergeRequest),
		withRequestOpts(options...),
	)
	return resp, err
}
