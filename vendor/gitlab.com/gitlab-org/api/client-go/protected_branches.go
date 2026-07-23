//
// Copyright 2021, Sander van Harmelen, Michael Lihs
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

import "net/http"

type (
	ProtectedBranchesServiceInterface interface {
		ListProtectedBranches(pid any, opt *ListProtectedBranchesOptions, options ...RequestOptionFunc) ([]*ProtectedBranch, *Response, error)
		GetProtectedBranch(pid any, branch string, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
		ProtectRepositoryBranches(pid any, opt *ProtectRepositoryBranchesOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
		UnprotectRepositoryBranches(pid any, branch string, options ...RequestOptionFunc) (*Response, error)
		UpdateProtectedBranch(pid any, branch string, opt *UpdateProtectedBranchOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
	}

	// ProtectedBranchesService handles communication with the protected branch
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/protected_branches/
	ProtectedBranchesService struct {
		client *Client
	}
)

var _ ProtectedBranchesServiceInterface = (*ProtectedBranchesService)(nil)

// ProtectedBranch represents a protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type ProtectedBranch struct {
	ID                        int64                      `json:"id"`
	Name                      string                     `json:"name"`
	PushAccessLevels          []*BranchAccessDescription `json:"push_access_levels"`
	MergeAccessLevels         []*BranchAccessDescription `json:"merge_access_levels"`
	UnprotectAccessLevels     []*BranchAccessDescription `json:"unprotect_access_levels"`
	AllowForcePush            bool                       `json:"allow_force_push"`
	CodeOwnerApprovalRequired bool                       `json:"code_owner_approval_required"`
}

// BranchAccessDescription represents the access description for a protected
// branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type BranchAccessDescription struct {
	ID                     int64            `json:"id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
	DeployKeyID            int64            `json:"deploy_key_id"`
	UserID                 int64            `json:"user_id"`
	GroupID                int64            `json:"group_id"`
}

// ListProtectedBranchesOptions represents the available ListProtectedBranches()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type ListProtectedBranchesOptions struct {
	ListOptions
	Search *string `url:"search,omitempty" json:"search,omitempty"`
}

// ListProtectedBranches gets a list of protected branches from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
func (s *ProtectedBranchesService) ListProtectedBranches(pid any, opt *ListProtectedBranchesOptions, options ...RequestOptionFunc) ([]*ProtectedBranch, *Response, error) {
	return do[[]*ProtectedBranch](s.client,
		withPath("projects/%s/protected_branches", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProtectedBranch gets a single protected branch or wildcard protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#get-a-single-protected-branch-or-wildcard-protected-branch
func (s *ProtectedBranchesService) GetProtectedBranch(pid any, branch string, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	return do[*ProtectedBranch](s.client,
		withPath("projects/%s/protected_branches/%s", ProjectID{pid}, branch),
		withRequestOpts(options...),
	)
}

// ProtectRepositoryBranchesOptions represents the available
// ProtectRepositoryBranches() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
type ProtectRepositoryBranchesOptions struct {
	Name                      *string                     `url:"name,omitempty" json:"name,omitempty"`
	PushAccessLevel           *AccessLevelValue           `url:"push_access_level,omitempty" json:"push_access_level,omitempty"`
	MergeAccessLevel          *AccessLevelValue           `url:"merge_access_level,omitempty" json:"merge_access_level,omitempty"`
	UnprotectAccessLevel      *AccessLevelValue           `url:"unprotect_access_level,omitempty" json:"unprotect_access_level,omitempty"`
	AllowForcePush            *bool                       `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	AllowedToPush             *[]*BranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*BranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*BranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
	CodeOwnerApprovalRequired *bool                       `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
}

// BranchPermissionOptions represents a branch permission option.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
type BranchPermissionOptions struct {
	ID          *int64            `url:"id,omitempty" json:"id,omitempty"`
	UserID      *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID     *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	DeployKeyID *int64            `url:"deploy_key_id,omitempty" json:"deploy_key_id,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	Destroy     *bool             `url:"_destroy,omitempty" json:"_destroy,omitempty"`
}

// ProtectRepositoryBranches protects a single repository branch or several
// project repository branches using a wildcard protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
func (s *ProtectedBranchesService) ProtectRepositoryBranches(pid any, opt *ProtectRepositoryBranchesOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	return do[*ProtectedBranch](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/protected_branches", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UnprotectRepositoryBranches unprotects the given protected branch or wildcard
// protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#unprotect-repository-branches
func (s *ProtectedBranchesService) UnprotectRepositoryBranches(pid any, branch string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/protected_branches/%s", ProjectID{pid}, branch),
		withRequestOpts(options...),
	)
	return resp, err
}

// UpdateProtectedBranchOptions represents the available
// UpdateProtectedBranch() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
type UpdateProtectedBranchOptions struct {
	Name                      *string                     `url:"name,omitempty" json:"name,omitempty"`
	AllowForcePush            *bool                       `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	CodeOwnerApprovalRequired *bool                       `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
	AllowedToPush             *[]*BranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*BranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*BranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
}

// UpdateProtectedBranch updates a protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
func (s *ProtectedBranchesService) UpdateProtectedBranch(pid any, branch string, opt *UpdateProtectedBranchOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	return do[*ProtectedBranch](s.client,
		withMethod(http.MethodPatch),
		withPath("projects/%s/protected_branches/%s", ProjectID{pid}, branch),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
