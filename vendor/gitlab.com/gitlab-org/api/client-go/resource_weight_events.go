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
	ResourceWeightEventsServiceInterface interface {
		ListIssueWeightEvents(pid any, issue int64, opt *ListWeightEventsOptions, options ...RequestOptionFunc) ([]*WeightEvent, *Response, error)
	}

	// ResourceWeightEventsService handles communication with the event related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/resource_weight_events/
	ResourceWeightEventsService struct {
		client *Client
	}
)

var _ ResourceWeightEventsServiceInterface = (*ResourceWeightEventsService)(nil)

// WeightEvent represents a resource weight event.
//
// GitLab API docs: https://docs.gitlab.com/api/resource_weight_events/
type WeightEvent struct {
	ID           int64          `json:"id"`
	User         *BasicUser     `json:"user"`
	CreatedAt    *time.Time     `json:"created_at"`
	ResourceType string         `json:"resource_type"`
	ResourceID   int64          `json:"resource_id"`
	State        EventTypeValue `json:"state"`
	IssueID      int64          `json:"issue_id"`
	Weight       int64          `json:"weight"`
}

// ListWeightEventsOptions represents the options for all resource weight events
// list methods.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_weight_events/#list-project-issue-weight-events
type ListWeightEventsOptions struct {
	ListOptions
}

// ListIssueWeightEvents retrieves resource weight events for the specified
// project and issue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/resource_weight_events/#list-project-issue-weight-events
func (s *ResourceWeightEventsService) ListIssueWeightEvents(pid any, issue int64, opt *ListWeightEventsOptions, options ...RequestOptionFunc) ([]*WeightEvent, *Response, error) {
	return do[[]*WeightEvent](s.client,
		withPath("projects/%s/issues/%d/resource_weight_events", ProjectID{pid}, issue),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
