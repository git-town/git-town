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
	IssueBoardsServiceInterface interface {
		// CreateIssueBoard creates a new issue board.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#create-an-issue-board
		// CreateIssueBoard creates a new issue board.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#create-an-issue-board
		CreateIssueBoard(pid any, opt *CreateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error)

		// UpdateIssueBoard update an issue board.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#update-an-issue-board
		UpdateIssueBoard(pid any, board int64, opt *UpdateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error)

		// DeleteIssueBoard deletes an issue board.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#delete-an-issue-board
		DeleteIssueBoard(pid any, board int64, options ...RequestOptionFunc) (*Response, error)

		// ListIssueBoards gets a list of all issue boards in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#list-project-issue-boards
		ListIssueBoards(pid any, opt *ListIssueBoardsOptions, options ...RequestOptionFunc) ([]*IssueBoard, *Response, error)

		// GetIssueBoard gets a single issue board of a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#show-a-single-issue-board
		GetIssueBoard(pid any, board int64, options ...RequestOptionFunc) (*IssueBoard, *Response, error)

		// GetIssueBoardLists gets a list of the issue board's lists. Does not include
		// backlog and closed lists.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#list-board-lists-in-a-project-issue-board
		GetIssueBoardLists(pid any, board int64, opt *GetIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error)

		// GetIssueBoardList gets a single issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#show-a-single-board-list
		GetIssueBoardList(pid any, board, list int64, options ...RequestOptionFunc) (*BoardList, *Response, error)

		// CreateIssueBoardList creates a new issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#create-a-board-list
		CreateIssueBoardList(pid any, board int64, opt *CreateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error)

		// UpdateIssueBoardList updates the position of an existing issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#reorder-a-list-in-a-board
		UpdateIssueBoardList(pid any, board, list int64, opt *UpdateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error)

		// DeleteIssueBoardList soft deletes an issue board list. Only for admins and
		// project owners.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/boards/#delete-a-board-list-from-a-board
		DeleteIssueBoardList(pid any, board, list int64, options ...RequestOptionFunc) (*Response, error)
	}

	// IssueBoardsService handles communication with the issue board related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/boards/
	IssueBoardsService struct {
		client *Client
	}
)

var _ IssueBoardsServiceInterface = (*IssueBoardsService)(nil)

// IssueBoard represents a GitLab issue board.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/
type IssueBoard struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Project         *Project        `json:"project"`
	Milestone       *Milestone      `json:"milestone"`
	Assignee        *BasicUser      `json:"assignee"`
	Lists           []*BoardList    `json:"lists"`
	Weight          int64           `json:"weight"`
	Labels          []*LabelDetails `json:"labels"`
	HideBacklogList bool            `json:"hide_backlog_list"`
	HideClosedList  bool            `json:"hide_closed_list"`
}

func (b IssueBoard) String() string {
	return Stringify(b)
}

// BoardList represents a GitLab board list.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/
type BoardList struct {
	ID             int64              `json:"id"`
	Assignee       *BoardListAssignee `json:"assignee"`
	Iteration      *ProjectIteration  `json:"iteration"`
	Label          *Label             `json:"label"`
	MaxIssueCount  int64              `json:"max_issue_count"`
	MaxIssueWeight int64              `json:"max_issue_weight"`
	Milestone      *Milestone         `json:"milestone"`
	Position       int64              `json:"position"`
}

func (b BoardList) String() string {
	return Stringify(b)
}

// BoardListAssignee represents a GitLab board list assignee.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/
type BoardListAssignee struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (a BoardListAssignee) String() string {
	return Stringify(a)
}

// CreateIssueBoardOptions represents the available CreateIssueBoard() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-an-issue-board
type CreateIssueBoardOptions struct {
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

func (s *IssueBoardsService) CreateIssueBoard(pid any, opt *CreateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	return do[*IssueBoard](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/boards", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateIssueBoardOptions represents the available UpdateIssueBoard() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#update-an-issue-board
type UpdateIssueBoardOptions struct {
	Name            *string       `url:"name,omitempty" json:"name,omitempty"`
	AssigneeID      *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MilestoneID     *int64        `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	Labels          *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	Weight          *int64        `url:"weight,omitempty" json:"weight,omitempty"`
	HideBacklogList *bool         `url:"hide_backlog_list,omitempty" json:"hide_backlog_list,omitempty"`
	HideClosedList  *bool         `url:"hide_closed_list,omitempty" json:"hide_closed_list,omitempty"`
}

func (s *IssueBoardsService) UpdateIssueBoard(pid any, board int64, opt *UpdateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	return do[*IssueBoard](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/boards/%d", ProjectID{pid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IssueBoardsService) DeleteIssueBoard(pid any, board int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/boards/%d", ProjectID{pid}, board),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListIssueBoardsOptions represents the available ListIssueBoards() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-project-issue-boards
type ListIssueBoardsOptions struct {
	ListOptions
}

func (s *IssueBoardsService) ListIssueBoards(pid any, opt *ListIssueBoardsOptions, options ...RequestOptionFunc) ([]*IssueBoard, *Response, error) {
	return do[[]*IssueBoard](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/boards", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IssueBoardsService) GetIssueBoard(pid any, board int64, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	return do[*IssueBoard](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/boards/%d", ProjectID{pid}, board),
		withRequestOpts(options...),
	)
}

// GetIssueBoardListsOptions represents the available GetIssueBoardLists() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-board-lists-in-a-project-issue-board
type GetIssueBoardListsOptions struct {
	ListOptions
}

func (s *IssueBoardsService) GetIssueBoardLists(pid any, board int64, opt *GetIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error) {
	return do[[]*BoardList](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/boards/%d/lists", ProjectID{pid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IssueBoardsService) GetIssueBoardList(pid any, board, list int64, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	return do[*BoardList](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/boards/%d/lists/%d", ProjectID{pid}, board, list),
		withRequestOpts(options...),
	)
}

// CreateIssueBoardListOptions represents the available CreateIssueBoardList()
// options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-a-board-list
type CreateIssueBoardListOptions struct {
	LabelID     *int64 `url:"label_id,omitempty" json:"label_id,omitempty"`
	AssigneeID  *int64 `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MilestoneID *int64 `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	IterationID *int64 `url:"iteration_id,omitempty" json:"iteration_id,omitempty"`
}

func (s *IssueBoardsService) CreateIssueBoardList(pid any, board int64, opt *CreateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	return do[*BoardList](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/boards/%d/lists", ProjectID{pid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateIssueBoardListOptions represents the available UpdateIssueBoardList()
// options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#reorder-a-list-in-a-board
type UpdateIssueBoardListOptions struct {
	Position *int64 `url:"position" json:"position"`
}

func (s *IssueBoardsService) UpdateIssueBoardList(pid any, board, list int64, opt *UpdateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	return do[*BoardList](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/boards/%d/lists/%d", ProjectID{pid}, board, list),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IssueBoardsService) DeleteIssueBoardList(pid any, board, list int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/boards/%d/lists/%d", ProjectID{pid}, board, list),
		withRequestOpts(options...),
	)
	return resp, err
}
