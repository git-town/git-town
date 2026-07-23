package gitlab

import (
	"net/http"
	"net/url"
)

type (
	GroupProtectedBranchesServiceInterface interface {
		// ListProtectedBranches returns a list of protected branches from a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_protected_branches/#list-protected-branches
		ListProtectedBranches(gid any, opt *ListGroupProtectedBranchesOptions, options ...RequestOptionFunc) ([]*GroupProtectedBranch, *Response, error)

		// GetProtectedBranch returns a single group-level protected branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_protected_branches/#get-a-single-protected-branch-or-wildcard-protected-branch
		GetProtectedBranch(gid any, branch string, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error)

		// ProtectRepositoryBranches protects a single group-level branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_protected_branches/#protect-repository-branches
		ProtectRepositoryBranches(gid any, opt *ProtectGroupRepositoryBranchesOptions, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error)

		// UpdateProtectedBranch updates a single group-level protected branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_protected_branches/#update-a-protected-branch
		UpdateProtectedBranch(gid any, branch string, opt *UpdateGroupProtectedBranchOptions, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error)

		// UnprotectRepositoryBranches unprotects the given protected group-level branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_protected_branches/#unprotect-repository-branches
		UnprotectRepositoryBranches(gid any, branch string, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupProtectedBranchesService handles communication with the group-level
	// protected branch methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_protected_branches/
	GroupProtectedBranchesService struct {
		client *Client
	}
)

var _ GroupProtectedBranchesServiceInterface = (*GroupProtectedBranchesService)(nil)

// GroupProtectedBranch represents a group protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#list-protected-branches
type GroupProtectedBranch struct {
	ID                        int64                           `json:"id"`
	Name                      string                          `json:"name"`
	PushAccessLevels          []*GroupBranchAccessDescription `json:"push_access_levels"`
	MergeAccessLevels         []*GroupBranchAccessDescription `json:"merge_access_levels"`
	UnprotectAccessLevels     []*GroupBranchAccessDescription `json:"unprotect_access_levels"`
	AllowForcePush            bool                            `json:"allow_force_push"`
	CodeOwnerApprovalRequired bool                            `json:"code_owner_approval_required"`
}

// GroupBranchAccessDescription represents the access description for a group protected
// branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#list-protected-branches
type GroupBranchAccessDescription struct {
	ID                     int64            `json:"id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
	DeployKeyID            int64            `json:"deploy_key_id"`
	UserID                 int64            `json:"user_id"`
	GroupID                int64            `json:"group_id"`
}

// ListGroupProtectedBranchesOptions represents the available ListProtectedBranches()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#list-protected-branches
type ListGroupProtectedBranchesOptions struct {
	ListOptions
	Search *string `url:"search,omitempty" json:"search,omitempty"`
}

func (s *GroupProtectedBranchesService) ListProtectedBranches(gid any, opt *ListGroupProtectedBranchesOptions, options ...RequestOptionFunc) ([]*GroupProtectedBranch, *Response, error) {
	return do[[]*GroupProtectedBranch](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/protected_branches", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupProtectedBranchesService) GetProtectedBranch(gid any, branch string, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error) {
	return do[*GroupProtectedBranch](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/protected_branches/%s", GroupID{gid}, url.PathEscape(branch)),
		withRequestOpts(options...),
	)
}

// ProtectGroupRepositoryBranchesOptions represents the available
// ProtectRepositoryBranches() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#protect-repository-branches
type ProtectGroupRepositoryBranchesOptions struct {
	Name                      *string                          `url:"name,omitempty" json:"name,omitempty"`
	PushAccessLevel           *AccessLevelValue                `url:"push_access_level,omitempty" json:"push_access_level,omitempty"`
	MergeAccessLevel          *AccessLevelValue                `url:"merge_access_level,omitempty" json:"merge_access_level,omitempty"`
	UnprotectAccessLevel      *AccessLevelValue                `url:"unprotect_access_level,omitempty" json:"unprotect_access_level,omitempty"`
	AllowForcePush            *bool                            `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	AllowedToPush             *[]*GroupBranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*GroupBranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*GroupBranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
	CodeOwnerApprovalRequired *bool                            `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
}

// GroupBranchPermissionOptions represents a branch permission option.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#protect-repository-branches
type GroupBranchPermissionOptions struct {
	ID          *int64            `url:"id,omitempty" json:"id,omitempty"`
	UserID      *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID     *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	DeployKeyID *int64            `url:"deploy_key_id,omitempty" json:"deploy_key_id,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	Destroy     *bool             `url:"_destroy,omitempty" json:"_destroy,omitempty"`
}

func (s *GroupProtectedBranchesService) ProtectRepositoryBranches(gid any, opt *ProtectGroupRepositoryBranchesOptions, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error) {
	return do[*GroupProtectedBranch](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/protected_branches", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateGroupProtectedBranchOptions represents the available
// UpdateProtectedBranch() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_protected_branches/#update-a-protected-branch
type UpdateGroupProtectedBranchOptions struct {
	Name                      *string                          `url:"name,omitempty" json:"name,omitempty"`
	AllowForcePush            *bool                            `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	CodeOwnerApprovalRequired *bool                            `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
	AllowedToPush             *[]*GroupBranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*GroupBranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*GroupBranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
}

func (s *GroupProtectedBranchesService) UpdateProtectedBranch(gid any, branch string, opt *UpdateGroupProtectedBranchOptions, options ...RequestOptionFunc) (*GroupProtectedBranch, *Response, error) {
	return do[*GroupProtectedBranch](s.client,
		withMethod(http.MethodPatch),
		withPath("groups/%s/protected_branches/%s", GroupID{gid}, url.PathEscape(branch)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupProtectedBranchesService) UnprotectRepositoryBranches(gid any, branch string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/protected_branches/%s", GroupID{gid}, url.PathEscape(branch)),
		withRequestOpts(options...),
	)
	return resp, err
}
