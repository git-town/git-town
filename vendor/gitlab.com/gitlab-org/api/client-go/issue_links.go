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
	"net/http"
	"time"
)

type (
	IssueLinksServiceInterface interface {
		ListIssueRelations(pid any, issue int64, options ...RequestOptionFunc) ([]*IssueRelation, *Response, error)
		GetIssueLink(pid any, issue, issueLink int64, options ...RequestOptionFunc) (*IssueLink, *Response, error)
		CreateIssueLink(pid any, issue int64, opt *CreateIssueLinkOptions, options ...RequestOptionFunc) (*IssueLink, *Response, error)
		DeleteIssueLink(pid any, issue, issueLink int64, options ...RequestOptionFunc) (*IssueLink, *Response, error)
	}

	// IssueLinksService handles communication with the issue relations related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/issue_links/
	IssueLinksService struct {
		client *Client
	}
)

var _ IssueLinksServiceInterface = (*IssueLinksService)(nil)

// IssueLink represents a two-way relation between two issues.
//
// GitLab API docs: https://docs.gitlab.com/api/issue_links/
type IssueLink struct {
	ID          int64  `json:"id"`
	SourceIssue *Issue `json:"source_issue"`
	TargetIssue *Issue `json:"target_issue"`
	LinkType    string `json:"link_type"`
}

// IssueRelation gets a relation between two issues.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issue_links/#list-issue-relations
type IssueRelation struct {
	ID             int64            `json:"id"`
	IID            int64            `json:"iid"`
	State          string           `json:"state"`
	Description    string           `json:"description"`
	Confidential   bool             `json:"confidential"`
	Author         *IssueAuthor     `json:"author"`
	Milestone      *Milestone       `json:"milestone"`
	ProjectID      int64            `json:"project_id"`
	Assignees      []*IssueAssignee `json:"assignees"`
	Assignee       *IssueAssignee   `json:"assignee"`
	UpdatedAt      *time.Time       `json:"updated_at"`
	Title          string           `json:"title"`
	CreatedAt      *time.Time       `json:"created_at"`
	Labels         Labels           `json:"labels"`
	DueDate        *ISOTime         `json:"due_date"`
	WebURL         string           `json:"web_url"`
	References     *IssueReferences `json:"references"`
	Weight         int64            `json:"weight"`
	UserNotesCount int64            `json:"user_notes_count"`
	IssueLinkID    int64            `json:"issue_link_id"`
	LinkType       string           `json:"link_type"`
	LinkCreatedAt  *time.Time       `json:"link_created_at"`
	LinkUpdatedAt  *time.Time       `json:"link_updated_at"`
}

// ListIssueRelations gets a list of related issues of a given issue,
// sorted by the relationship creation datetime (ascending).
//
// Issues will be filtered according to the user authorizations.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issue_links/#list-issue-relations
func (s *IssueLinksService) ListIssueRelations(pid any, issue int64, options ...RequestOptionFunc) ([]*IssueRelation, *Response, error) {
	// Use explicit format string for the path
	return do[[]*IssueRelation](s.client,
		withPath("projects/%s/issues/%d/links", ProjectID{pid}, issue),
		withRequestOpts(options...),
	)
}

// GetIssueLink gets a specific issue link.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issue_links/#get-an-issue-link
func (s *IssueLinksService) GetIssueLink(pid any, issue, issueLink int64, options ...RequestOptionFunc) (*IssueLink, *Response, error) {
	// Use explicit format string for the path
	return do[*IssueLink](s.client,
		withPath("projects/%s/issues/%d/links/%d", ProjectID{pid}, issue, issueLink),
		withRequestOpts(options...),
	)
}

// CreateIssueLinkOptions represents the available CreateIssueLink() options.
//
// GitLab API docs: https://docs.gitlab.com/api/issue_links/#create-an-issue-link
type CreateIssueLinkOptions struct {
	TargetProjectID *string `json:"target_project_id"`
	TargetIssueIID  *string `json:"target_issue_iid"`
	LinkType        *string `json:"link_type"`
}

// CreateIssueLink creates a two-way relation between two issues.
// User must be allowed to update both issues in order to succeed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issue_links/#create-an-issue-link
func (s *IssueLinksService) CreateIssueLink(pid any, issue int64, opt *CreateIssueLinkOptions, options ...RequestOptionFunc) (*IssueLink, *Response, error) {
	return do[*IssueLink](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/issues/%d/links", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteIssueLink deletes an issue link, thus removes the two-way relationship.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issue_links/#delete-an-issue-link
func (s *IssueLinksService) DeleteIssueLink(pid any, issue, issueLink int64, options ...RequestOptionFunc) (*IssueLink, *Response, error) {
	return do[*IssueLink](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/issues/%d/links/%d", ProjectID{pid}, issue, issueLink),
		withRequestOpts(options...),
	)
}
