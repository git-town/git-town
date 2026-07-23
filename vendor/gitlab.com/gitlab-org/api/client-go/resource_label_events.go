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
	"time"
)

type (
	ResourceLabelEventsServiceInterface interface {
		ListIssueLabelEvents(pid any, issue int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error)
		GetIssueLabelEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error)
		ListMergeRequestsLabelEvents(pid any, request int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error)
		GetMergeRequestLabelEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error)

		// Will be removed in v5, use Work Items API instead
		ListGroupEpicLabelEvents(gid any, epic int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error)
		// Will be removed in v5, use Work Items API instead
		GetGroupEpicLabelEvent(gid any, epic int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error)
	}

	// ResourceLabelEventsService handles communication with the event related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/resource_label_events/
	ResourceLabelEventsService struct {
		client *Client
	}
)

var _ ResourceLabelEventsServiceInterface = (*ResourceLabelEventsService)(nil)

// LabelEvent represents a resource label event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#get-single-issue-label-event
type LabelEvent struct {
	ID           int64           `json:"id"`
	Action       string          `json:"action"`
	CreatedAt    *time.Time      `json:"created_at"`
	ResourceType string          `json:"resource_type"`
	ResourceID   int64           `json:"resource_id"`
	User         BasicUser       `json:"user"`
	Label        LabelEventLabel `json:"label"`
}

// LabelEventLabel represents a resource label event label.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#get-single-issue-label-event
type LabelEventLabel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	TextColor   string `json:"text_color"`
	Description string `json:"description"`
}

// ListLabelEventsOptions represents the options for all resource label events
// list methods.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#list-project-issue-label-events
type ListLabelEventsOptions struct {
	ListOptions
}

// ListIssueLabelEvents retrieves resource label events for the
// specified project and issue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#list-project-issue-label-events
func (s *ResourceLabelEventsService) ListIssueLabelEvents(pid any, issue int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error) {
	return do[[]*LabelEvent](s.client,
		withPath("projects/%s/issues/%d/resource_label_events", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetIssueLabelEvent gets a single issue-label-event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#get-single-issue-label-event
func (s *ResourceLabelEventsService) GetIssueLabelEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error) {
	return do[*LabelEvent](s.client,
		withPath("projects/%s/issues/%d/resource_label_events/%d", ProjectID{pid}, issue, event),
		withRequestOpts(options...),
	)
}

// ListGroupEpicLabelEvents retrieves resource label events for the specified
// group and epic.
// Will be removed in v5, use Work Items API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#list-group-epic-label-events
func (s *ResourceLabelEventsService) ListGroupEpicLabelEvents(gid any, epic int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error) {
	return do[[]*LabelEvent](s.client,
		withPath("groups/%s/epics/%d/resource_label_events", GroupID{gid}, epic),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupEpicLabelEvent gets a single group epic label event.
// Will be removed in v5, use Work Items API instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#get-single-epic-label-event
func (s *ResourceLabelEventsService) GetGroupEpicLabelEvent(gid any, epic int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error) {
	return do[*LabelEvent](s.client,
		withPath("groups/%s/epics/%d/resource_label_events/%d", GroupID{gid}, epic, event),
		withRequestOpts(options...),
	)
}

// ListMergeRequestsLabelEvents retrieves resource label events for the specified
// project and merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#list-project-merge-request-label-events
func (s *ResourceLabelEventsService) ListMergeRequestsLabelEvents(pid any, request int64, opt *ListLabelEventsOptions, options ...RequestOptionFunc) ([]*LabelEvent, *Response, error) {
	return do[[]*LabelEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_label_events", ProjectID{pid}, request),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetMergeRequestLabelEvent gets a single merge request label event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_label_events/#get-single-merge-request-label-event
func (s *ResourceLabelEventsService) GetMergeRequestLabelEvent(pid any, request int64, event int64, options ...RequestOptionFunc) (*LabelEvent, *Response, error) {
	return do[*LabelEvent](s.client,
		withPath("projects/%s/merge_requests/%d/resource_label_events/%d", ProjectID{pid}, request, event),
		withRequestOpts(options...),
	)
}
