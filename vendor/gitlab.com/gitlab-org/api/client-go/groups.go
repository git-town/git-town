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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type (
	GroupsServiceInterface interface {
		ListGroups(opt *ListGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error)
		ListSubGroups(gid any, opt *ListSubGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error)
		ListDescendantGroups(gid any, opt *ListDescendantGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error)
		ListGroupProjects(gid any, opt *ListGroupProjectsOptions, options ...RequestOptionFunc) ([]*Project, *Response, error)
		GetGroup(gid any, opt *GetGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error)
		DownloadAvatar(gid any, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		CreateGroup(opt *CreateGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error)
		TransferGroup(gid any, pid any, options ...RequestOptionFunc) (*Group, *Response, error)
		TransferSubGroup(gid any, opt *TransferSubGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error)
		UpdateGroup(gid any, opt *UpdateGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error)
		UploadAvatar(gid any, avatar io.Reader, filename string, options ...RequestOptionFunc) (*Group, *Response, error)
		DeleteGroup(gid any, opt *DeleteGroupOptions, options ...RequestOptionFunc) (*Response, error)
		RestoreGroup(gid any, options ...RequestOptionFunc) (*Group, *Response, error)
		SearchGroup(query string, options ...RequestOptionFunc) ([]*Group, *Response, error)
		ListProvisionedUsers(gid any, opt *ListProvisionedUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
		ListGroupLDAPLinks(gid any, options ...RequestOptionFunc) ([]*LDAPGroupLink, *Response, error)
		AddGroupLDAPLink(gid any, opt *AddGroupLDAPLinkOptions, options ...RequestOptionFunc) (*LDAPGroupLink, *Response, error)
		DeleteGroupLDAPLink(gid any, cn string, options ...RequestOptionFunc) (*Response, error)
		DeleteGroupLDAPLinkWithCNOrFilter(gid any, opts *DeleteGroupLDAPLinkWithCNOrFilterOptions, options ...RequestOptionFunc) (*Response, error)
		DeleteGroupLDAPLinkForProvider(gid any, provider, cn string, options ...RequestOptionFunc) (*Response, error)
		ListGroupSAMLLinks(gid any, options ...RequestOptionFunc) ([]*SAMLGroupLink, *Response, error)
		ListGroupSharedProjects(gid any, opt *ListGroupSharedProjectsOptions, options ...RequestOptionFunc) ([]*Project, *Response, error)
		GetGroupSAMLLink(gid any, samlGroupName string, options ...RequestOptionFunc) (*SAMLGroupLink, *Response, error)
		AddGroupSAMLLink(gid any, opt *AddGroupSAMLLinkOptions, options ...RequestOptionFunc) (*SAMLGroupLink, *Response, error)
		DeleteGroupSAMLLink(gid any, samlGroupName string, options ...RequestOptionFunc) (*Response, error)
		ShareGroupWithGroup(gid any, opt *ShareGroupWithGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error)
		UnshareGroupFromGroup(gid any, groupID int, options ...RequestOptionFunc) (*Response, error)
		GetGroupPushRules(gid any, options ...RequestOptionFunc) (*GroupPushRules, *Response, error)
		AddGroupPushRule(gid any, opt *AddGroupPushRuleOptions, options ...RequestOptionFunc) (*GroupPushRules, *Response, error)
		EditGroupPushRule(gid any, opt *EditGroupPushRuleOptions, options ...RequestOptionFunc) (*GroupPushRules, *Response, error)
		DeleteGroupPushRule(gid any, options ...RequestOptionFunc) (*Response, error)

		// group_hooks.go
		ListGroupHooks(gid any, opt *ListGroupHooksOptions, options ...RequestOptionFunc) ([]*GroupHook, *Response, error)
		GetGroupHook(gid any, hook int, options ...RequestOptionFunc) (*GroupHook, *Response, error)
		ResendGroupHookEvent(gid any, hook int, hookEventID int, options ...RequestOptionFunc) (*Response, error)
		AddGroupHook(gid any, opt *AddGroupHookOptions, options ...RequestOptionFunc) (*GroupHook, *Response, error)
		EditGroupHook(gid any, hook int, opt *EditGroupHookOptions, options ...RequestOptionFunc) (*GroupHook, *Response, error)
		DeleteGroupHook(gid any, hook int, options ...RequestOptionFunc) (*Response, error)
		TriggerTestGroupHook(pid any, hook int, trigger GroupHookTrigger, options ...RequestOptionFunc) (*Response, error)
		SetGroupCustomHeader(gid any, hook int, key string, opt *SetHookCustomHeaderOptions, options ...RequestOptionFunc) (*Response, error)
		DeleteGroupCustomHeader(gid any, hook int, key string, options ...RequestOptionFunc) (*Response, error)
		SetGroupHookURLVariable(gid any, hook int, key string, opt *SetHookURLVariableOptions, options ...RequestOptionFunc) (*Response, error)
		DeleteGroupHookURLVariable(gid any, hook int, key string, options ...RequestOptionFunc) (*Response, error)

		// group_serviceaccounts.go
		ListServiceAccounts(gid any, opt *ListServiceAccountsOptions, options ...RequestOptionFunc) ([]*GroupServiceAccount, *Response, error)
		CreateServiceAccount(gid any, opt *CreateServiceAccountOptions, options ...RequestOptionFunc) (*GroupServiceAccount, *Response, error)
		DeleteServiceAccount(gid any, serviceAccount int, opt *DeleteServiceAccountOptions, options ...RequestOptionFunc) (*Response, error)
		CreateServiceAccountPersonalAccessToken(gid any, serviceAccount int, opt *CreateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		RotateServiceAccountPersonalAccessToken(gid any, serviceAccount, token int, opt *RotateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)

		// group_members.go
		ListGroupMembers(gid any, opt *ListGroupMembersOptions, options ...RequestOptionFunc) ([]*GroupMember, *Response, error)
		ListAllGroupMembers(gid any, opt *ListGroupMembersOptions, options ...RequestOptionFunc) ([]*GroupMember, *Response, error)
		ListBillableGroupMembers(gid any, opt *ListBillableGroupMembersOptions, options ...RequestOptionFunc) ([]*BillableGroupMember, *Response, error)
		ListMembershipsForBillableGroupMember(gid any, user int, opt *ListMembershipsForBillableGroupMemberOptions, options ...RequestOptionFunc) ([]*BillableUserMembership, *Response, error)
		RemoveBillableGroupMember(gid any, user int, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupsService handles communication with the group related methods of
	// the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/groups/
	GroupsService struct {
		client *Client
	}
)

var _ GroupsServiceInterface = (*GroupsService)(nil)

// Group represents a GitLab group.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/
type Group struct {
	ID                              int                        `json:"id"`
	Name                            string                     `json:"name"`
	Path                            string                     `json:"path"`
	Description                     string                     `json:"description"`
	MembershipLock                  bool                       `json:"membership_lock"`
	Visibility                      VisibilityValue            `json:"visibility"`
	LFSEnabled                      bool                       `json:"lfs_enabled"`
	DefaultBranch                   string                     `json:"default_branch"`
	DefaultBranchProtectionDefaults *BranchProtectionDefaults  `json:"default_branch_protection_defaults"`
	AvatarURL                       string                     `json:"avatar_url"`
	WebURL                          string                     `json:"web_url"`
	RequestAccessEnabled            bool                       `json:"request_access_enabled"`
	RepositoryStorage               string                     `json:"repository_storage"`
	FullName                        string                     `json:"full_name"`
	FullPath                        string                     `json:"full_path"`
	FileTemplateProjectID           int                        `json:"file_template_project_id"`
	ParentID                        int                        `json:"parent_id"`
	Statistics                      *Statistics                `json:"statistics"`
	CustomAttributes                []*CustomAttribute         `json:"custom_attributes"`
	ShareWithGroupLock              bool                       `json:"share_with_group_lock"`
	RequireTwoFactorAuth            bool                       `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod            int                        `json:"two_factor_grace_period"`
	ProjectCreationLevel            ProjectCreationLevelValue  `json:"project_creation_level"`
	AutoDevopsEnabled               bool                       `json:"auto_devops_enabled"`
	SubGroupCreationLevel           SubGroupCreationLevelValue `json:"subgroup_creation_level"`
	EmailsEnabled                   bool                       `json:"emails_enabled"`
	MentionsDisabled                bool                       `json:"mentions_disabled"`
	RunnersToken                    string                     `json:"runners_token"`
	SharedRunnersSetting            SharedRunnersSettingValue  `json:"shared_runners_setting"`
	SharedWithGroups                []SharedWithGroup          `json:"shared_with_groups"`
	LDAPCN                          string                     `json:"ldap_cn"`
	LDAPAccess                      AccessLevelValue           `json:"ldap_access"`
	LDAPGroupLinks                  []*LDAPGroupLink           `json:"ldap_group_links"`
	SAMLGroupLinks                  []*SAMLGroupLink           `json:"saml_group_links"`
	SharedRunnersMinutesLimit       int                        `json:"shared_runners_minutes_limit"`
	ExtraSharedRunnersMinutesLimit  int                        `json:"extra_shared_runners_minutes_limit"`
	PreventForkingOutsideGroup      bool                       `json:"prevent_forking_outside_group"`
	MarkedForDeletionOn             *ISOTime                   `json:"marked_for_deletion_on"`
	CreatedAt                       *time.Time                 `json:"created_at"`
	IPRestrictionRanges             string                     `json:"ip_restriction_ranges"`
	AllowedEmailDomainsList         string                     `json:"allowed_email_domains_list"`
	WikiAccessLevel                 AccessControlValue         `json:"wiki_access_level"`

	// Deprecated: will be removed in v5 of the API, use ListGroupProjects instead
	Projects []*Project `json:"projects"`

	// Deprecated: will be removed in v5 of the API, use ListGroupSharedProjects instead
	SharedProjects []*Project `json:"shared_projects"`

	// Deprecated: Use EmailsEnabled instead
	EmailsDisabled bool `json:"emails_disabled"`

	// Deprecated: Use DefaultBranchProtectionDefaults instead
	DefaultBranchProtection int `json:"default_branch_protection"`
}

// SharedWithGroup represents a GitLab group shared with a group.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/
type SharedWithGroup struct {
	GroupID          int      `json:"group_id"`
	GroupName        string   `json:"group_name"`
	GroupFullPath    string   `json:"group_full_path"`
	GroupAccessLevel int      `json:"group_access_level"`
	ExpiresAt        *ISOTime `json:"expires_at"`
	MemberRoleID     int      `json:"member_role_id"`
}

// BranchProtectionDefaults represents default Git protected branch permissions.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#options-for-default_branch_protection_defaults
type BranchProtectionDefaults struct {
	AllowedToPush           []*GroupAccessLevel `json:"allowed_to_push,omitempty"`
	AllowForcePush          bool                `json:"allow_force_push,omitempty"`
	AllowedToMerge          []*GroupAccessLevel `json:"allowed_to_merge,omitempty"`
	DeveloperCanInitialPush bool                `json:"developer_can_initial_push,omitempty"`
}

// GroupAccessLevel represents default branch protection defaults access levels.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#options-for-default_branch_protection_defaults
type GroupAccessLevel struct {
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
}

// GroupAvatar represents a GitLab group avatar.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/
type GroupAvatar struct {
	Filename string
	Image    io.Reader
}

// MarshalJSON implements the json.Marshaler interface.
func (a *GroupAvatar) MarshalJSON() ([]byte, error) {
	if a.Filename == "" && a.Image == nil {
		return []byte(`""`), nil
	}
	type alias GroupAvatar
	return json.Marshal((*alias)(a))
}

// LDAPGroupLink represents a GitLab LDAP group link.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/
type LDAPGroupLink struct {
	CN           string           `json:"cn"`
	Filter       string           `json:"filter"`
	GroupAccess  AccessLevelValue `json:"group_access"`
	Provider     string           `json:"provider"`
	MemberRoleID int64            `json:"member_role_id"`
}

// SAMLGroupLink represents a GitLab SAML group link.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#saml-group-links
type SAMLGroupLink struct {
	Name         string           `json:"name"`
	AccessLevel  AccessLevelValue `json:"access_level"`
	MemberRoleID int              `json:"member_role_id,omitempty"`
}

// ListGroupsOptions represents the available ListGroups() options.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#list-groups
type ListGroupsOptions struct {
	ListOptions
	SkipGroups           *[]int            `url:"skip_groups,omitempty" del:"," json:"skip_groups,omitempty"`
	AllAvailable         *bool             `url:"all_available,omitempty" json:"all_available,omitempty"`
	Search               *string           `url:"search,omitempty" json:"search,omitempty"`
	OrderBy              *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                 *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Statistics           *bool             `url:"statistics,omitempty" json:"statistics,omitempty"`
	WithCustomAttributes *bool             `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
	Owned                *bool             `url:"owned,omitempty" json:"owned,omitempty"`
	MinAccessLevel       *AccessLevelValue `url:"min_access_level,omitempty" json:"min_access_level,omitempty"`
	TopLevelOnly         *bool             `url:"top_level_only,omitempty" json:"top_level_only,omitempty"`
	RepositoryStorage    *string           `url:"repository_storage,omitempty" json:"repository_storage,omitempty"`
}

// ListGroups gets a list of groups (as user: my groups, as admin: all groups).
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-groups
func (s *GroupsService) ListGroups(opt *ListGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "groups", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var gs []*Group
	resp, err := s.client.Do(req, &gs)
	if err != nil {
		return nil, resp, err
	}

	return gs, resp, nil
}

// ListSubGroupsOptions represents the available ListSubGroups() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-subgroups
type ListSubGroupsOptions ListGroupsOptions

// ListSubGroups gets a list of subgroups for a given group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-subgroups
func (s *GroupsService) ListSubGroups(gid any, opt *ListSubGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/subgroups", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var gs []*Group
	resp, err := s.client.Do(req, &gs)
	if err != nil {
		return nil, resp, err
	}

	return gs, resp, nil
}

// ListDescendantGroupsOptions represents the available ListDescendantGroups()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-descendant-groups
type ListDescendantGroupsOptions ListGroupsOptions

// ListDescendantGroups gets a list of subgroups for a given project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-descendant-groups
func (s *GroupsService) ListDescendantGroups(gid any, opt *ListDescendantGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/descendant_groups", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var gs []*Group
	resp, err := s.client.Do(req, &gs)
	if err != nil {
		return nil, resp, err
	}

	return gs, resp, nil
}

// ListGroupProjectsOptions represents the available ListGroup() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-projects
type ListGroupProjectsOptions struct {
	ListOptions
	Archived                 *bool             `url:"archived,omitempty" json:"archived,omitempty"`
	IncludeSubGroups         *bool             `url:"include_subgroups,omitempty" json:"include_subgroups,omitempty"`
	MinAccessLevel           *AccessLevelValue `url:"min_access_level,omitempty" json:"min_access_level,omitempty"`
	OrderBy                  *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Owned                    *bool             `url:"owned,omitempty" json:"owned,omitempty"`
	Search                   *string           `url:"search,omitempty" json:"search,omitempty"`
	Simple                   *bool             `url:"simple,omitempty" json:"simple,omitempty"`
	Sort                     *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Starred                  *bool             `url:"starred,omitempty" json:"starred,omitempty"`
	Topic                    *string           `url:"topic,omitempty" json:"topic,omitempty"`
	Visibility               *VisibilityValue  `url:"visibility,omitempty" json:"visibility,omitempty"`
	WithCustomAttributes     *bool             `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
	WithIssuesEnabled        *bool             `url:"with_issues_enabled,omitempty" json:"with_issues_enabled,omitempty"`
	WithMergeRequestsEnabled *bool             `url:"with_merge_requests_enabled,omitempty" json:"with_merge_requests_enabled,omitempty"`
	WithSecurityReports      *bool             `url:"with_security_reports,omitempty" json:"with_security_reports,omitempty"`
	WithShared               *bool             `url:"with_shared,omitempty" json:"with_shared,omitempty"`
}

// ListGroupProjects get a list of group projects
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-projects
func (s *GroupsService) ListGroupProjects(gid any, opt *ListGroupProjectsOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/projects", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ps []*Project
	resp, err := s.client.Do(req, &ps)
	if err != nil {
		return nil, resp, err
	}

	return ps, resp, nil
}

// GetGroupOptions represents the available GetGroup() options.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#get-a-single-group
type GetGroupOptions struct {
	ListOptions
	WithCustomAttributes *bool `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`

	// Deprecated: will be removed in v5 of the API, use ListGroupProjects instead
	WithProjects *bool `url:"with_projects,omitempty" json:"with_projects,omitempty"`
}

// GetGroup gets all details of a group.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#get-a-single-group
func (s *GroupsService) GetGroup(gid any, opt *GetGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// DownloadAvatar downloads a group avatar.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#download-a-group-avatar
func (s *GroupsService) DownloadAvatar(gid any, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/avatar", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	avatar := new(bytes.Buffer)
	resp, err := s.client.Do(req, avatar)
	if err != nil {
		return nil, resp, err
	}

	return bytes.NewReader(avatar.Bytes()), resp, err
}

// CreateGroupOptions represents the available CreateGroup() options.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#create-a-group
type CreateGroupOptions struct {
	Name                            *string                                 `url:"name,omitempty" json:"name,omitempty"`
	Path                            *string                                 `url:"path,omitempty" json:"path,omitempty"`
	Avatar                          *GroupAvatar                            `url:"-" json:"-"`
	DefaultBranch                   *string                                 `url:"default_branch,omitempty" json:"default_branch,omitempty"`
	Description                     *string                                 `url:"description,omitempty" json:"description,omitempty"`
	MembershipLock                  *bool                                   `url:"membership_lock,omitempty" json:"membership_lock,omitempty"`
	Visibility                      *VisibilityValue                        `url:"visibility,omitempty" json:"visibility,omitempty"`
	ShareWithGroupLock              *bool                                   `url:"share_with_group_lock,omitempty" json:"share_with_group_lock,omitempty"`
	RequireTwoFactorAuth            *bool                                   `url:"require_two_factor_authentication,omitempty" json:"require_two_factor_authentication,omitempty"`
	TwoFactorGracePeriod            *int                                    `url:"two_factor_grace_period,omitempty" json:"two_factor_grace_period,omitempty"`
	ProjectCreationLevel            *ProjectCreationLevelValue              `url:"project_creation_level,omitempty" json:"project_creation_level,omitempty"`
	AutoDevopsEnabled               *bool                                   `url:"auto_devops_enabled,omitempty" json:"auto_devops_enabled,omitempty"`
	SubGroupCreationLevel           *SubGroupCreationLevelValue             `url:"subgroup_creation_level,omitempty" json:"subgroup_creation_level,omitempty"`
	EmailsEnabled                   *bool                                   `url:"emails_enabled,omitempty" json:"emails_enabled,omitempty"`
	MentionsDisabled                *bool                                   `url:"mentions_disabled,omitempty" json:"mentions_disabled,omitempty"`
	LFSEnabled                      *bool                                   `url:"lfs_enabled,omitempty" json:"lfs_enabled,omitempty"`
	DefaultBranchProtectionDefaults *DefaultBranchProtectionDefaultsOptions `url:"default_branch_protection_defaults,omitempty" json:"default_branch_protection_defaults,omitempty"`
	RequestAccessEnabled            *bool                                   `url:"request_access_enabled,omitempty" json:"request_access_enabled,omitempty"`
	ParentID                        *int                                    `url:"parent_id,omitempty" json:"parent_id,omitempty"`
	SharedRunnersMinutesLimit       *int                                    `url:"shared_runners_minutes_limit,omitempty" json:"shared_runners_minutes_limit,omitempty"`
	ExtraSharedRunnersMinutesLimit  *int                                    `url:"extra_shared_runners_minutes_limit,omitempty" json:"extra_shared_runners_minutes_limit,omitempty"`
	WikiAccessLevel                 *AccessControlValue                     `url:"wiki_access_level,omitempty" json:"wiki_access_level,omitempty"`

	// Deprecated: Use EmailsEnabled instead
	EmailsDisabled *bool `url:"emails_disabled,omitempty" json:"emails_disabled,omitempty"`

	// Deprecated: User DefaultBranchProtectionDefaults instead
	DefaultBranchProtection *int `url:"default_branch_protection,omitempty" json:"default_branch_protection,omitempty"`
}

// DefaultBranchProtectionDefaultsOptions represents the available options for
// using default_branch_protection_defaults in CreateGroup() or UpdateGroup()
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#options-for-default_branch_protection_defaults
type DefaultBranchProtectionDefaultsOptions struct {
	AllowedToPush           *[]*GroupAccessLevel `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowForcePush          *bool                `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	AllowedToMerge          *[]*GroupAccessLevel `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	DeveloperCanInitialPush *bool                `url:"developer_can_initial_push,omitempty" json:"developer_can_initial_push,omitempty"`
}

// EncodeValues implements the query.Encoder interface
func (d *DefaultBranchProtectionDefaultsOptions) EncodeValues(key string, v *url.Values) error {
	if d.AllowForcePush != nil {
		v.Add(key+"[allow_force_push]", strconv.FormatBool(*d.AllowForcePush))
	}
	if d.DeveloperCanInitialPush != nil {
		v.Add(key+"[developer_can_initial_push]", strconv.FormatBool(*d.DeveloperCanInitialPush))
	}
	// The GitLab API only accepts one value for `allowed_to_merge` even when multiples are
	// provided on the request.  The API will take the highest permission level.  For instance,
	// if 'developer' and 'maintainer' are provided, the API will take 'maintainer'.
	if d.AllowedToMerge != nil {
		for _, atm := range *d.AllowedToMerge {
			if atm != nil {
				v.Add(key+"[allowed_to_merge][][access_level]", strconv.Itoa((int)(*atm.AccessLevel)))
			}
		}
	}
	// The GitLab API only accepts one value for `allowed_to_push` even when multiples are
	// provided on the request.  The API will take the highest permission level.  For instance,
	// if 'developer' and 'maintainer' are provided, the API will take 'maintainer'.
	if d.AllowedToPush != nil {
		for _, atp := range *d.AllowedToPush {
			if atp != nil {
				v.Add(key+"[allowed_to_push][][access_level]", strconv.Itoa((int)(*atp.AccessLevel)))
			}
		}
	}
	return nil
}

// CreateGroup creates a new project group. Available only for users who can
// create groups.
//
// When `default_branch_protection_defaults` are defined with an `avatar` value,
// only one value for `allowed_to_push` and `allowed_to_merge` will be used as
// the GitLab API only accepts one value for those attributes even when multiples
// are provided on the request. The API will take the highest permission level.
// For instance, if 'developer' and 'maintainer' are provided, the API will take 'maintainer'.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#create-a-group
func (s *GroupsService) CreateGroup(opt *CreateGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error) {
	var err error
	var req *retryablehttp.Request

	if opt.Avatar == nil {
		req, err = s.client.NewRequest(http.MethodPost, "groups", opt, options)
	} else {
		// since the Avatar is provided, check allowed_to_push and
		// allowed_to_merge access levels and error if multiples are provided
		if opt.DefaultBranchProtectionDefaults != nil && (len(*opt.DefaultBranchProtectionDefaults.AllowedToMerge) > 1 || len(*opt.DefaultBranchProtectionDefaults.AllowedToPush) > 1) {
			return nil, nil, fmt.Errorf("multiple access levels for allowed_to_merge or allowed_to_push are not permitted when an Avatar is also specified as it will result in unexpected behavior")
		}
		req, err = s.client.UploadRequest(
			http.MethodPost,
			"groups",
			opt.Avatar.Image,
			opt.Avatar.Filename,
			UploadAvatar,
			opt,
			options,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// TransferGroup transfers a project to the Group namespace. Available only
// for admin.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#transfer-a-project-to-a-group
func (s *GroupsService) TransferGroup(gid any, pid any, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/projects/%s", PathEscape(group), PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// TransferSubGroupOptions represents the available TransferSubGroup() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#transfer-a-group
type TransferSubGroupOptions struct {
	GroupID *int `url:"group_id,omitempty" json:"group_id,omitempty"`
}

// TransferSubGroup transfers a group to a new parent group or turn a subgroup
// to a top-level group. Available to administrators and users.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#transfer-a-group
func (s *GroupsService) TransferSubGroup(gid any, opt *TransferSubGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/transfer", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// UpdateGroupOptions represents the available UpdateGroup() options.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#update-group-attributes
type UpdateGroupOptions struct {
	Name                                 *string                                 `url:"name,omitempty" json:"name,omitempty"`
	Path                                 *string                                 `url:"path,omitempty" json:"path,omitempty"`
	Avatar                               *GroupAvatar                            `url:"-" json:"avatar,omitempty"`
	DefaultBranch                        *string                                 `url:"default_branch,omitempty" json:"default_branch,omitempty"`
	Description                          *string                                 `url:"description,omitempty" json:"description,omitempty"`
	MembershipLock                       *bool                                   `url:"membership_lock,omitempty" json:"membership_lock,omitempty"`
	Visibility                           *VisibilityValue                        `url:"visibility,omitempty" json:"visibility,omitempty"`
	ShareWithGroupLock                   *bool                                   `url:"share_with_group_lock,omitempty" json:"share_with_group_lock,omitempty"`
	RequireTwoFactorAuth                 *bool                                   `url:"require_two_factor_authentication,omitempty" json:"require_two_factor_authentication,omitempty"`
	TwoFactorGracePeriod                 *int                                    `url:"two_factor_grace_period,omitempty" json:"two_factor_grace_period,omitempty"`
	ProjectCreationLevel                 *ProjectCreationLevelValue              `url:"project_creation_level,omitempty" json:"project_creation_level,omitempty"`
	AutoDevopsEnabled                    *bool                                   `url:"auto_devops_enabled,omitempty" json:"auto_devops_enabled,omitempty"`
	SubGroupCreationLevel                *SubGroupCreationLevelValue             `url:"subgroup_creation_level,omitempty" json:"subgroup_creation_level,omitempty"`
	EmailsEnabled                        *bool                                   `url:"emails_enabled,omitempty" json:"emails_enabled,omitempty"`
	MentionsDisabled                     *bool                                   `url:"mentions_disabled,omitempty" json:"mentions_disabled,omitempty"`
	LFSEnabled                           *bool                                   `url:"lfs_enabled,omitempty" json:"lfs_enabled,omitempty"`
	RequestAccessEnabled                 *bool                                   `url:"request_access_enabled,omitempty" json:"request_access_enabled,omitempty"`
	DefaultBranchProtectionDefaults      *DefaultBranchProtectionDefaultsOptions `url:"default_branch_protection_defaults,omitempty" json:"default_branch_protection_defaults,omitempty"`
	FileTemplateProjectID                *int                                    `url:"file_template_project_id,omitempty" json:"file_template_project_id,omitempty"`
	SharedRunnersMinutesLimit            *int                                    `url:"shared_runners_minutes_limit,omitempty" json:"shared_runners_minutes_limit,omitempty"`
	ExtraSharedRunnersMinutesLimit       *int                                    `url:"extra_shared_runners_minutes_limit,omitempty" json:"extra_shared_runners_minutes_limit,omitempty"`
	PreventForkingOutsideGroup           *bool                                   `url:"prevent_forking_outside_group,omitempty" json:"prevent_forking_outside_group,omitempty"`
	SharedRunnersSetting                 *SharedRunnersSettingValue              `url:"shared_runners_setting,omitempty" json:"shared_runners_setting,omitempty"`
	PreventSharingGroupsOutsideHierarchy *bool                                   `url:"prevent_sharing_groups_outside_hierarchy,omitempty" json:"prevent_sharing_groups_outside_hierarchy,omitempty"`
	IPRestrictionRanges                  *string                                 `url:"ip_restriction_ranges,omitempty" json:"ip_restriction_ranges,omitempty"`
	AllowedEmailDomainsList              *string                                 `url:"allowed_email_domains_list,omitempty" json:"allowed_email_domains_list,omitempty"`
	WikiAccessLevel                      *AccessControlValue                     `url:"wiki_access_level,omitempty" json:"wiki_access_level,omitempty"`

	// Deprecated: Use EmailsEnabled instead
	EmailsDisabled *bool `url:"emails_disabled,omitempty" json:"emails_disabled,omitempty"`

	// Deprecated: Use DefaultBranchProtectionDefaults instead
	DefaultBranchProtection *int `url:"default_branch_protection,omitempty" json:"default_branch_protection,omitempty"`
}

// UpdateGroup updates an existing group; only available to group owners and
// administrators.
//
// When `default_branch_protection_defaults` are defined with an `avatar` value,
// only one value for `allowed_to_push` and `allowed_to_merge` will be used as
// the GitLab API only accepts one value for those attributes even when multiples
// are provided on the request. The API will take the highest permission level.
// For instance, if 'developer' and 'maintainer' are provided, the API will take 'maintainer'.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#update-group-attributes
func (s *GroupsService) UpdateGroup(gid any, opt *UpdateGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s", PathEscape(group))

	var req *retryablehttp.Request

	if opt.Avatar == nil || (opt.Avatar.Filename == "" && opt.Avatar.Image == nil) {
		req, err = s.client.NewRequest(http.MethodPut, u, opt, options)
	} else {
		// since the Avatar is provided, check allowed_to_push and
		// allowed_to_merge access levels and error if multiples are provided
		if opt.DefaultBranchProtectionDefaults != nil && (len(*opt.DefaultBranchProtectionDefaults.AllowedToMerge) > 1 || len(*opt.DefaultBranchProtectionDefaults.AllowedToPush) > 1) {
			return nil, nil, fmt.Errorf("multiple access levels for allowed_to_merge or allowed_to_push are not permitted when an Avatar is also specified as it will result in unexpected behavior")
		}
		req, err = s.client.UploadRequest(
			http.MethodPut,
			u,
			opt.Avatar.Image,
			opt.Avatar.Filename,
			UploadAvatar,
			opt,
			options,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// UploadAvatar uploads a group avatar.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#upload-a-group-avatar
func (s *GroupsService) UploadAvatar(gid any, avatar io.Reader, filename string, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s", PathEscape(group))

	req, err := s.client.UploadRequest(
		http.MethodPut,
		u,
		avatar,
		filename,
		UploadAvatar,
		nil,
		options,
	)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// DeleteGroupOptions represents the available DeleteGroup() options.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#delete-a-group
type DeleteGroupOptions struct {
	PermanentlyRemove *bool   `url:"permanently_remove,omitempty" json:"permanently_remove,omitempty"`
	FullPath          *string `url:"full_path,omitempty" json:"full_path,omitempty"`
}

// DeleteGroup removes group with all projects inside.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#delete-a-group
func (s *GroupsService) DeleteGroup(gid any, opt *DeleteGroupOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RestoreGroup restores a previously deleted group
//
// GitLap API docs:
// https://docs.gitlab.com/api/groups/#restore-a-group-marked-for-deletion
func (s *GroupsService) RestoreGroup(gid any, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/restore", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// SearchGroup get all groups that match your string in their name or path.
//
// GitLab API docs: https://docs.gitlab.com/api/groups/#search-for-a-group
func (s *GroupsService) SearchGroup(query string, options ...RequestOptionFunc) ([]*Group, *Response, error) {
	var q struct {
		Search string `url:"search,omitempty" json:"search,omitempty"`
	}
	q.Search = query

	req, err := s.client.NewRequest(http.MethodGet, "groups", &q, options)
	if err != nil {
		return nil, nil, err
	}

	var gs []*Group
	resp, err := s.client.Do(req, &gs)
	if err != nil {
		return nil, resp, err
	}

	return gs, resp, nil
}

// ListProvisionedUsersOptions represents the available ListProvisionedUsers()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-provisioned-users
type ListProvisionedUsersOptions struct {
	ListOptions
	Username      *string    `url:"username,omitempty" json:"username,omitempty"`
	Search        *string    `url:"search,omitempty" json:"search,omitempty"`
	Active        *bool      `url:"active,omitempty" json:"active,omitempty"`
	Blocked       *bool      `url:"blocked,omitempty" json:"blocked,omitempty"`
	CreatedAfter  *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
}

// ListProvisionedUsers gets a list of users provisioned by the given group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-provisioned-users
func (s *GroupsService) ListProvisionedUsers(gid any, opt *ListProvisionedUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/provisioned_users", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var us []*User
	resp, err := s.client.Do(req, &us)
	if err != nil {
		return nil, resp, err
	}

	return us, resp, nil
}

// ListGroupLDAPLinks lists the group's LDAP links. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#list-ldap-group-links
func (s *GroupsService) ListGroupLDAPLinks(gid any, options ...RequestOptionFunc) ([]*LDAPGroupLink, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/ldap_group_links", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var gls []*LDAPGroupLink
	resp, err := s.client.Do(req, &gls)
	if err != nil {
		return nil, resp, err
	}

	return gls, resp, nil
}

// AddGroupLDAPLinkOptions represents the available AddGroupLDAPLink() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#add-an-ldap-group-link-with-cn-or-filter
type AddGroupLDAPLinkOptions struct {
	CN           *string           `url:"cn,omitempty" json:"cn,omitempty"`
	Filter       *string           `url:"filter,omitempty" json:"filter,omitempty"`
	GroupAccess  *AccessLevelValue `url:"group_access,omitempty" json:"group_access,omitempty"`
	Provider     *string           `url:"provider,omitempty" json:"provider,omitempty"`
	MemberRoleID *int64            `url:"member_role_id,omitempty" json:"member_role_id,omitempty"`
}

// AddGroupLDAPLink creates a new group LDAP link. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#add-an-ldap-group-link-with-cn-or-filter
func (s *GroupsService) AddGroupLDAPLink(gid any, opt *AddGroupLDAPLinkOptions, options ...RequestOptionFunc) (*LDAPGroupLink, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/ldap_group_links", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gl := new(LDAPGroupLink)
	resp, err := s.client.Do(req, gl)
	if err != nil {
		return nil, resp, err
	}

	return gl, resp, nil
}

// DeleteGroupLDAPLink deletes a group LDAP link. Available only for users who
// can edit groups.
// Deprecated as upstream API is deprecated. Use DeleteGroupLDAPLinkWithCNOrFilter() instead.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#delete-an-ldap-group-link-deprecated
func (s *GroupsService) DeleteGroupLDAPLink(gid any, cn string, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/ldap_group_links/%s", PathEscape(group), PathEscape(cn))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteGroupLDAPLinkWithCNOrFilterOptions represents the available DeleteGroupLDAPLinkWithCNOrFilter() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#delete-an-ldap-group-link-with-cn-or-filter
type DeleteGroupLDAPLinkWithCNOrFilterOptions struct {
	CN       *string `url:"cn,omitempty" json:"cn,omitempty"`
	Filter   *string `url:"filter,omitempty" json:"filter,omitempty"`
	Provider *string `url:"provider,omitempty" json:"provider,omitempty"`
}

// DeleteGroupLDAPLinkWithCNOrFilter deletes a group LDAP link. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#delete-an-ldap-group-link-with-cn-or-filter
func (s *GroupsService) DeleteGroupLDAPLinkWithCNOrFilter(gid any, opts *DeleteGroupLDAPLinkWithCNOrFilterOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/ldap_group_links", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, opts, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteGroupLDAPLinkForProvider deletes a group LDAP link from a specific
// provider. Available only for users who can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ldap_links/#delete-an-ldap-group-link-deprecated
func (s *GroupsService) DeleteGroupLDAPLinkForProvider(gid any, provider, cn string, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(
		"groups/%s/ldap_group_links/%s/%s",
		PathEscape(group),
		PathEscape(provider),
		PathEscape(cn),
	)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListGroupSAMLLinks lists the group's SAML links. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/saml/#list-saml-group-links
func (s *GroupsService) ListGroupSAMLLinks(gid any, options ...RequestOptionFunc) ([]*SAMLGroupLink, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/saml_group_links", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var gl []*SAMLGroupLink
	resp, err := s.client.Do(req, &gl)
	if err != nil {
		return nil, resp, err
	}

	return gl, resp, nil
}

// ListGroupSharedProjectsOptions represents the available ListGroupSharedProjects() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-shared-projects
type ListGroupSharedProjectsOptions struct {
	ListOptions
	Archived                 *bool             `url:"archived,omitempty" json:"archived,omitempty"`
	MinAccessLevel           *AccessLevelValue `url:"min_access_level,omitempty" json:"min_access_level,omitempty"`
	OrderBy                  *string           `url:"order_by,omitempty" json:"order_by,omitempty"`
	Search                   *string           `url:"search,omitempty" json:"search,omitempty"`
	Simple                   *bool             `url:"simple,omitempty" json:"simple,omitempty"`
	Sort                     *string           `url:"sort,omitempty" json:"sort,omitempty"`
	Starred                  *bool             `url:"starred,omitempty" json:"starred,omitempty"`
	Visibility               *VisibilityValue  `url:"visibility,omitempty" json:"visibility,omitempty"`
	WithCustomAttributes     *bool             `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
	WithIssuesEnabled        *bool             `url:"with_issues_enabled,omitempty" json:"with_issues_enabled,omitempty"`
	WithMergeRequestsEnabled *bool             `url:"with_merge_requests_enabled,omitempty" json:"with_merge_requests_enabled,omitempty"`
}

// ListGroupSharedProjects gets a list of projects shared to this group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-shared-projects
func (s *GroupsService) ListGroupSharedProjects(gid any, opt *ListGroupSharedProjectsOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/projects/shared", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var p []*Project
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// GetGroupSAMLLink get a specific group SAML link. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/saml/#get-a-saml-group-link
func (s *GroupsService) GetGroupSAMLLink(gid any, samlGroupName string, options ...RequestOptionFunc) (*SAMLGroupLink, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/saml_group_links/%s", PathEscape(group), PathEscape(samlGroupName))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	gl := new(SAMLGroupLink)
	resp, err := s.client.Do(req, &gl)
	if err != nil {
		return nil, resp, err
	}

	return gl, resp, nil
}

// AddGroupSAMLLinkOptions represents the available AddGroupSAMLLink() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/saml/#add-a-saml-group-link
type AddGroupSAMLLinkOptions struct {
	SAMLGroupName *string           `url:"saml_group_name,omitempty" json:"saml_group_name,omitempty"`
	AccessLevel   *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	MemberRoleID  *int              `url:"member_role_id,omitempty" json:"member_role_id,omitempty"`
}

// AddGroupSAMLLink creates a new group SAML link. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/saml/#add-a-saml-group-link
func (s *GroupsService) AddGroupSAMLLink(gid any, opt *AddGroupSAMLLinkOptions, options ...RequestOptionFunc) (*SAMLGroupLink, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/saml_group_links", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gl := new(SAMLGroupLink)
	resp, err := s.client.Do(req, &gl)
	if err != nil {
		return nil, resp, err
	}

	return gl, resp, nil
}

// DeleteGroupSAMLLink deletes a group SAML link. Available only for users who
// can edit groups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/saml/#delete-a-saml-group-link
func (s *GroupsService) DeleteGroupSAMLLink(gid any, samlGroupName string, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/saml_group_links/%s", PathEscape(group), PathEscape(samlGroupName))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ShareGroupWithGroupOptions represents the available ShareGroupWithGroup() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#share-groups-with-groups
type ShareGroupWithGroupOptions struct {
	GroupID      *int              `url:"group_id,omitempty" json:"group_id,omitempty"`
	GroupAccess  *AccessLevelValue `url:"group_access,omitempty" json:"group_access,omitempty"`
	ExpiresAt    *ISOTime          `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	MemberRoleID *int              `url:"member_role_id,omitempty" json:"member_role_id,omitempty"`
}

// ShareGroupWithGroup shares a group with another group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#create-a-link-to-share-a-group-with-another-group
func (s *GroupsService) ShareGroupWithGroup(gid any, opt *ShareGroupWithGroupOptions, options ...RequestOptionFunc) (*Group, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/share", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	g := new(Group)
	resp, err := s.client.Do(req, g)
	if err != nil {
		return nil, resp, err
	}

	return g, resp, nil
}

// UnshareGroupFromGroup unshares a group from another group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#delete-the-link-that-shares-a-group-with-another-group
func (s *GroupsService) UnshareGroupFromGroup(gid any, groupID int, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/share/%d", PathEscape(group), groupID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GroupPushRules represents a group push rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#get-the-push-rules-of-a-group
type GroupPushRules struct {
	ID                         int        `json:"id"`
	CreatedAt                  *time.Time `json:"created_at"`
	CommitMessageRegex         string     `json:"commit_message_regex"`
	CommitMessageNegativeRegex string     `json:"commit_message_negative_regex"`
	BranchNameRegex            string     `json:"branch_name_regex"`
	DenyDeleteTag              bool       `json:"deny_delete_tag"`
	MemberCheck                bool       `json:"member_check"`
	PreventSecrets             bool       `json:"prevent_secrets"`
	AuthorEmailRegex           string     `json:"author_email_regex"`
	FileNameRegex              string     `json:"file_name_regex"`
	MaxFileSize                int        `json:"max_file_size"`
	CommitCommitterCheck       bool       `json:"commit_committer_check"`
	CommitCommitterNameCheck   bool       `json:"commit_committer_name_check"`
	RejectUnsignedCommits      bool       `json:"reject_unsigned_commits"`
	RejectNonDCOCommits        bool       `json:"reject_non_dco_commits"`
}

// GetGroupPushRules gets the push rules of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#get-the-push-rules-of-a-group
func (s *GroupsService) GetGroupPushRules(gid any, options ...RequestOptionFunc) (*GroupPushRules, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/push_rule", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	gpr := new(GroupPushRules)
	resp, err := s.client.Do(req, gpr)
	if err != nil {
		return nil, resp, err
	}

	return gpr, resp, nil
}

// AddGroupPushRuleOptions represents the available AddGroupPushRule()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#add-push-rules-to-a-group
type AddGroupPushRuleOptions struct {
	AuthorEmailRegex           *string `url:"author_email_regex,omitempty" json:"author_email_regex,omitempty"`
	BranchNameRegex            *string `url:"branch_name_regex,omitempty" json:"branch_name_regex,omitempty"`
	CommitCommitterCheck       *bool   `url:"commit_committer_check,omitempty" json:"commit_committer_check,omitempty"`
	CommitCommitterNameCheck   *bool   `url:"commit_committer_name_check,omitempty" json:"commit_committer_name_check,omitempty"`
	CommitMessageNegativeRegex *string `url:"commit_message_negative_regex,omitempty" json:"commit_message_negative_regex,omitempty"`
	CommitMessageRegex         *string `url:"commit_message_regex,omitempty" json:"commit_message_regex,omitempty"`
	DenyDeleteTag              *bool   `url:"deny_delete_tag,omitempty" json:"deny_delete_tag,omitempty"`
	FileNameRegex              *string `url:"file_name_regex,omitempty" json:"file_name_regex,omitempty"`
	MaxFileSize                *int    `url:"max_file_size,omitempty" json:"max_file_size,omitempty"`
	MemberCheck                *bool   `url:"member_check,omitempty" json:"member_check,omitempty"`
	PreventSecrets             *bool   `url:"prevent_secrets,omitempty" json:"prevent_secrets,omitempty"`
	RejectUnsignedCommits      *bool   `url:"reject_unsigned_commits,omitempty" json:"reject_unsigned_commits,omitempty"`
	RejectNonDCOCommits        *bool   `url:"reject_non_dco_commits,omitempty" json:"reject_non_dco_commits,omitempty"`
}

// AddGroupPushRule adds push rules to the specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#add-push-rules-to-a-group
func (s *GroupsService) AddGroupPushRule(gid any, opt *AddGroupPushRuleOptions, options ...RequestOptionFunc) (*GroupPushRules, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/push_rule", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gpr := new(GroupPushRules)
	resp, err := s.client.Do(req, gpr)
	if err != nil {
		return nil, resp, err
	}

	return gpr, resp, nil
}

// EditGroupPushRuleOptions represents the available EditGroupPushRule()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#edit-the-push-rules-of-a-group
type EditGroupPushRuleOptions struct {
	AuthorEmailRegex           *string `url:"author_email_regex,omitempty" json:"author_email_regex,omitempty"`
	BranchNameRegex            *string `url:"branch_name_regex,omitempty" json:"branch_name_regex,omitempty"`
	CommitCommitterCheck       *bool   `url:"commit_committer_check,omitempty" json:"commit_committer_check,omitempty"`
	CommitCommitterNameCheck   *bool   `url:"commit_committer_name_check,omitempty" json:"commit_committer_name_check,omitempty"`
	CommitMessageNegativeRegex *string `url:"commit_message_negative_regex,omitempty" json:"commit_message_negative_regex,omitempty"`
	CommitMessageRegex         *string `url:"commit_message_regex,omitempty" json:"commit_message_regex,omitempty"`
	DenyDeleteTag              *bool   `url:"deny_delete_tag,omitempty" json:"deny_delete_tag,omitempty"`
	FileNameRegex              *string `url:"file_name_regex,omitempty" json:"file_name_regex,omitempty"`
	MaxFileSize                *int    `url:"max_file_size,omitempty" json:"max_file_size,omitempty"`
	MemberCheck                *bool   `url:"member_check,omitempty" json:"member_check,omitempty"`
	PreventSecrets             *bool   `url:"prevent_secrets,omitempty" json:"prevent_secrets,omitempty"`
	RejectUnsignedCommits      *bool   `url:"reject_unsigned_commits,omitempty" json:"reject_unsigned_commits,omitempty"`
	RejectNonDCOCommits        *bool   `url:"reject_non_dco_commits,omitempty" json:"reject_non_dco_commits,omitempty"`
}

// EditGroupPushRule edits a push rule for a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#edit-the-push-rules-of-a-group
func (s *GroupsService) EditGroupPushRule(gid any, opt *EditGroupPushRuleOptions, options ...RequestOptionFunc) (*GroupPushRules, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/push_rule", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gpr := new(GroupPushRules)
	resp, err := s.client.Do(req, gpr)
	if err != nil {
		return nil, resp, err
	}

	return gpr, resp, nil
}

// DeleteGroupPushRule deletes the push rules of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_push_rules/#delete-the-push-rules-of-a-group
func (s *GroupsService) DeleteGroupPushRule(gid any, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/push_rule", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
