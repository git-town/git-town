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
	ProjectMembersServiceInterface interface {
		ListProjectMembers(pid any, opt *ListProjectMembersOptions, options ...RequestOptionFunc) ([]*ProjectMember, *Response, error)
		ListAllProjectMembers(pid any, opt *ListProjectMembersOptions, options ...RequestOptionFunc) ([]*ProjectMember, *Response, error)
		GetProjectMember(pid any, user int64, options ...RequestOptionFunc) (*ProjectMember, *Response, error)
		GetInheritedProjectMember(pid any, user int64, options ...RequestOptionFunc) (*ProjectMember, *Response, error)
		AddProjectMember(pid any, opt *AddProjectMemberOptions, options ...RequestOptionFunc) (*ProjectMember, *Response, error)
		EditProjectMember(pid any, user int64, opt *EditProjectMemberOptions, options ...RequestOptionFunc) (*ProjectMember, *Response, error)
		DeleteProjectMember(pid any, user int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectMembersService handles communication with the project members
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/members/
	ProjectMembersService struct {
		client *Client
	}
)

var _ ProjectMembersServiceInterface = (*ProjectMembersService)(nil)

// ProjectMember represents a project member.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/
type ProjectMember struct {
	ID          int64            `json:"id"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	Name        string           `json:"name"`
	State       string           `json:"state"`
	CreatedAt   *time.Time       `json:"created_at"`
	CreatedBy   *MemberCreatedBy `json:"created_by"`
	ExpiresAt   *ISOTime         `json:"expires_at"`
	AccessLevel AccessLevelValue `json:"access_level"`
	WebURL      string           `json:"web_url"`
	AvatarURL   string           `json:"avatar_url"`
	MemberRole  *MemberRole      `json:"member_role"`
	IsUsingSeat bool             `json:"is_using_seat,omitempty"`
}

type MemberCreatedBy struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

// ListProjectMembersOptions represents the available ListProjectMembers() and
// ListAllProjectMembers() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#list-all-members-of-a-group-or-project
type ListProjectMembersOptions struct {
	ListOptions
	Query        *string  `url:"query,omitempty" json:"query,omitempty"`
	UserIDs      *[]int64 `url:"user_ids[],omitempty" json:"user_ids,omitempty"`
	ShowSeatInfo *bool    `url:"show_seat_info,omitempty" json:"show_seat_info,omitempty"`
}

// ListProjectMembers gets a list of a project's team members viewable by the
// authenticated user. Returns only direct members and not inherited members
// through ancestors groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#list-all-members-of-a-group-or-project
func (s *ProjectMembersService) ListProjectMembers(pid any, opt *ListProjectMembersOptions, options ...RequestOptionFunc) ([]*ProjectMember, *Response, error) {
	return do[[]*ProjectMember](s.client,
		withPath("projects/%s/members", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListAllProjectMembers gets a list of a project's team members viewable by the
// authenticated user. Returns a list including inherited members through
// ancestor groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#list-all-members-of-a-group-or-project-including-inherited-and-invited-members
func (s *ProjectMembersService) ListAllProjectMembers(pid any, opt *ListProjectMembersOptions, options ...RequestOptionFunc) ([]*ProjectMember, *Response, error) {
	return do[[]*ProjectMember](s.client,
		withPath("projects/%s/members/all", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProjectMember gets a project team member.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#get-a-member-of-a-group-or-project
func (s *ProjectMembersService) GetProjectMember(pid any, user int64, options ...RequestOptionFunc) (*ProjectMember, *Response, error) {
	return do[*ProjectMember](s.client,
		withPath("projects/%s/members/%d", ProjectID{pid}, user),
		withRequestOpts(options...),
	)
}

// GetInheritedProjectMember gets a project team member, including inherited
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#get-a-member-of-a-group-or-project-including-inherited-and-invited-members
func (s *ProjectMembersService) GetInheritedProjectMember(pid any, user int64, options ...RequestOptionFunc) (*ProjectMember, *Response, error) {
	return do[*ProjectMember](s.client,
		withPath("projects/%s/members/all/%d", ProjectID{pid}, user),
		withRequestOpts(options...),
	)
}

// AddProjectMemberOptions represents the available AddProjectMember() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#add-a-member-to-a-group-or-project
type AddProjectMemberOptions struct {
	UserID       any               `url:"user_id,omitempty" json:"user_id,omitempty"`
	Username     *string           `url:"username,omitempty" json:"username,omitempty"`
	AccessLevel  *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ExpiresAt    *string           `url:"expires_at,omitempty" json:"expires_at"`
	MemberRoleID *int64            `url:"member_role_id,omitempty" json:"member_role_id,omitempty"`
}

// AddProjectMember adds a user to a project team. This is an idempotent
// method and can be called multiple times with the same parameters. Adding
// team membership to a user that is already a member does not affect the
// existing membership.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#add-a-member-to-a-group-or-project
func (s *ProjectMembersService) AddProjectMember(pid any, opt *AddProjectMemberOptions, options ...RequestOptionFunc) (*ProjectMember, *Response, error) {
	return do[*ProjectMember](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/members", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditProjectMemberOptions represents the available EditProjectMember() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#edit-a-member-of-a-group-or-project
type EditProjectMemberOptions struct {
	AccessLevel  *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ExpiresAt    *string           `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	MemberRoleID *int64            `url:"member_role_id,omitempty" json:"member_role_id,omitempty"`
}

// EditProjectMember updates a project team member to a specified access level..
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#edit-a-member-of-a-group-or-project
func (s *ProjectMembersService) EditProjectMember(pid any, user int64, opt *EditProjectMemberOptions, options ...RequestOptionFunc) (*ProjectMember, *Response, error) {
	return do[*ProjectMember](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/members/%d", ProjectID{pid}, user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteProjectMember removes a user from a project team.
//
// GitLab API docs:
// https://docs.gitlab.com/api/members/#remove-a-member-from-a-group-or-project
func (s *ProjectMembersService) DeleteProjectMember(pid any, user int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/members/%d", ProjectID{pid}, user),
		withRequestOpts(options...),
	)
	return resp, err
}
