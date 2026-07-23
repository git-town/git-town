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
	InvitesServiceInterface interface {
		ListPendingGroupInvitations(gid any, opt *ListPendingInvitationsOptions, options ...RequestOptionFunc) ([]*PendingInvite, *Response, error)
		ListPendingProjectInvitations(pid any, opt *ListPendingInvitationsOptions, options ...RequestOptionFunc) ([]*PendingInvite, *Response, error)
		GroupInvites(gid any, opt *InvitesOptions, options ...RequestOptionFunc) (*InvitesResult, *Response, error)
		ProjectInvites(pid any, opt *InvitesOptions, options ...RequestOptionFunc) (*InvitesResult, *Response, error)
	}

	// InvitesService handles communication with the invitation related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/invitations/
	InvitesService struct {
		client *Client
	}
)

var _ InvitesServiceInterface = (*InvitesService)(nil)

// PendingInvite represents a pending invite.
//
// GitLab API docs: https://docs.gitlab.com/api/invitations/
type PendingInvite struct {
	ID            int64            `json:"id"`
	InviteEmail   string           `json:"invite_email"`
	CreatedAt     *time.Time       `json:"created_at"`
	AccessLevel   AccessLevelValue `json:"access_level"`
	ExpiresAt     *time.Time       `json:"expires_at"`
	UserName      string           `json:"user_name"`
	CreatedByName string           `json:"created_by_name"`
}

// ListPendingInvitationsOptions represents the available
// ListPendingInvitations() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#list-all-invitations-pending-for-a-group-or-project
type ListPendingInvitationsOptions struct {
	ListOptions
	Query *string `url:"query,omitempty" json:"query,omitempty"`
}

// ListPendingGroupInvitations gets a list of invited group members.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#list-all-invitations-pending-for-a-group-or-project
func (s *InvitesService) ListPendingGroupInvitations(gid any, opt *ListPendingInvitationsOptions, options ...RequestOptionFunc) ([]*PendingInvite, *Response, error) {
	return do[[]*PendingInvite](s.client,
		withPath("groups/%s/invitations", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListPendingProjectInvitations gets a list of invited project members.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#list-all-invitations-pending-for-a-group-or-project
func (s *InvitesService) ListPendingProjectInvitations(pid any, opt *ListPendingInvitationsOptions, options ...RequestOptionFunc) ([]*PendingInvite, *Response, error) {
	return do[[]*PendingInvite](s.client,
		withPath("projects/%s/invitations", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// InvitesOptions represents the available GroupInvites() and ProjectInvites()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#add-a-member-to-a-group-or-project
type InvitesOptions struct {
	ID          any               `url:"id,omitempty" json:"id,omitempty"`
	Email       *string           `url:"email,omitempty" json:"email,omitempty"`
	UserID      any               `url:"user_id,omitempty" json:"user_id,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ExpiresAt   *ISOTime          `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// InvitesResult represents an invitations result.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#add-a-member-to-a-group-or-project
type InvitesResult struct {
	Status  string            `json:"status"`
	Message map[string]string `json:"message,omitempty"`
}

// GroupInvites invites new users by email to join a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#add-a-member-to-a-group-or-project
func (s *InvitesService) GroupInvites(gid any, opt *InvitesOptions, options ...RequestOptionFunc) (*InvitesResult, *Response, error) {
	return do[*InvitesResult](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/invitations", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ProjectInvites invites new users by email to join a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/invitations/#add-a-member-to-a-group-or-project
func (s *InvitesService) ProjectInvites(pid any, opt *InvitesOptions, options ...RequestOptionFunc) (*InvitesResult, *Response, error) {
	return do[*InvitesResult](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/invitations", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
