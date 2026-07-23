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
	"net/http"
	"time"
)

type (
	// EpicsServiceInterface defines all the API methods for the EpicsService
	// Will be removed in v5 of the API, use Work Items API instead
	EpicsServiceInterface interface {
		// ListGroupEpics gets a list of group epics. This function accepts pagination
		// parameters page and per_page to return the list of group epics.
		// Will be removed in v5 of the API, use Work Items API instead
		//
		// GitLab API docs: https://docs.gitlab.com/api/epics/#list-epics-for-a-group
		ListGroupEpics(gid any, opt *ListGroupEpicsOptions, options ...RequestOptionFunc) ([]*Epic, *Response, error)

		// GetEpic gets a single group epic.
		// Will be removed in v5 of the API, use Work Items API instead
		GetEpic(gid any, epic int64, options ...RequestOptionFunc) (*Epic, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		GetEpicLinks(gid any, epic int64, options ...RequestOptionFunc) ([]*Epic, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		//
		// GitLab API docs: https://docs.gitlab.com/api/epics/#new-epic
		CreateEpic(gid any, opt *CreateEpicOptions, options ...RequestOptionFunc) (*Epic, *Response, error)

		// UpdateEpic updates an existing group epic. This function is also used
		// to mark an epic as closed.
		// Will be removed in v5 of the API, use Work Items API instead
		UpdateEpic(gid any, epic int64, opt *UpdateEpicOptions, options ...RequestOptionFunc) (*Epic, *Response, error)
		// Will be removed in v5 of the API, use Work Items API instead
		DeleteEpic(gid any, epic int64, options ...RequestOptionFunc) (*Response, error)
	}

	// EpicsService handles communication with the epic related methods
	// of the GitLab API.
	// Will be removed in v5 of the API, use Work Items API instead
	//
	// GitLab API docs: https://docs.gitlab.com/api/epics/
	EpicsService struct {
		client *Client
	}
)

// Will be removed in v5 of the API, use Work Items API instead
var _ EpicsServiceInterface = (*EpicsService)(nil)

// EpicAuthor represents a author of the epic.
// Will be removed in v5 of the API, use Work Items API instead
type EpicAuthor struct {
	ID        int64  `json:"id"`
	State     string `json:"state"`
	WebURL    string `json:"web_url"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Username  string `json:"username"`
}

// Epic represents a GitLab epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/
type Epic struct {
	ID                      int64       `json:"id"`
	IID                     int64       `json:"iid"`
	GroupID                 int64       `json:"group_id"`
	ParentID                int64       `json:"parent_id"`
	Title                   string      `json:"title"`
	Description             string      `json:"description"`
	State                   string      `json:"state"`
	Confidential            bool        `json:"confidential"`
	WebURL                  string      `json:"web_url"`
	Author                  *EpicAuthor `json:"author"`
	StartDate               *ISOTime    `json:"start_date"`
	StartDateIsFixed        bool        `json:"start_date_is_fixed"`
	StartDateFixed          *ISOTime    `json:"start_date_fixed"`
	StartDateFromMilestones *ISOTime    `json:"start_date_from_milestones"`
	DueDate                 *ISOTime    `json:"due_date"`
	DueDateIsFixed          bool        `json:"due_date_is_fixed"`
	DueDateFixed            *ISOTime    `json:"due_date_fixed"`
	DueDateFromMilestones   *ISOTime    `json:"due_date_from_milestones"`
	CreatedAt               *time.Time  `json:"created_at"`
	UpdatedAt               *time.Time  `json:"updated_at"`
	ClosedAt                *time.Time  `json:"closed_at"`
	Labels                  []string    `json:"labels"`
	Upvotes                 int64       `json:"upvotes"`
	Downvotes               int64       `json:"downvotes"`
	UserNotesCount          int64       `json:"user_notes_count"`
	URL                     string      `json:"url"`
}

// String gets a string representation of an Epic.
//
// Will be removed in v5 of the API, use Work Items API instead
func (e Epic) String() string {
	return Stringify(e)
}

// ListGroupEpicsOptions represents the available ListGroupEpics() options.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#list-epics-for-a-group
type ListGroupEpicsOptions struct {
	ListOptions
	AuthorID                *int64        `url:"author_id,omitempty" json:"author_id,omitempty"`
	Labels                  *LabelOptions `url:"labels,comma,omitempty" json:"labels,omitempty"`
	WithLabelDetails        *bool         `url:"with_labels_details,omitempty" json:"with_labels_details,omitempty"`
	OrderBy                 *string       `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                    *string       `url:"sort,omitempty" json:"sort,omitempty"`
	Search                  *string       `url:"search,omitempty" json:"search,omitempty"`
	State                   *string       `url:"state,omitempty" json:"state,omitempty"`
	CreatedAfter            *time.Time    `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore           *time.Time    `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter            *time.Time    `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore           *time.Time    `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	IncludeAncestorGroups   *bool         `url:"include_ancestor_groups,omitempty" json:"include_ancestor_groups,omitempty"`
	IncludeDescendantGroups *bool         `url:"include_descendant_groups,omitempty" json:"include_descendant_groups,omitempty"`
	MyReactionEmoji         *string       `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
}

func (s *EpicsService) ListGroupEpics(gid any, opt *ListGroupEpicsOptions, options ...RequestOptionFunc) ([]*Epic, *Response, error) {
	return do[[]*Epic](s.client,
		withPath("groups/%s/epics", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetEpic gets a single group epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#single-epic
func (s *EpicsService) GetEpic(gid any, epic int64, options ...RequestOptionFunc) (*Epic, *Response, error) {
	return do[*Epic](s.client,
		withPath("groups/%s/epics/%d", GroupID{gid}, epic),
		withRequestOpts(options...),
	)
}

// GetEpicLinks gets all child epics of an epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epic_links/
func (s *EpicsService) GetEpicLinks(gid any, epic int64, options ...RequestOptionFunc) ([]*Epic, *Response, error) {
	return do[[]*Epic](s.client,
		withPath("groups/%s/epics/%d/epics", GroupID{gid}, epic),
		withRequestOpts(options...),
	)
}

// CreateEpicOptions represents the available CreateEpic() options.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#new-epic
type CreateEpicOptions struct {
	Title            *string       `url:"title,omitempty" json:"title,omitempty"`
	Labels           *LabelOptions `url:"labels,comma,omitempty" json:"labels,omitempty"`
	Description      *string       `url:"description,omitempty" json:"description,omitempty"`
	Color            *string       `url:"color,omitempty" json:"color,omitempty"`
	Confidential     *bool         `url:"confidential,omitempty" json:"confidential,omitempty"`
	CreatedAt        *time.Time    `url:"created_at,omitempty" json:"created_at,omitempty"`
	StartDateIsFixed *bool         `url:"start_date_is_fixed,omitempty" json:"start_date_is_fixed,omitempty"`
	StartDateFixed   *ISOTime      `url:"start_date_fixed,omitempty" json:"start_date_fixed,omitempty"`
	DueDateIsFixed   *bool         `url:"due_date_is_fixed,omitempty" json:"due_date_is_fixed,omitempty"`
	DueDateFixed     *ISOTime      `url:"due_date_fixed,omitempty" json:"due_date_fixed,omitempty"`
	ParentID         *int64        `url:"parent_id,omitempty" json:"parent_id,omitempty"`
}

func (s *EpicsService) CreateEpic(gid any, opt *CreateEpicOptions, options ...RequestOptionFunc) (*Epic, *Response, error) {
	return do[*Epic](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/epics", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateEpicOptions represents the available UpdateEpic() options.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#update-epic
type UpdateEpicOptions struct {
	AddLabels        *LabelOptions `url:"add_labels,omitempty" json:"add_labels,omitempty"`
	Confidential     *bool         `url:"confidential,omitempty" json:"confidential,omitempty"`
	Description      *string       `url:"description,omitempty" json:"description,omitempty"`
	DueDateFixed     *ISOTime      `url:"due_date_fixed,omitempty" json:"due_date_fixed,omitempty"`
	DueDateIsFixed   *bool         `url:"due_date_is_fixed,omitempty" json:"due_date_is_fixed,omitempty"`
	Labels           *LabelOptions `url:"labels,comma,omitempty" json:"labels,omitempty"`
	ParentID         *int64        `url:"parent_id,omitempty" json:"parent_id,omitempty"`
	RemoveLabels     *LabelOptions `url:"remove_labels,omitempty" json:"remove_labels,omitempty"`
	StartDateFixed   *ISOTime      `url:"start_date_fixed,omitempty" json:"start_date_fixed,omitempty"`
	StartDateIsFixed *bool         `url:"start_date_is_fixed,omitempty" json:"start_date_is_fixed,omitempty"`
	StateEvent       *string       `url:"state_event,omitempty" json:"state_event,omitempty"`
	Title            *string       `url:"title,omitempty" json:"title,omitempty"`
	UpdatedAt        *time.Time    `url:"updated_at,omitempty" json:"updated_at,omitempty"`
	Color            *string       `url:"color,omitempty" json:"color,omitempty"`
}

// UpdateEpic updates an existing group epic. This function is also used
// to mark an epic as closed.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#update-epic
func (s *EpicsService) UpdateEpic(gid any, epic int64, opt *UpdateEpicOptions, options ...RequestOptionFunc) (*Epic, *Response, error) {
	return do[*Epic](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/epics/%d", GroupID{gid}, epic),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteEpic deletes a single group epic.
// Will be removed in v5 of the API, use Work Items API instead
//
// GitLab API docs: https://docs.gitlab.com/api/epics/#delete-epic
func (s *EpicsService) DeleteEpic(gid any, epic int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/epics/%d", GroupID{gid}, epic),
		withRequestOpts(options...),
	)
	return resp, err
}
