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
	"time"
)

type (
	GroupMilestonesServiceInterface interface {
		ListGroupMilestones(gid any, opt *ListGroupMilestonesOptions, options ...RequestOptionFunc) ([]*GroupMilestone, *Response, error)
		GetGroupMilestone(gid any, milestone int, options ...RequestOptionFunc) (*GroupMilestone, *Response, error)
		CreateGroupMilestone(gid any, opt *CreateGroupMilestoneOptions, options ...RequestOptionFunc) (*GroupMilestone, *Response, error)
		UpdateGroupMilestone(gid any, milestone int, opt *UpdateGroupMilestoneOptions, options ...RequestOptionFunc) (*GroupMilestone, *Response, error)
		DeleteGroupMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Response, error)
		GetGroupMilestoneIssues(gid any, milestone int, opt *GetGroupMilestoneIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		GetGroupMilestoneMergeRequests(gid any, milestone int, opt *GetGroupMilestoneMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)
		GetGroupMilestoneBurndownChartEvents(gid any, milestone int, opt *GetGroupMilestoneBurndownChartEventsOptions, options ...RequestOptionFunc) ([]*BurndownChartEvent, *Response, error)
	}

	// GroupMilestonesService handles communication with the milestone related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_milestones/
	GroupMilestonesService struct {
		client *Client
	}
)

var _ GroupMilestonesServiceInterface = (*GroupMilestonesService)(nil)

// GroupMilestone represents a GitLab milestone.
//
// GitLab API docs: https://docs.gitlab.com/api/group_milestones/
type GroupMilestone struct {
	ID          int        `json:"id"`
	IID         int        `json:"iid"`
	GroupID     int        `json:"group_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *ISOTime   `json:"start_date"`
	DueDate     *ISOTime   `json:"due_date"`
	State       string     `json:"state"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
	Expired     *bool      `json:"expired"`
}

func (m GroupMilestone) String() string {
	return Stringify(m)
}

// ListGroupMilestonesOptions represents the available
// ListGroupMilestones() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#list-group-milestones
type ListGroupMilestonesOptions struct {
	ListOptions
	IIDs               *[]int   `url:"iids[],omitempty" json:"iids,omitempty"`
	State              *string  `url:"state,omitempty" json:"state,omitempty"`
	Title              *string  `url:"title,omitempty" json:"title,omitempty"`
	Search             *string  `url:"search,omitempty" json:"search,omitempty"`
	SearchTitle        *string  `url:"search_title,omitempty" json:"search_title,omitempty"`
	IncludeAncestors   *bool    `url:"include_ancestors,omitempty" json:"include_ancestors,omitempty"`
	IncludeDescendents *bool    `url:"include_descendents,omitempty" json:"include_descendents,omitempty"`
	UpdatedBefore      *ISOTime `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	UpdatedAfter       *ISOTime `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	ContainingDate     *ISOTime `url:"containing_date,omitempty" json:"containing_date,omitempty"`
	StartDate          *ISOTime `url:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate            *ISOTime `url:"end_date,omitempty" json:"end_date,omitempty"`

	// Deprecated: in GitLab 16.7, use IncludeAncestors instead
	IncludeParentMilestones *bool `url:"include_parent_milestones,omitempty" json:"include_parent_milestones,omitempty"`
}

// ListGroupMilestones returns a list of group milestones.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#list-group-milestones
func (s *GroupMilestonesService) ListGroupMilestones(gid any, opt *ListGroupMilestonesOptions, options ...RequestOptionFunc) ([]*GroupMilestone, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var m []*GroupMilestone
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// GetGroupMilestone gets a single group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-single-milestone
func (s *GroupMilestonesService) GetGroupMilestone(gid any, milestone int, options ...RequestOptionFunc) (*GroupMilestone, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d", PathEscape(group), milestone)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(GroupMilestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateGroupMilestoneOptions represents the available CreateGroupMilestone() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#create-new-milestone
type CreateGroupMilestoneOptions struct {
	Title       *string  `url:"title,omitempty" json:"title,omitempty"`
	Description *string  `url:"description,omitempty" json:"description,omitempty"`
	StartDate   *ISOTime `url:"start_date,omitempty" json:"start_date,omitempty"`
	DueDate     *ISOTime `url:"due_date,omitempty" json:"due_date,omitempty"`
}

// CreateGroupMilestone creates a new group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#create-new-milestone
func (s *GroupMilestonesService) CreateGroupMilestone(gid any, opt *CreateGroupMilestoneOptions, options ...RequestOptionFunc) (*GroupMilestone, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(GroupMilestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// UpdateGroupMilestoneOptions represents the available UpdateGroupMilestone() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#edit-milestone
type UpdateGroupMilestoneOptions struct {
	Title       *string  `url:"title,omitempty" json:"title,omitempty"`
	Description *string  `url:"description,omitempty" json:"description,omitempty"`
	StartDate   *ISOTime `url:"start_date,omitempty" json:"start_date,omitempty"`
	DueDate     *ISOTime `url:"due_date,omitempty" json:"due_date,omitempty"`
	StateEvent  *string  `url:"state_event,omitempty" json:"state_event,omitempty"`
}

// UpdateGroupMilestone updates an existing group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#edit-milestone
func (s *GroupMilestonesService) UpdateGroupMilestone(gid any, milestone int, opt *UpdateGroupMilestoneOptions, options ...RequestOptionFunc) (*GroupMilestone, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d", PathEscape(group), milestone)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(GroupMilestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteGroupMilestone deletes a specified group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#delete-group-milestone
func (s *GroupMilestonesService) DeleteGroupMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d", PathEscape(project), milestone)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req, nil)
}

// GetGroupMilestoneIssuesOptions represents the available GetGroupMilestoneIssues() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-issues-assigned-to-a-single-milestone
type GetGroupMilestoneIssuesOptions ListOptions

// GetGroupMilestoneIssues gets all issues assigned to a single group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-issues-assigned-to-a-single-milestone
func (s *GroupMilestonesService) GetGroupMilestoneIssues(gid any, milestone int, opt *GetGroupMilestoneIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d/issues", PathEscape(group), milestone)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var i []*Issue
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, nil
}

// GetGroupMilestoneMergeRequestsOptions represents the available
// GetGroupMilestoneMergeRequests() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-merge-requests-assigned-to-a-single-milestone
type GetGroupMilestoneMergeRequestsOptions ListOptions

// GetGroupMilestoneMergeRequests gets all merge requests assigned to a
// single group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-merge-requests-assigned-to-a-single-milestone
func (s *GroupMilestonesService) GetGroupMilestoneMergeRequests(gid any, milestone int, opt *GetGroupMilestoneMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d/merge_requests", PathEscape(group), milestone)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var mr []*BasicMergeRequest
	resp, err := s.client.Do(req, &mr)
	if err != nil {
		return nil, resp, err
	}

	return mr, resp, nil
}

// BurndownChartEvent reprensents a burnout chart event
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-burndown-chart-events-for-a-single-milestone
type BurndownChartEvent struct {
	CreatedAt *time.Time `json:"created_at"`
	Weight    *int       `json:"weight"`
	Action    *string    `json:"action"`
}

// GetGroupMilestoneBurndownChartEventsOptions represents the available
// GetGroupMilestoneBurndownChartEventsOptions() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-burndown-chart-events-for-a-single-milestone
type GetGroupMilestoneBurndownChartEventsOptions ListOptions

// GetGroupMilestoneBurndownChartEvents gets all merge requests assigned to a
// single group milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_milestones/#get-all-burndown-chart-events-for-a-single-milestone
func (s *GroupMilestonesService) GetGroupMilestoneBurndownChartEvents(gid any, milestone int, opt *GetGroupMilestoneBurndownChartEventsOptions, options ...RequestOptionFunc) ([]*BurndownChartEvent, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/milestones/%d/burndown_events", PathEscape(group), milestone)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var be []*BurndownChartEvent
	resp, err := s.client.Do(req, &be)
	if err != nil {
		return nil, resp, err
	}

	return be, resp, nil
}
