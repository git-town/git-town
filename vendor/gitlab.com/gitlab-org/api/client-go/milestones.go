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
	MilestonesServiceInterface interface {
		ListMilestones(pid any, opt *ListMilestonesOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error)
		GetMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Milestone, *Response, error)
		CreateMilestone(pid any, opt *CreateMilestoneOptions, options ...RequestOptionFunc) (*Milestone, *Response, error)
		UpdateMilestone(pid any, milestone int, opt *UpdateMilestoneOptions, options ...RequestOptionFunc) (*Milestone, *Response, error)
		DeleteMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Response, error)
		GetMilestoneIssues(pid any, milestone int, opt *GetMilestoneIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		GetMilestoneMergeRequests(pid any, milestone int, opt *GetMilestoneMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error)
	}

	// MilestonesService handles communication with the milestone related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/milestones/
	MilestonesService struct {
		client *Client
	}
)

var _ MilestonesServiceInterface = (*MilestonesService)(nil)

// Milestone represents a GitLab milestone.
//
// GitLab API docs: https://docs.gitlab.com/api/milestones/
type Milestone struct {
	ID          int        `json:"id"`
	IID         int        `json:"iid"`
	GroupID     int        `json:"group_id"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *ISOTime   `json:"start_date"`
	DueDate     *ISOTime   `json:"due_date"`
	State       string     `json:"state"`
	WebURL      string     `json:"web_url"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
	Expired     *bool      `json:"expired"`
}

func (m Milestone) String() string {
	return Stringify(m)
}

// ListMilestonesOptions represents the available ListMilestones() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#list-project-milestones
type ListMilestonesOptions struct {
	ListOptions
	IIDs             *[]int  `url:"iids[],omitempty" json:"iids,omitempty"`
	Title            *string `url:"title,omitempty" json:"title,omitempty"`
	State            *string `url:"state,omitempty" json:"state,omitempty"`
	Search           *string `url:"search,omitempty" json:"search,omitempty"`
	IncludeAncestors *bool   `url:"include_ancestors,omitempty" json:"include_ancestors,omitempty"`

	// Deprecated: in GitLab 16,7, use IncludeAncestors instead
	IncludeParentMilestones *bool `url:"include_parent_milestones,omitempty" json:"include_parent_milestones,omitempty"`
}

// ListMilestones returns a list of project milestones.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#list-project-milestones
func (s *MilestonesService) ListMilestones(pid any, opt *ListMilestonesOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var m []*Milestone
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// GetMilestone gets a single project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#get-single-milestone
func (s *MilestonesService) GetMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Milestone, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones/%d", PathEscape(project), milestone)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(Milestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// CreateMilestoneOptions represents the available CreateMilestone() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#create-new-milestone
type CreateMilestoneOptions struct {
	Title       *string  `url:"title,omitempty" json:"title,omitempty"`
	Description *string  `url:"description,omitempty" json:"description,omitempty"`
	StartDate   *ISOTime `url:"start_date,omitempty" json:"start_date,omitempty"`
	DueDate     *ISOTime `url:"due_date,omitempty" json:"due_date,omitempty"`
}

// CreateMilestone creates a new project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#create-new-milestone
func (s *MilestonesService) CreateMilestone(pid any, opt *CreateMilestoneOptions, options ...RequestOptionFunc) (*Milestone, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(Milestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// UpdateMilestoneOptions represents the available UpdateMilestone() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#edit-milestone
type UpdateMilestoneOptions struct {
	Title       *string  `url:"title,omitempty" json:"title,omitempty"`
	Description *string  `url:"description,omitempty" json:"description,omitempty"`
	StartDate   *ISOTime `url:"start_date,omitempty" json:"start_date,omitempty"`
	DueDate     *ISOTime `url:"due_date,omitempty" json:"due_date,omitempty"`
	StateEvent  *string  `url:"state_event,omitempty" json:"state_event,omitempty"`
}

// UpdateMilestone updates an existing project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#edit-milestone
func (s *MilestonesService) UpdateMilestone(pid any, milestone int, opt *UpdateMilestoneOptions, options ...RequestOptionFunc) (*Milestone, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones/%d", PathEscape(project), milestone)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(Milestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// DeleteMilestone deletes a specified project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#delete-project-milestone
func (s *MilestonesService) DeleteMilestone(pid any, milestone int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones/%d", PathEscape(project), milestone)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req, nil)
}

// GetMilestoneIssuesOptions represents the available GetMilestoneIssues() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#get-all-issues-assigned-to-a-single-milestone
type GetMilestoneIssuesOptions ListOptions

// GetMilestoneIssues gets all issues assigned to a single project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#get-all-issues-assigned-to-a-single-milestone
func (s *MilestonesService) GetMilestoneIssues(pid any, milestone int, opt *GetMilestoneIssuesOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones/%d/issues", PathEscape(project), milestone)

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

// GetMilestoneMergeRequestsOptions represents the available
// GetMilestoneMergeRequests() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#get-all-merge-requests-assigned-to-a-single-milestone
type GetMilestoneMergeRequestsOptions ListOptions

// GetMilestoneMergeRequests gets all merge requests assigned to a single
// project milestone.
//
// GitLab API docs:
// https://docs.gitlab.com/api/milestones/#get-all-merge-requests-assigned-to-a-single-milestone
func (s *MilestonesService) GetMilestoneMergeRequests(pid any, milestone int, opt *GetMilestoneMergeRequestsOptions, options ...RequestOptionFunc) ([]*BasicMergeRequest, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/milestones/%d/merge_requests", PathEscape(project), milestone)

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
