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
	// EpicIssuesServiceInterface defines all the API methods for the EpicIssuesService
	// Will be removed in v5 of the API, use Work Items API instead
	EpicIssuesServiceInterface interface {
		// Will be removed in v5 of the API, use Work Items API instead
		ListEpicIssues(gid any, epic int, opt *ListOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		AssignEpicIssue(gid any, epic, issue int, options ...RequestOptionFunc) (*EpicIssueAssignment, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		RemoveEpicIssue(gid any, epic, epicIssue int, options ...RequestOptionFunc) (*EpicIssueAssignment, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		UpdateEpicIssueAssignment(gid any, epic, epicIssue int, opt *UpdateEpicIssueAssignmentOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
	}

	// EpicIssuesService handles communication with the epic issue related methods
	// of the GitLab API.
	// Will be removed in v5 of the API, use Work Items API instead
	//
	// GitLab API docs: https://docs.gitlab.com/api/epic_issues/
	EpicIssuesService struct {
		client *Client
	}
)

// Will be removed in v5 of the API, use Work Items API instead
var _ EpicIssuesServiceInterface = (*EpicIssuesService)(nil)

// EpicIssueAssignment contains both the epic and issue objects returned from
// Gitlab with the assignment ID.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epic_issues/
type EpicIssueAssignment struct {
	ID    int    `json:"id"`
	Epic  *Epic  `json:"epic"`
	Issue *Issue `json:"issue"`
}

// ListEpicIssues get a list of epic issues.
// Will be removed in v5 of the API, use Work Items API instead
//
// Gitlab API docs:
// https://docs.gitlab.com/api/epic_issues/#list-issues-for-an-epic
func (s *EpicIssuesService) ListEpicIssues(gid any, epic int, opt *ListOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/epics/%d/issues", PathEscape(group), epic)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var is []*Issue
	resp, err := s.client.Do(req, &is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}

// AssignEpicIssue assigns an existing issue to an epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// Gitlab API Docs:
// https://docs.gitlab.com/api/epic_issues/#assign-an-issue-to-the-epic
func (s *EpicIssuesService) AssignEpicIssue(gid any, epic, issue int, options ...RequestOptionFunc) (*EpicIssueAssignment, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/epics/%d/issues/%d", PathEscape(group), epic, issue)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(EpicIssueAssignment)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// RemoveEpicIssue removes an issue from an epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// Gitlab API Docs:
// https://docs.gitlab.com/api/epic_issues/#remove-an-issue-from-the-epic
func (s *EpicIssuesService) RemoveEpicIssue(gid any, epic, epicIssue int, options ...RequestOptionFunc) (*EpicIssueAssignment, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/epics/%d/issues/%d", PathEscape(group), epic, epicIssue)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	a := new(EpicIssueAssignment)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// UpdateEpicIssueAssignmentOptions describes the UpdateEpicIssueAssignment()
// options.
// Will be removed in v5 of the API, use Work Items API instead
//
// Gitlab API Docs:
// https://docs.gitlab.com/api/epic_issues/#update-epic---issue-association
type UpdateEpicIssueAssignmentOptions struct {
	*ListOptions
	MoveBeforeID *int `url:"move_before_id,omitempty" json:"move_before_id,omitempty"`
	MoveAfterID  *int `url:"move_after_id,omitempty" json:"move_after_id,omitempty"`
}

// UpdateEpicIsssueAssignmentOptions is kept for backwards compatibility.
// Deprecated: use UpdateEpicIssueAssignmentOptions instead.
type UpdateEpicIsssueAssignmentOptions = UpdateEpicIssueAssignmentOptions

// UpdateEpicIssueAssignment moves an issue before or after another issue in an
// epic issue list.
// Will be removed in v5 of the API, use Work Items API instead
//
// Gitlab API Docs:
// https://docs.gitlab.com/api/epic_issues/#update-epic---issue-association
func (s *EpicIssuesService) UpdateEpicIssueAssignment(gid any, epic, epicIssue int, opt *UpdateEpicIssueAssignmentOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/epics/%d/issues/%d", PathEscape(group), epic, epicIssue)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var is []*Issue
	resp, err := s.client.Do(req, &is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}
