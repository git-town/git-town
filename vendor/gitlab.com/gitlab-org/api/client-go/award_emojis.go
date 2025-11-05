//
// Copyright 2021, Arkbriar
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
	"fmt"
	"net/http"
	"time"
)

type (
	AwardEmojiServiceInterface interface {
		// ListMergeRequestAwardEmoji gets a list of all award emoji on the merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-an-awardables-emoji-reactions
		ListMergeRequestAwardEmoji(pid any, mergeRequestIID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// ListIssueAwardEmoji gets a list of all award emoji on the issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-an-awardables-emoji-reactions
		ListIssueAwardEmoji(pid any, issueIID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// ListSnippetAwardEmoji gets a list of all award emoji on the snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-an-awardables-emoji-reactions
		ListSnippetAwardEmoji(pid any, snippetID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// GetMergeRequestAwardEmoji get an award emoji from merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-single-emoji-reaction
		GetMergeRequestAwardEmoji(pid any, mergeRequestIID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// GetIssueAwardEmoji get an award emoji from issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-single-emoji-reaction
		GetIssueAwardEmoji(pid any, issueIID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// GetSnippetAwardEmoji get an award emoji from snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-single-emoji-reaction
		GetSnippetAwardEmoji(pid any, snippetID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateMergeRequestAwardEmoji get an award emoji from merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction
		CreateMergeRequestAwardEmoji(pid any, mergeRequestIID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateIssueAwardEmoji get an award emoji from issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction
		CreateIssueAwardEmoji(pid any, issueIID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateSnippetAwardEmoji get an award emoji from snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction
		CreateSnippetAwardEmoji(pid any, snippetID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// DeleteIssueAwardEmoji delete award emoji on an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction
		DeleteIssueAwardEmoji(pid any, issueIID, awardID int, options ...RequestOptionFunc) (*Response, error)

		// DeleteMergeRequestAwardEmoji delete award emoji on a merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction
		DeleteMergeRequestAwardEmoji(pid any, mergeRequestIID, awardID int, options ...RequestOptionFunc) (*Response, error)

		// DeleteSnippetAwardEmoji delete award emoji on a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction
		DeleteSnippetAwardEmoji(pid any, snippetID, awardID int, options ...RequestOptionFunc) (*Response, error)

		// ListIssuesAwardEmojiOnNote gets a list of all award emoji on a note from the
		// issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-a-comments-emoji-reactions
		ListIssuesAwardEmojiOnNote(pid any, issueID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// ListMergeRequestAwardEmojiOnNote gets a list of all award emoji on a note
		// from the merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-a-comments-emoji-reactions
		ListMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// ListSnippetAwardEmojiOnNote gets a list of all award emoji on a note from the
		// snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#list-a-comments-emoji-reactions
		ListSnippetAwardEmojiOnNote(pid any, snippetIID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error)

		// GetIssuesAwardEmojiOnNote gets an award emoji on a note from an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-an-emoji-reaction-for-a-comment
		GetIssuesAwardEmojiOnNote(pid any, issueID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// GetMergeRequestAwardEmojiOnNote gets an award emoji on a note from a
		// merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-an-emoji-reaction-for-a-comment
		GetMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// GetSnippetAwardEmojiOnNote gets an award emoji on a note from a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#get-an-emoji-reaction-for-a-comment
		GetSnippetAwardEmojiOnNote(pid any, snippetIID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateIssuesAwardEmojiOnNote gets an award emoji on a note from an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction-to-a-comment
		CreateIssuesAwardEmojiOnNote(pid any, issueID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateMergeRequestAwardEmojiOnNote gets an award emoji on a note from a
		// merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction-to-a-comment
		CreateMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// CreateSnippetAwardEmojiOnNote gets an award emoji on a note from a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction-to-a-comment
		CreateSnippetAwardEmojiOnNote(pid any, snippetIID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error)

		// DeleteIssuesAwardEmojiOnNote deletes an award emoji on a note from an issue.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction-from-a-comment
		DeleteIssuesAwardEmojiOnNote(pid any, issueID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error)

		// DeleteMergeRequestAwardEmojiOnNote deletes an award emoji on a note from a
		// merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction-from-a-comment
		DeleteMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error)

		// DeleteSnippetAwardEmojiOnNote deletes an award emoji on a note from a snippet.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction-from-a-comment
		DeleteSnippetAwardEmojiOnNote(pid any, snippetIID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error)
	}

	// AwardEmojiService handles communication with the emoji awards related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/emoji_reactions/
	AwardEmojiService struct {
		client *Client
	}
)

var _ AwardEmojiServiceInterface = (*AwardEmojiService)(nil)

// AwardEmoji represents a GitLab Award Emoji.
//
// GitLab API docs: https://docs.gitlab.com/api/emoji_reactions/
type AwardEmoji struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	User struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		ID        int    `json:"id"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"user"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	AwardableID   int        `json:"awardable_id"`
	AwardableType string     `json:"awardable_type"`
}

const (
	awardMergeRequest = "merge_requests"
	awardIssue        = "issues"
	awardSnippets     = "snippets"
)

// ListAwardEmojiOptions represents the available options for listing emoji
// for each resource
//
// GitLab API docs:
// https://docs.gitlab.com/api/emoji_reactions/
type ListAwardEmojiOptions ListOptions

func (s *AwardEmojiService) ListMergeRequestAwardEmoji(pid any, mergeRequestIID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmoji(pid, awardMergeRequest, mergeRequestIID, opt, options...)
}

func (s *AwardEmojiService) ListIssueAwardEmoji(pid any, issueIID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmoji(pid, awardIssue, issueIID, opt, options...)
}

func (s *AwardEmojiService) ListSnippetAwardEmoji(pid any, snippetID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmoji(pid, awardSnippets, snippetID, opt, options...)
}

func (s *AwardEmojiService) listAwardEmoji(pid any, resource string, resourceID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/award_emoji",
		PathEscape(project),
		resource,
		resourceID,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var as []*AwardEmoji
	resp, err := s.client.Do(req, &as)
	if err != nil {
		return nil, resp, err
	}

	return as, resp, nil
}

func (s *AwardEmojiService) GetMergeRequestAwardEmoji(pid any, mergeRequestIID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getAwardEmoji(pid, awardMergeRequest, mergeRequestIID, awardID, options...)
}

func (s *AwardEmojiService) GetIssueAwardEmoji(pid any, issueIID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getAwardEmoji(pid, awardIssue, issueIID, awardID, options...)
}

func (s *AwardEmojiService) GetSnippetAwardEmoji(pid any, snippetID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getAwardEmoji(pid, awardSnippets, snippetID, awardID, options...)
}

func (s *AwardEmojiService) getAwardEmoji(pid any, resource string, resourceID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/award_emoji/%d",
		PathEscape(project),
		resource,
		resourceID,
		awardID,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(AwardEmoji)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// CreateAwardEmojiOptions represents the available options for awarding emoji
// for a resource
//
// GitLab API docs:
// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction
type CreateAwardEmojiOptions struct {
	Name string `json:"name"`
}

func (s *AwardEmojiService) CreateMergeRequestAwardEmoji(pid any, mergeRequestIID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmoji(pid, awardMergeRequest, mergeRequestIID, opt, options...)
}

func (s *AwardEmojiService) CreateIssueAwardEmoji(pid any, issueIID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmoji(pid, awardIssue, issueIID, opt, options...)
}

func (s *AwardEmojiService) CreateSnippetAwardEmoji(pid any, snippetID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmoji(pid, awardSnippets, snippetID, opt, options...)
}

func (s *AwardEmojiService) createAwardEmoji(pid any, resource string, resourceID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/award_emoji",
		PathEscape(project),
		resource,
		resourceID,
	)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(AwardEmoji)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

func (s *AwardEmojiService) DeleteIssueAwardEmoji(pid any, issueIID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmoji(pid, awardIssue, issueIID, awardID, options...)
}

func (s *AwardEmojiService) DeleteMergeRequestAwardEmoji(pid any, mergeRequestIID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmoji(pid, awardMergeRequest, mergeRequestIID, awardID, options...)
}

func (s *AwardEmojiService) DeleteSnippetAwardEmoji(pid any, snippetID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmoji(pid, awardSnippets, snippetID, awardID, options...)
}

// DeleteAwardEmoji Delete an award emoji on the specified resource.
//
// GitLab API docs:
// https://docs.gitlab.com/api/emoji_reactions/#delete-an-emoji-reaction
func (s *AwardEmojiService) deleteAwardEmoji(pid any, resource string, resourceID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/award_emoji/%d", PathEscape(project), resource,
		resourceID, awardID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req, nil)
}

func (s *AwardEmojiService) ListIssuesAwardEmojiOnNote(pid any, issueID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmojiOnNote(pid, awardIssue, issueID, noteID, opt, options...)
}

func (s *AwardEmojiService) ListMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmojiOnNote(pid, awardMergeRequest, mergeRequestIID, noteID, opt, options...)
}

func (s *AwardEmojiService) ListSnippetAwardEmojiOnNote(pid any, snippetIID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	return s.listAwardEmojiOnNote(pid, awardSnippets, snippetIID, noteID, opt, options...)
}

func (s *AwardEmojiService) listAwardEmojiOnNote(pid any, resources string, resourceID, noteID int, opt *ListAwardEmojiOptions, options ...RequestOptionFunc) ([]*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/notes/%d/award_emoji", PathEscape(project), resources,
		resourceID, noteID)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var as []*AwardEmoji
	resp, err := s.client.Do(req, &as)
	if err != nil {
		return nil, resp, err
	}

	return as, resp, nil
}

func (s *AwardEmojiService) GetIssuesAwardEmojiOnNote(pid any, issueID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getSingleNoteAwardEmoji(pid, awardIssue, issueID, noteID, awardID, options...)
}

func (s *AwardEmojiService) GetMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getSingleNoteAwardEmoji(pid, awardMergeRequest, mergeRequestIID, noteID, awardID,
		options...)
}

func (s *AwardEmojiService) GetSnippetAwardEmojiOnNote(pid any, snippetIID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.getSingleNoteAwardEmoji(pid, awardSnippets, snippetIID, noteID, awardID, options...)
}

func (s *AwardEmojiService) getSingleNoteAwardEmoji(pid any, resource string, resourceID, noteID, awardID int, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/notes/%d/award_emoji/%d",
		PathEscape(project),
		resource,
		resourceID,
		noteID,
		awardID,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(AwardEmoji)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

func (s *AwardEmojiService) CreateIssuesAwardEmojiOnNote(pid any, issueID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmojiOnNote(pid, awardIssue, issueID, noteID, opt, options...)
}

func (s *AwardEmojiService) CreateMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmojiOnNote(pid, awardMergeRequest, mergeRequestIID, noteID, opt, options...)
}

func (s *AwardEmojiService) CreateSnippetAwardEmojiOnNote(pid any, snippetIID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	return s.createAwardEmojiOnNote(pid, awardSnippets, snippetIID, noteID, opt, options...)
}

// CreateAwardEmojiOnNote award emoji on a note.
//
// GitLab API docs:
// https://docs.gitlab.com/api/emoji_reactions/#add-a-new-emoji-reaction-to-a-comment
func (s *AwardEmojiService) createAwardEmojiOnNote(pid any, resource string, resourceID, noteID int, opt *CreateAwardEmojiOptions, options ...RequestOptionFunc) (*AwardEmoji, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/notes/%d/award_emoji",
		PathEscape(project),
		resource,
		resourceID,
		noteID,
	)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(AwardEmoji)
	resp, err := s.client.Do(req, &a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

func (s *AwardEmojiService) DeleteIssuesAwardEmojiOnNote(pid any, issueID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmojiOnNote(pid, awardIssue, issueID, noteID, awardID, options...)
}

func (s *AwardEmojiService) DeleteMergeRequestAwardEmojiOnNote(pid any, mergeRequestIID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmojiOnNote(pid, awardMergeRequest, mergeRequestIID, noteID, awardID,
		options...)
}

func (s *AwardEmojiService) DeleteSnippetAwardEmojiOnNote(pid any, snippetIID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	return s.deleteAwardEmojiOnNote(pid, awardSnippets, snippetIID, noteID, awardID, options...)
}

func (s *AwardEmojiService) deleteAwardEmojiOnNote(pid any, resource string, resourceID, noteID, awardID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/notes/%d/award_emoji/%d",
		PathEscape(project),
		resource,
		resourceID,
		noteID,
		awardID,
	)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
