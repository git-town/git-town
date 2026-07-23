//
// Copyright 2022, Mai Lapyst
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
	"time"
)

type (
	ResourceMilestoneEventsServiceInterface interface {
		ListIssueMilestoneEvents(pid any, issue int64, opt *ListMilestoneEventsOptions, options ...RequestOptionFunc) ([]*MilestoneEvent, *Response, error)
		GetIssueMilestoneEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*MilestoneEvent, *Response, error)
		ListMergeMilestoneEvents(pid any, request int64, opt *ListMilestoneEventsOptions, options ...RequestOptionFunc) ([]*MilestoneEvent, *Response, error)
		GetMergeRequestMilestoneEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*MilestoneEvent, *Response, error)
	}

	// ResourceMilestoneEventsService handles communication with the event related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/resource_milestone_events/
	ResourceMilestoneEventsService struct {
		client *Client
	}
)

var _ ResourceMilestoneEventsServiceInterface = (*ResourceMilestoneEventsService)(nil)

// MilestoneEvent represents a resource milestone event.
//
// GitLab API docs: https://docs.gitlab.com/api/resource_milestone_events/
type MilestoneEvent struct {
	ID           int64      `json:"id"`
	User         *BasicUser `json:"user"`
	CreatedAt    *time.Time `json:"created_at"`
	ResourceType string     `json:"resource_type"`
	ResourceID   int64      `json:"resource_id"`
	Milestone    *Milestone `json:"milestone"`
	Action       string     `json:"action"`
}

// ListMilestoneEventsOptions represents the options for all resource state events
// list methods.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_milestone_events/#list-project-issue-milestone-events
type ListMilestoneEventsOptions struct {
	ListOptions
}

// ListIssueMilestoneEvents retrieves resource milestone events for the specified
// project and issue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_milestone_events/#list-project-issue-milestone-events
func (s *ResourceMilestoneEventsService) ListIssueMilestoneEvents(pid any, issue int64, opt *ListMilestoneEventsOptions, options ...RequestOptionFunc) ([]*MilestoneEvent, *Response, error) {
	return do[[]*MilestoneEvent](s.client,
		withPath("projects/%s/issues/%d/resource_milestone_events", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetIssueMilestoneEvent gets a single issue milestone event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_milestone_events/#get-single-issue-milestone-event
func (s *ResourceMilestoneEventsService) GetIssueMilestoneEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*MilestoneEvent, *Response, error) {
	return do[*MilestoneEvent](s.client,
		withPath("projects/%s/issues/%d/resource_milestone_events/%d", ProjectID{pid}, issue, event),
		withRequestOpts(options...),
	)
}

// ListMergeMilestoneEvents retrieves resource milestone events for the specified
// project and merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_milestone_events/#list-project-merge-request-milestone-events
func (s *ResourceMilestoneEventsService) ListMergeMilestoneEvents(pid any, request int64, opt *ListMilestoneEventsOptions, options ...RequestOptionFunc) ([]*MilestoneEvent, *Response, error) {
	return do[[]*MilestoneEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_milestone_events", ProjectID{pid}, request),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetMergeRequestMilestoneEvent gets a single merge request milestone event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_milestone_events/#get-single-merge-request-milestone-event
func (s *ResourceMilestoneEventsService) GetMergeRequestMilestoneEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*MilestoneEvent, *Response, error) {
	return do[*MilestoneEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_milestone_events/%d", ProjectID{pid}, request, event),
		withRequestOpts(options...),
	)
}
