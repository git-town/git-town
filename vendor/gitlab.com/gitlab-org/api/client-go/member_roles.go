package gitlab

import (
	"net/http"
)

type (
	MemberRolesServiceInterface interface {
		ListInstanceMemberRoles(options ...RequestOptionFunc) ([]*MemberRole, *Response, error)
		CreateInstanceMemberRole(opt *CreateMemberRoleOptions, options ...RequestOptionFunc) (*MemberRole, *Response, error)
		DeleteInstanceMemberRole(memberRoleID int64, options ...RequestOptionFunc) (*Response, error)
		ListMemberRoles(gid any, options ...RequestOptionFunc) ([]*MemberRole, *Response, error)
		CreateMemberRole(gid any, opt *CreateMemberRoleOptions, options ...RequestOptionFunc) (*MemberRole, *Response, error)
		DeleteMemberRole(gid any, memberRole int64, options ...RequestOptionFunc) (*Response, error)
	}

	// MemberRolesService handles communication with the member roles related
	// methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/member_roles/#list-all-member-roles-of-a-group
	MemberRolesService struct {
		client *Client
	}
)

var _ MemberRolesServiceInterface = (*MemberRolesService)(nil)

// MemberRole represents a GitLab member role.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#list-all-member-roles-of-a-group
type MemberRole struct {
	ID                         int64            `json:"id"`
	Name                       string           `json:"name"`
	Description                string           `json:"description,omitempty"`
	GroupID                    int64            `json:"group_id"`
	BaseAccessLevel            AccessLevelValue `json:"base_access_level"`
	AdminCICDVariables         bool             `json:"admin_cicd_variables,omitempty"`
	AdminComplianceFramework   bool             `json:"admin_compliance_framework,omitempty"`
	AdminGroupMembers          bool             `json:"admin_group_member,omitempty"`
	AdminMergeRequests         bool             `json:"admin_merge_request,omitempty"`
	AdminPushRules             bool             `json:"admin_push_rules,omitempty"`
	AdminTerraformState        bool             `json:"admin_terraform_state,omitempty"`
	AdminVulnerability         bool             `json:"admin_vulnerability,omitempty"`
	AdminWebHook               bool             `json:"admin_web_hook,omitempty"`
	ArchiveProject             bool             `json:"archive_project,omitempty"`
	ManageDeployTokens         bool             `json:"manage_deploy_tokens,omitempty"`
	ManageGroupAccessTokens    bool             `json:"manage_group_access_tokens,omitempty"`
	ManageMergeRequestSettings bool             `json:"manage_merge_request_settings,omitempty"`
	ManageProjectAccessTokens  bool             `json:"manage_project_access_tokens,omitempty"`
	ManageSecurityPolicyLink   bool             `json:"manage_security_policy_link,omitempty"`
	ReadCode                   bool             `json:"read_code,omitempty"`
	ReadRunners                bool             `json:"read_runners,omitempty"`
	ReadDependency             bool             `json:"read_dependency,omitempty"`
	ReadVulnerability          bool             `json:"read_vulnerability,omitempty"`
	RemoveGroup                bool             `json:"remove_group,omitempty"`
	RemoveProject              bool             `json:"remove_project,omitempty"`
}

// ListInstanceMemberRoles gets all member roles in an instance.
// Authentication as Administrator is required.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#get-all-instance-member-roles
func (s *MemberRolesService) ListInstanceMemberRoles(options ...RequestOptionFunc) ([]*MemberRole, *Response, error) {
	return do[[]*MemberRole](s.client,
		withMethod(http.MethodGet),
		withPath("member_roles"),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// CreateMemberRoleOptions represents the available CreateInstanceMemberRole()
// and CreateMemberRole() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#create-a-instance-member-role
// https://docs.gitlab.com/api/member_roles/#add-a-member-role-to-a-group
type CreateMemberRoleOptions struct {
	Name                       *string           `url:"name,omitempty" json:"name,omitempty"`
	BaseAccessLevel            *AccessLevelValue `url:"base_access_level,omitempty" json:"base_access_level,omitempty"`
	Description                *string           `url:"description,omitempty" json:"description,omitempty"`
	AdminCICDVariables         *bool             `url:"admin_cicd_variables" json:"admin_cicd_variables,omitempty"`
	AdminComplianceFramework   *bool             `url:"admin_compliance_framework" json:"admin_compliance_framework,omitempty"`
	AdminGroupMembers          *bool             `url:"admin_group_member" json:"admin_group_member,omitempty"`
	AdminMergeRequest          *bool             `url:"admin_merge_request,omitempty" json:"admin_merge_request,omitempty"`
	AdminPushRules             *bool             `url:"admin_push_rules" json:"admin_push_rules,omitempty"`
	AdminTerraformState        *bool             `url:"admin_terraform_state" json:"admin_terraform_state,omitempty"`
	AdminVulnerability         *bool             `url:"admin_vulnerability,omitempty" json:"admin_vulnerability,omitempty"`
	AdminWebHook               *bool             `url:"admin_web_hook" json:"admin_web_hook,omitempty"`
	ArchiveProject             *bool             `url:"archive_project" json:"archive_project,omitempty"`
	ManageDeployTokens         *bool             `url:"manage_deploy_tokens" json:"manage_deploy_tokens,omitempty"`
	ManageGroupAccessTokens    *bool             `url:"manage_group_access_tokens" json:"manage_group_access_tokens,omitempty"`
	ManageMergeRequestSettings *bool             `url:"manage_merge_request_settings" json:"manage_merge_request_settings,omitempty"`
	ManageProjectAccessTokens  *bool             `url:"manage_project_access_tokens" json:"manage_project_access_tokens,omitempty"`
	ManageSecurityPolicyLink   *bool             `url:"manage_security_policy_link" json:"manage_security_policy_link,omitempty"`
	ReadCode                   *bool             `url:"read_code,omitempty" json:"read_code,omitempty"`
	ReadRunners                *bool             `url:"read_runners" json:"read_runners,omitempty"`
	ReadDependency             *bool             `url:"read_dependency,omitempty" json:"read_dependency,omitempty"`
	ReadVulnerability          *bool             `url:"read_vulnerability,omitempty" json:"read_vulnerability,omitempty"`
	RemoveGroup                *bool             `url:"remove_group" json:"remove_group,omitempty"`
	RemoveProject              *bool             `url:"remove_project" json:"remove_project,omitempty"`
}

// CreateInstanceMemberRole creates an instance-wide member role.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#create-a-instance-member-role
func (s *MemberRolesService) CreateInstanceMemberRole(opt *CreateMemberRoleOptions, options ...RequestOptionFunc) (*MemberRole, *Response, error) {
	return do[*MemberRole](s.client,
		withMethod(http.MethodPost),
		withPath("member_roles"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteInstanceMemberRole deletes a member role from a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#delete-an-instance-member-role
func (s *MemberRolesService) DeleteInstanceMemberRole(memberRoleID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("member_roles/%d", memberRoleID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListMemberRoles gets a list of member roles for a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#get-all-group-member-roles
func (s *MemberRolesService) ListMemberRoles(gid any, options ...RequestOptionFunc) ([]*MemberRole, *Response, error) {
	return do[[]*MemberRole](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/member_roles", GroupID{gid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// CreateMemberRole creates a new member role for a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#add-a-member-role-to-a-group
func (s *MemberRolesService) CreateMemberRole(gid any, opt *CreateMemberRoleOptions, options ...RequestOptionFunc) (*MemberRole, *Response, error) {
	return do[*MemberRole](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/member_roles", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteMemberRole deletes a member role from a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/member_roles/#remove-member-role-of-a-group
func (s *MemberRolesService) DeleteMemberRole(gid any, memberRole int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/member_roles/%d", GroupID{gid}, memberRole),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}
