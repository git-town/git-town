//
// Copyright 2023, Hakki Ceylan, Yavuz Turk
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
	ResourceIterationEventsServiceInterface interface {
		ListIssueIterationEvents(pid any, issue int64, opt *ListIterationEventsOptions, options ...RequestOptionFunc) ([]*IterationEvent, *Response, error)
		GetIssueIterationEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*IterationEvent, *Response, error)
	}

	// ResourceIterationEventsService handles communication with the event related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/resource_iteration_events/
	ResourceIterationEventsService struct {
		client *Client
	}
)

var _ ResourceIterationEventsServiceInterface = (*ResourceIterationEventsService)(nil)

// IterationEvent represents a resource iteration event.
//
// GitLab API docs: https://docs.gitlab.com/api/resource_iteration_events/
type IterationEvent struct {
	ID           int64      `json:"id"`
	User         *BasicUser `json:"user"`
	CreatedAt    *time.Time `json:"created_at"`
	ResourceType string     `json:"resource_type"`
	ResourceID   int64      `json:"resource_id"`
	Iteration    *Iteration `json:"iteration"`
	Action       string     `json:"action"`
}

// Iteration represents a project issue iteration.
//
// GitLab API docs: https://docs.gitlab.com/api/resource_iteration_events/
type Iteration struct {
	ID          int64      `json:"id"`
	IID         int64      `json:"iid"`
	Sequence    int64      `json:"sequence"`
	GroupID     int64      `json:"group_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       int64      `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DueDate     *ISOTime   `json:"due_date"`
	StartDate   *ISOTime   `json:"start_date"`
	WebURL      string     `json:"web_url"`
}

// ListIterationEventsOptions represents the options for all resource state
// events list methods.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_iteration_events/#list-project-issue-iteration-events
type ListIterationEventsOptions struct {
	ListOptions
}

// ListIssueIterationEvents retrieves resource iteration events for the
// specified project and issue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_iteration_events/#list-project-issue-iteration-events
func (s *ResourceIterationEventsService) ListIssueIterationEvents(pid any, issue int64, opt *ListIterationEventsOptions, options ...RequestOptionFunc) ([]*IterationEvent, *Response, error) {
	return do[[]*IterationEvent](s.client,
		withPath("projects/%s/issues/%d/resource_iteration_events", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetIssueIterationEvent gets a single issue iteration event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_iteration_events/#get-single-issue-iteration-event
func (s *ResourceIterationEventsService) GetIssueIterationEvent(pid any, issue int64, event int64, options ...RequestOptionFunc) (*IterationEvent, *Response, error) {
	return do[*IterationEvent](s.client,
		withPath("projects/%s/issues/%d/resource_iteration_events/%d", ProjectID{pid}, issue, event),
		withRequestOpts(options...),
	)
}
