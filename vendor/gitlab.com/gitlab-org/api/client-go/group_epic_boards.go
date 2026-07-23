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

type (
	GroupEpicBoardsServiceInterface interface {
		ListGroupEpicBoards(gid any, opt *ListGroupEpicBoardsOptions, options ...RequestOptionFunc) ([]*GroupEpicBoard, *Response, error)
		GetGroupEpicBoard(gid any, board int64, options ...RequestOptionFunc) (*GroupEpicBoard, *Response, error)
	}

	// GroupEpicBoardsService handles communication with the group epic board
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_epic_boards/
	GroupEpicBoardsService struct {
		client *Client
	}
)

var _ GroupEpicBoardsServiceInterface = (*GroupEpicBoardsService)(nil)

// GroupEpicBoard represents a GitLab group epic board.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_epic_boards/
type GroupEpicBoard struct {
	ID     int64           `json:"id"`
	Name   string          `json:"name"`
	Group  *Group          `json:"group"`
	Labels []*LabelDetails `json:"labels"`
	Lists  []*BoardList    `json:"lists"`
}

func (b GroupEpicBoard) String() string {
	return Stringify(b)
}

// ListGroupEpicBoardsOptions represents the available
// ListGroupEpicBoards() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_epic_boards/#list-all-epic-boards-in-a-group
type ListGroupEpicBoardsOptions struct {
	ListOptions
}

// ListGroupEpicBoards gets a list of all epic boards in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_epic_boards/#list-all-epic-boards-in-a-group
func (s *GroupEpicBoardsService) ListGroupEpicBoards(gid any, opt *ListGroupEpicBoardsOptions, options ...RequestOptionFunc) ([]*GroupEpicBoard, *Response, error) {
	return do[[]*GroupEpicBoard](s.client,
		withPath("groups/%s/epic_boards", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupEpicBoard gets a single epic board of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_epic_boards/#single-group-epic-board
func (s *GroupEpicBoardsService) GetGroupEpicBoard(gid any, board int64, options ...RequestOptionFunc) (*GroupEpicBoard, *Response, error) {
	return do[*GroupEpicBoard](s.client,
		withPath("groups/%s/epic_boards/%d", GroupID{gid}, board),
		withRequestOpts(options...),
	)
}
