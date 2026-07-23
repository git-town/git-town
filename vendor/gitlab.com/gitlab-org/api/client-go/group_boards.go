//
// Copyright 2021, Patrick Webster
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
	GroupIssueBoardsServiceInterface interface {
		// ListGroupIssueBoards gets a list of all issue boards in a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#list-all-group-issue-boards-in-a-group
		ListGroupIssueBoards(gid any, opt *ListGroupIssueBoardsOptions, options ...RequestOptionFunc) ([]*GroupIssueBoard, *Response, error)
		// CreateGroupIssueBoard creates a new issue board.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#create-a-group-issue-board
		CreateGroupIssueBoard(gid any, opt *CreateGroupIssueBoardOptions, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error)
		// GetGroupIssueBoard gets a single issue board of a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#single-group-issue-board
		GetGroupIssueBoard(gid any, board int64, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error)
		// UpdateIssueBoard updates a single issue board of a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#update-a-group-issue-board
		UpdateIssueBoard(gid any, board int64, opt *UpdateGroupIssueBoardOptions, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error)
		// DeleteIssueBoard deletes a single issue board of a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#delete-a-group-issue-board
		DeleteIssueBoard(gid any, board int64, options ...RequestOptionFunc) (*Response, error)
		// ListGroupIssueBoardLists gets a list of the issue board's lists. Does not include
		// backlog and closed lists.
		//
		// GitLab API docs: https://docs.gitlab.com/api/group_boards/#list-group-issue-board-lists
		ListGroupIssueBoardLists(gid any, board int64, opt *ListGroupIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error)
		// GetGroupIssueBoardList gets a single issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#single-group-issue-board-list
		GetGroupIssueBoardList(gid any, board, list int64, options ...RequestOptionFunc) (*BoardList, *Response, error)
		// CreateGroupIssueBoardList creates a new issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#new-group-issue-board-list
		CreateGroupIssueBoardList(gid any, board int64, opt *CreateGroupIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error)
		// UpdateIssueBoardList updates the position of an existing
		// group issue board list.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_boards/#edit-group-issue-board-list
		UpdateIssueBoardList(gid any, board, list int64, opt *UpdateGroupIssueBoardListOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error)
		DeleteGroupIssueBoardList(gid any, board, list int64, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupIssueBoardsService handles communication with the group issue board
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_boards/
	GroupIssueBoardsService struct {
		client *Client
	}
)

var _ GroupIssueBoardsServiceInterface = (*GroupIssueBoardsService)(nil)

// GroupIssueBoard represents a GitLab group issue board.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/
type GroupIssueBoard struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Group     *Group        `json:"group"`
	Milestone *Milestone    `json:"milestone"`
	Labels    []*GroupLabel `json:"labels"`
	Lists     []*BoardList  `json:"lists"`
}

func (b GroupIssueBoard) String() string {
	return Stringify(b)
}

// ListGroupIssueBoardsOptions represents the available
// ListGroupIssueBoards() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#list-all-group-issue-boards-in-a-group
type ListGroupIssueBoardsOptions struct {
	ListOptions
}

func (s *GroupIssueBoardsService) ListGroupIssueBoards(gid any, opt *ListGroupIssueBoardsOptions, options ...RequestOptionFunc) ([]*GroupIssueBoard, *Response, error) {
	return do[[]*GroupIssueBoard](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/boards", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateGroupIssueBoardOptions represents the available
// CreateGroupIssueBoard() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#create-a-group-issue-board
type CreateGroupIssueBoardOptions struct {
	Name *string `url:"name" json:"name"`
}

func (s *GroupIssueBoardsService) CreateGroupIssueBoard(gid any, opt *CreateGroupIssueBoardOptions, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error) {
	return do[*GroupIssueBoard](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/boards", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupIssueBoardsService) GetGroupIssueBoard(gid any, board int64, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error) {
	return do[*GroupIssueBoard](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/boards/%d", GroupID{gid}, board),
		withRequestOpts(options...),
	)
}

// UpdateGroupIssueBoardOptions represents a group issue board.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#update-a-group-issue-board
type UpdateGroupIssueBoardOptions struct {
	Name        *string       `url:"name,omitempty" json:"name,omitempty"`
	AssigneeID  *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MilestoneID *int64        `url:"milestone_id,omitempty" json:"milestone_id,omitempty"`
	Labels      *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	Weight      *int64        `url:"weight,omitempty" json:"weight,omitempty"`
}

func (s *GroupIssueBoardsService) UpdateIssueBoard(gid any, board int64, opt *UpdateGroupIssueBoardOptions, options ...RequestOptionFunc) (*GroupIssueBoard, *Response, error) {
	return do[*GroupIssueBoard](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/boards/%d", GroupID{gid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupIssueBoardsService) DeleteIssueBoard(gid any, board int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/boards/%d", GroupID{gid}, board),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListGroupIssueBoardListsOptions represents the available
// ListGroupIssueBoardLists() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#list-group-issue-board-lists
type ListGroupIssueBoardListsOptions struct {
	ListOptions
}

func (s *GroupIssueBoardsService) ListGroupIssueBoardLists(gid any, board int64, opt *ListGroupIssueBoardListsOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error) {
	return do[[]*BoardList](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/boards/%d/lists", GroupID{gid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupIssueBoardsService) GetGroupIssueBoardList(gid any, board, list int64, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	return do[*BoardList](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/boards/%d/lists/%d", GroupID{gid}, board, list),
		withRequestOpts(options...),
	)
}

// CreateGroupIssueBoardListOptions represents the available
// CreateGroupIssueBoardList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#new-group-issue-board-list
type CreateGroupIssueBoardListOptions struct {
	LabelID *int64 `url:"label_id" json:"label_id"`
}

func (s *GroupIssueBoardsService) CreateGroupIssueBoardList(gid any, board int64, opt *CreateGroupIssueBoardListOptions, options ...RequestOptionFunc) (*BoardList, *Response, error) {
	return do[*BoardList](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/boards/%d/lists", GroupID{gid}, board),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateGroupIssueBoardListOptions represents the available
// UpdateGroupIssueBoardList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_boards/#edit-group-issue-board-list
type UpdateGroupIssueBoardListOptions struct {
	Position *int64 `url:"position" json:"position"`
}

func (s *GroupIssueBoardsService) UpdateIssueBoardList(gid any, board, list int64, opt *UpdateGroupIssueBoardListOptions, options ...RequestOptionFunc) ([]*BoardList, *Response, error) {
	return do[[]*BoardList](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/boards/%d/lists/%d", GroupID{gid}, board, list),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupIssueBoardsService) DeleteGroupIssueBoardList(gid any, board, list int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/boards/%d/lists/%d", GroupID{gid}, board, list),
		withRequestOpts(options...),
	)
	return resp, err
}
