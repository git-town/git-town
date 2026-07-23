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
	// EventsServiceInterface defines all the API methods for the EventsService
	EventsServiceInterface interface {
		// ListCurrentUserContributionEvents  retrieves all events
		// for the currently authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/events/#list-all-events
		ListCurrentUserContributionEvents(opt *ListContributionEventsOptions, options ...RequestOptionFunc) ([]*ContributionEvent, *Response, error)

		// ListProjectVisibleEvents gets the events for the specified project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/events/#list-a-projects-visible-events
		ListProjectVisibleEvents(pid any, opt *ListProjectVisibleEventsOptions, options ...RequestOptionFunc) ([]*ProjectEvent, *Response, error)
	}

	// EventsService handles communication with the event related methods of
	// the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/events/
	EventsService struct {
		client *Client
	}
)

var _ EventsServiceInterface = (*EventsService)(nil)

// ContributionEvent represents a user's contribution
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#get-user-contribution-events
type ContributionEvent struct {
	ID             int64                     `json:"id"`
	Title          string                    `json:"title"`
	ProjectID      int64                     `json:"project_id"`
	ActionName     string                    `json:"action_name"`
	TargetID       int64                     `json:"target_id"`
	TargetIID      int64                     `json:"target_iid"`
	TargetType     string                    `json:"target_type"`
	AuthorID       int64                     `json:"author_id"`
	TargetTitle    string                    `json:"target_title"`
	CreatedAt      *time.Time                `json:"created_at"`
	PushData       ContributionEventPushData `json:"push_data"`
	Note           *Note                     `json:"note"`
	Author         BasicUser                 `json:"author"`
	AuthorUsername string                    `json:"author_username"`
}

// ContributionEventPushData represents a user's contribution push data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#get-contribution-events-for-a-user
type ContributionEventPushData struct {
	CommitCount int64  `json:"commit_count"`
	Action      string `json:"action"`
	RefType     string `json:"ref_type"`
	CommitFrom  string `json:"commit_from"`
	CommitTo    string `json:"commit_to"`
	Ref         string `json:"ref"`
	CommitTitle string `json:"commit_title"`
}

// ListContributionEventsOptions represents the options for GetUserContributionEvents
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#get-user-contribution-events
type ListContributionEventsOptions struct {
	ListOptions
	Action     *EventTypeValue       `url:"action,omitempty" json:"action,omitempty"`
	TargetType *EventTargetTypeValue `url:"target_type,omitempty" json:"target_type,omitempty"`
	Before     *ISOTime              `url:"before,omitempty" json:"before,omitempty"`
	After      *ISOTime              `url:"after,omitempty" json:"after,omitempty"`
	Sort       *string               `url:"sort,omitempty" json:"sort,omitempty"`
	Scope      *string               `url:"scope,omitempty" json:"scope,omitempty"`
}

func (s *UsersService) ListUserContributionEvents(uid any, opt *ListContributionEventsOptions, options ...RequestOptionFunc) ([]*ContributionEvent, *Response, error) {
	user, err := parseID(uid)
	if err != nil {
		return nil, nil, err
	}

	return do[[]*ContributionEvent](s.client,
		withPath("users/%s/events", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *EventsService) ListCurrentUserContributionEvents(opt *ListContributionEventsOptions, options ...RequestOptionFunc) ([]*ContributionEvent, *Response, error) {
	return do[[]*ContributionEvent](s.client,
		withPath("events"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ProjectEvent represents a GitLab project event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ProjectEvent struct {
	ID             int64                `json:"id"`
	Title          string               `json:"title"`
	ProjectID      int64                `json:"project_id"`
	ActionName     string               `json:"action_name"`
	TargetID       int64                `json:"target_id"`
	TargetIID      int64                `json:"target_iid"`
	TargetType     string               `json:"target_type"`
	AuthorID       int64                `json:"author_id"`
	TargetTitle    string               `json:"target_title"`
	CreatedAt      string               `json:"created_at"`
	Author         BasicUser            `json:"author"`
	AuthorUsername string               `json:"author_username"`
	Data           ProjectEventData     `json:"data"`
	Note           ProjectEventNote     `json:"note"`
	PushData       ProjectEventPushData `json:"push_data"`
}

func (s ProjectEvent) String() string {
	return Stringify(s)
}

// ProjectEventData represents the GitLab project event data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ProjectEventData struct {
	Before            string      `json:"before"`
	After             string      `json:"after"`
	Ref               string      `json:"ref"`
	UserID            int64       `json:"user_id"`
	UserName          string      `json:"user_name"`
	Repository        *Repository `json:"repository"`
	Commits           []*Commit   `json:"commits"`
	TotalCommitsCount int64       `json:"total_commits_count"`
}

func (d ProjectEventData) String() string {
	return Stringify(d)
}

// ProjectEventNote represents a GitLab project event note.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ProjectEventNote struct {
	ID           int64                  `json:"id"`
	Body         string                 `json:"body"`
	Attachment   string                 `json:"attachment"`
	Author       ProjectEventNoteAuthor `json:"author"`
	CreatedAt    *time.Time             `json:"created_at"`
	System       bool                   `json:"system"`
	NoteableID   int64                  `json:"noteable_id"`
	NoteableType string                 `json:"noteable_type"`
	NoteableIID  int64                  `json:"noteable_iid"`
}

func (n ProjectEventNote) String() string {
	return Stringify(n)
}

// ProjectEventNoteAuthor represents a GitLab project event note author.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ProjectEventNoteAuthor struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func (a ProjectEventNoteAuthor) String() string {
	return Stringify(a)
}

// ProjectEventPushData represents a GitLab project event push data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ProjectEventPushData struct {
	CommitCount int64  `json:"commit_count"`
	Action      string `json:"action"`
	RefType     string `json:"ref_type"`
	CommitFrom  string `json:"commit_from"`
	CommitTo    string `json:"commit_to"`
	Ref         string `json:"ref"`
	CommitTitle string `json:"commit_title"`
}

func (d ProjectEventPushData) String() string {
	return Stringify(d)
}

// ListProjectVisibleEventsOptions represents the available
// ListProjectVisibleEvents() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/events/#list-all-visible-events-for-a-project
type ListProjectVisibleEventsOptions struct {
	ListOptions
	Action     *EventTypeValue       `url:"action,omitempty" json:"action,omitempty"`
	TargetType *EventTargetTypeValue `url:"target_type,omitempty" json:"target_type,omitempty"`
	Before     *ISOTime              `url:"before,omitempty" json:"before,omitempty"`
	After      *ISOTime              `url:"after,omitempty" json:"after,omitempty"`
	Sort       *string               `url:"sort,omitempty" json:"sort,omitempty"`
}

func (s *EventsService) ListProjectVisibleEvents(pid any, opt *ListProjectVisibleEventsOptions, options ...RequestOptionFunc) ([]*ProjectEvent, *Response, error) {
	return do[[]*ProjectEvent](s.client,
		withPath("projects/%s/events", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
