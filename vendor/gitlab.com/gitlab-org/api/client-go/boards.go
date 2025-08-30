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
	"fmt"
	"net/http"
)

type (
	IssueBoardsServiceInterface interface {
		CreateIssueBoard(pid any, opt *CreateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error)
		UpdateIssueBoard(pid any, board int, opt *UpdateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error)
		DeleteIssueBoard(pid any, board int, options ...RequestOptionFunc) (*Response, error)
		ListIssueBoards(pid any, opt *ListIssueBoardsOptions, options ...RequestOptionFunc) ([]*IssueBoard, *Response, error)
		GetIssueBoard(pid any, board int, options ...RequestOptionFunc) (*IssueBoard, *Response, error)
		GetIssueBoardLists(pid any, board int, opt *GetIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error)
		GetIssueBoardList(pid any, board, list int, options ...RequestOptionFunc) (*BoardList, *Response, error)
		CreateIssueBoardList(pid any, board int, opt *CreateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error)
		UpdateIssueBoardList(pid any, board, list int, opt *UpdateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error)
		DeleteIssueBoardList(pid any, board, list int, options ...RequestOptionFunc) (*Response, error)
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
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Project   *Project   `json:"project"`
	Milestone *Milestone `json:"milestone"`
	Assignee  *struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignee"`
	Lists  []*BoardList    `json:"lists"`
	Weight int             `json:"weight"`
	Labels []*LabelDetails `json:"labels"`
}

func (b IssueBoard) String() string {
	return Stringify(b)
}

// BoardList represents a GitLab board list.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/
type BoardList struct {
	ID       int `json:"id"`
	Assignee *struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"assignee"`
	Iteration      *ProjectIteration `json:"iteration"`
	Label          *Label            `json:"label"`
	MaxIssueCount  int               `json:"max_issue_count"`
	MaxIssueWeight int               `json:"max_issue_weight"`
	Milestone      *Milestone        `json:"milestone"`
	Position       int               `json:"position"`
}

func (b BoardList) String() string {
	return Stringify(b)
}

// CreateIssueBoardOptions represents the available CreateIssueBoard() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-an-issue-board
type CreateIssueBoardOptions struct {
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

// CreateIssueBoard creates a new issue board.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-an-issue-board
func (s *IssueBoardsService) CreateIssueBoard(pid any, opt *CreateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	board := new(IssueBoard)
	resp, err := s.client.Do(req, board)
	if err != nil {
		return nil, resp, err
	}

	return board, resp, nil
}

// UpdateIssueBoardOptions represents the available UpdateIssueBoard() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#update-an-issue-board
type UpdateIssueBoardOptions struct {
	Name        *string       `url:"name,omitempty" json:"name,omitempty"`
	AssigneeID  *int          `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MilestoneID *int          `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	Labels      *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	Weight      *int          `url:"weight,omitempty" json:"weight,omitempty"`
}

// UpdateIssueBoard update an issue board.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#update-an-issue-board
func (s *IssueBoardsService) UpdateIssueBoard(pid any, board int, opt *UpdateIssueBoardOptions, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d", PathEscape(project), board)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	is := new(IssueBoard)
	resp, err := s.client.Do(req, is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}

// DeleteIssueBoard deletes an issue board.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#delete-an-issue-board
func (s *IssueBoardsService) DeleteIssueBoard(pid any, board int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d", PathEscape(project), board)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListIssueBoardsOptions represents the available ListIssueBoards() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-project-issue-boards
type ListIssueBoardsOptions ListOptions

// ListIssueBoards gets a list of all issue boards in a project.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-project-issue-boards
func (s *IssueBoardsService) ListIssueBoards(pid any, opt *ListIssueBoardsOptions, options ...RequestOptionFunc) ([]*IssueBoard, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var is []*IssueBoard
	resp, err := s.client.Do(req, &is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}

// GetIssueBoard gets a single issue board of a project.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#show-a-single-issue-board
func (s *IssueBoardsService) GetIssueBoard(pid any, board int, options ...RequestOptionFunc) (*IssueBoard, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d", PathEscape(project), board)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	ib := new(IssueBoard)
	resp, err := s.client.Do(req, ib)
	if err != nil {
		return nil, resp, err
	}

	return ib, resp, nil
}

// GetIssueBoardListsOptions represents the available GetIssueBoardLists() options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-board-lists-in-a-project-issue-board
type GetIssueBoardListsOptions ListOptions

// GetIssueBoardLists gets a list of the issue board's lists. Does not include
// backlog and closed lists.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#list-board-lists-in-a-project-issue-board
func (s *IssueBoardsService) GetIssueBoardLists(pid any, board int, opt *GetIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d/lists", PathEscape(project), board)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var bl []*BoardList
	resp, err := s.client.Do(req, &bl)
	if err != nil {
		return nil, resp, err
	}

	return bl, resp, nil
}

// GetIssueBoardList gets a single issue board list.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#show-a-single-board-list
func (s *IssueBoardsService) GetIssueBoardList(pid any, board, list int, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d/lists/%d",
		PathEscape(project),
		board,
		list,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	bl := new(BoardList)
	resp, err := s.client.Do(req, bl)
	if err != nil {
		return nil, resp, err
	}

	return bl, resp, nil
}

// CreateIssueBoardListOptions represents the available CreateIssueBoardList()
// options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-a-board-list
type CreateIssueBoardListOptions struct {
	LabelID     *int `url:"label_id,omitempty" json:"label_id,omitempty"`
	AssigneeID  *int `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MilestoneID *int `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	IterationID *int `url:"iteration_id,omitempty" json:"iteration_id,omitempty"`
}

// CreateIssueBoardList creates a new issue board list.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#create-a-board-list
func (s *IssueBoardsService) CreateIssueBoardList(pid any, board int, opt *CreateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d/lists", PathEscape(project), board)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	bl := new(BoardList)
	resp, err := s.client.Do(req, bl)
	if err != nil {
		return nil, resp, err
	}

	return bl, resp, nil
}

// UpdateIssueBoardListOptions represents the available UpdateIssueBoardList()
// options.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#reorder-a-list-in-a-board
type UpdateIssueBoardListOptions struct {
	Position *int `url:"position" json:"position"`
}

// UpdateIssueBoardList updates the position of an existing issue board list.
//
// GitLab API docs: https://docs.gitlab.com/api/boards/#reorder-a-list-in-a-board
func (s *IssueBoardsService) UpdateIssueBoardList(pid any, board, list int, opt *UpdateIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d/lists/%d",
		PathEscape(project),
		board,
		list,
	)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	bl := new(BoardList)
	resp, err := s.client.Do(req, bl)
	if err != nil {
		return nil, resp, err
	}

	return bl, resp, nil
}

// DeleteIssueBoardList soft deletes an issue board list. Only for admins and
// project owners.
//
// GitLab API docs:
// https://docs.gitlab.com/api/boards/#delete-a-board-list-from-a-board
func (s *IssueBoardsService) DeleteIssueBoardList(pid any, board, list int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/boards/%d/lists/%d",
		PathEscape(project),
		board,
		list,
	)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
