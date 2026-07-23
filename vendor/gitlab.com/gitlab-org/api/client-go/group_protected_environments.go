//
// Copyright 2023, Sander van Harmelen
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
)

type (
	GroupProtectedEnvironmentsServiceInterface interface {
		ListGroupProtectedEnvironments(gid any, opt *ListGroupProtectedEnvironmentsOptions, options ...RequestOptionFunc) ([]*GroupProtectedEnvironment, *Response, error)
		GetGroupProtectedEnvironment(gid any, environment string, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error)
		ProtectGroupEnvironment(gid any, opt *ProtectGroupEnvironmentOptions, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error)
		UpdateGroupProtectedEnvironment(gid any, environment string, opt *UpdateGroupProtectedEnvironmentOptions, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error)
		UnprotectGroupEnvironment(gid any, environment string, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupProtectedEnvironmentsService handles communication with the group-level
	// protected environment methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_protected_environments/
	GroupProtectedEnvironmentsService struct {
		client *Client
	}
)

var _ GroupProtectedEnvironmentsServiceInterface = (*GroupProtectedEnvironmentsService)(nil)

// GroupProtectedEnvironment represents a group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/
type GroupProtectedEnvironment struct {
	Name                  string                               `json:"name"`
	DeployAccessLevels    []*GroupEnvironmentAccessDescription `json:"deploy_access_levels"`
	RequiredApprovalCount int64                                `json:"required_approval_count"`
	ApprovalRules         []*GroupEnvironmentApprovalRule      `json:"approval_rules"`
}

// GroupEnvironmentAccessDescription represents the access description for a
// group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/
type GroupEnvironmentAccessDescription struct {
	ID                     int64            `json:"id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
	UserID                 int64            `json:"user_id"`
	GroupID                int64            `json:"group_id"`
	GroupInheritanceType   int64            `json:"group_inheritance_type"`
}

// GroupEnvironmentApprovalRule represents the approval rules for a group-level
// protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#protect-a-single-environment
type GroupEnvironmentApprovalRule struct {
	ID                     int64            `json:"id"`
	UserID                 int64            `json:"user_id"`
	GroupID                int64            `json:"group_id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
	RequiredApprovalCount  int64            `json:"required_approvals"`
	GroupInheritanceType   int64            `json:"group_inheritance_type"`
}

// ListGroupProtectedEnvironmentsOptions represents the available
// ListGroupProtectedEnvironments() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#list-group-level-protected-environments
type ListGroupProtectedEnvironmentsOptions struct {
	ListOptions
}

// ListGroupProtectedEnvironments returns a list of protected environments from
// a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#list-group-level-protected-environments
func (s *GroupProtectedEnvironmentsService) ListGroupProtectedEnvironments(gid any, opt *ListGroupProtectedEnvironmentsOptions, options ...RequestOptionFunc) ([]*GroupProtectedEnvironment, *Response, error) {
	return do[[]*GroupProtectedEnvironment](s.client,
		withPath("groups/%s/protected_environments", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupProtectedEnvironment returns a single group-level protected
// environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#get-a-single-protected-environment
func (s *GroupProtectedEnvironmentsService) GetGroupProtectedEnvironment(gid any, environment string, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error) {
	return do[*GroupProtectedEnvironment](s.client,
		withPath("groups/%s/protected_environments/%s", GroupID{gid}, environment),
		withRequestOpts(options...),
	)
}

// ProtectGroupEnvironmentOptions represents the available
// ProtectGroupEnvironment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#protect-a-single-environment
type ProtectGroupEnvironmentOptions struct {
	Name                  *string                                 `url:"name,omitempty" json:"name,omitempty"`
	DeployAccessLevels    *[]*GroupEnvironmentAccessOptions       `url:"deploy_access_levels,omitempty" json:"deploy_access_levels,omitempty"`
	RequiredApprovalCount *int64                                  `url:"required_approval_count,omitempty" json:"required_approval_count,omitempty"`
	ApprovalRules         *[]*GroupEnvironmentApprovalRuleOptions `url:"approval_rules,omitempty" json:"approval_rules,omitempty"`
}

// GroupEnvironmentAccessOptions represents the options for an access description
// for a group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#protect-a-single-environment
type GroupEnvironmentAccessOptions struct {
	AccessLevel          *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	UserID               *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID              *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	GroupInheritanceType *int64            `url:"group_inheritance_type,omitempty" json:"group_inheritance_type,omitempty"`
}

// GroupEnvironmentApprovalRuleOptions represents the approval rules for a
// group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#protect-a-single-environment
type GroupEnvironmentApprovalRuleOptions struct {
	UserID                 *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID                *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	AccessLevel            *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	AccessLevelDescription *string           `url:"access_level_description,omitempty" json:"access_level_description,omitempty"`
	RequiredApprovalCount  *int64            `url:"required_approvals,omitempty" json:"required_approvals,omitempty"`
	GroupInheritanceType   *int64            `url:"group_inheritance_type,omitempty" json:"group_inheritance_type,omitempty"`
}

// ProtectGroupEnvironment protects a single group-level environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#protect-a-single-environment
func (s *GroupProtectedEnvironmentsService) ProtectGroupEnvironment(gid any, opt *ProtectGroupEnvironmentOptions, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error) {
	return do[*GroupProtectedEnvironment](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/protected_environments", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateGroupProtectedEnvironmentOptions represents the available
// UpdateGroupProtectedEnvironment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#update-a-protected-environment
type UpdateGroupProtectedEnvironmentOptions struct {
	Name                  *string                                       `url:"name,omitempty" json:"name,omitempty"`
	DeployAccessLevels    *[]*UpdateGroupEnvironmentAccessOptions       `url:"deploy_access_levels,omitempty" json:"deploy_access_levels,omitempty"`
	RequiredApprovalCount *int64                                        `url:"required_approval_count,omitempty" json:"required_approval_count,omitempty"`
	ApprovalRules         *[]*UpdateGroupEnvironmentApprovalRuleOptions `url:"approval_rules,omitempty" json:"approval_rules,omitempty"`
}

// UpdateGroupEnvironmentAccessOptions represents the options for updates to the
// access description for a group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#update-a-protected-environment
type UpdateGroupEnvironmentAccessOptions struct {
	AccessLevel          *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ID                   *int64            `url:"id,omitempty" json:"id,omitempty"`
	UserID               *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID              *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	GroupInheritanceType *int64            `url:"group_inheritance_type,omitempty" json:"group_inheritance_type,omitempty"`
	Destroy              *bool             `url:"_destroy,omitempty" json:"_destroy,omitempty"`
}

// UpdateGroupEnvironmentApprovalRuleOptions represents the updates to the
// approval rules for a group-level protected environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#update-a-protected-environment
type UpdateGroupEnvironmentApprovalRuleOptions struct {
	ID                     *int64            `url:"id,omitempty" json:"id,omitempty"`
	UserID                 *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID                *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	AccessLevel            *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	AccessLevelDescription *string           `url:"access_level_description,omitempty" json:"access_level_description,omitempty"`
	RequiredApprovalCount  *int64            `url:"required_approvals,omitempty" json:"required_approvals,omitempty"`
	GroupInheritanceType   *int64            `url:"group_inheritance_type,omitempty" json:"group_inheritance_type,omitempty"`
	Destroy                *bool             `url:"_destroy,omitempty" json:"_destroy,omitempty"`
}

// UpdateGroupProtectedEnvironment updates a single group-level protected
// environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#update-a-protected-environment
func (s *GroupProtectedEnvironmentsService) UpdateGroupProtectedEnvironment(gid any, environment string, opt *UpdateGroupProtectedEnvironmentOptions, options ...RequestOptionFunc) (*GroupProtectedEnvironment, *Response, error) {
	return do[*GroupProtectedEnvironment](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/protected_environments/%s", GroupID{gid}, environment),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UnprotectGroupEnvironment unprotects the given protected group-level
// environment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_environments/#unprotect-a-single-environment
func (s *GroupProtectedEnvironmentsService) UnprotectGroupEnvironment(gid any, environment string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/protected_environments/%s", GroupID{gid}, environment),
		withRequestOpts(options...),
	)
	return resp, err
}
