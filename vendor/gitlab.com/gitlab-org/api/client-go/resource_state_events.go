//
// Copyright 2021, Matthias Simon
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
	ResourceStateEventsServiceInterface interface {
		ListIssueStateEvents(pid any, issue int64, opt *ListStateEventsOptions, options ...RequestOptionFunc) ([]*StateEvent, *Response, error)
		GetIssueStateEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*StateEvent, *Response, error)
		ListMergeStateEvents(pid any, request int64, opt *ListStateEventsOptions, options ...RequestOptionFunc) ([]*StateEvent, *Response, error)
		GetMergeRequestStateEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*StateEvent, *Response, error)
	}

	// ResourceStateEventsService handles communication with the event related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/resource_state_events/
	ResourceStateEventsService struct {
		client *Client
	}
)

var _ ResourceStateEventsServiceInterface = (*ResourceStateEventsService)(nil)

// StateEvent represents a resource state event.
//
// GitLab API docs: https://docs.gitlab.com/api/resource_state_events/
type StateEvent struct {
	ID           int64          `json:"id"`
	User         *BasicUser     `json:"user"`
	CreatedAt    *time.Time     `json:"created_at"`
	ResourceType string         `json:"resource_type"`
	ResourceID   int64          `json:"resource_id"`
	State        EventTypeValue `json:"state"`
}

// ListStateEventsOptions represents the options for all resource state events
// list methods.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_state_events/#list-project-issue-state-events
type ListStateEventsOptions struct {
	ListOptions
}

// ListIssueStateEvents retrieves resource state events for the specified
// project and issue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_state_events/#list-project-issue-state-events
func (s *ResourceStateEventsService) ListIssueStateEvents(pid any, issue int64, opt *ListStateEventsOptions, options ...RequestOptionFunc) ([]*StateEvent, *Response, error) {
	return do[[]*StateEvent](s.client,
		withPath("projects/%s/issues/%d/resource_state_events", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetIssueStateEvent gets a single issue-state-event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_state_events/#get-single-issue-state-event
func (s *ResourceStateEventsService) GetIssueStateEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*StateEvent, *Response, error) {
	return do[*StateEvent](s.client,
		withPath("projects/%s/issues/%d/resource_state_events/%d", ProjectID{pid}, issue, event),
		withRequestOpts(options...),
	)
}

// ListMergeStateEvents retrieves resource state events for the specified
// project and merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_state_events/#list-project-merge-request-state-events
func (s *ResourceStateEventsService) ListMergeStateEvents(pid any, request int64, opt *ListStateEventsOptions, options ...RequestOptionFunc) ([]*StateEvent, *Response, error) {
	return do[[]*StateEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_state_events", ProjectID{pid}, request),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetMergeRequestStateEvent gets a single merge request state event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_state_events/#get-single-merge-request-state-event
func (s *ResourceStateEventsService) GetMergeRequestStateEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*StateEvent, *Response, error) {
	return do[*StateEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_state_events/%d", ProjectID{pid}, request, event),
		withRequestOpts(options...),
	)
}
