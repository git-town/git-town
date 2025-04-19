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
	"fmt"
	"net/http"
)

type (
	MergeRequestApprovalSettingsServiceInterface interface {
		GetGroupMergeRequestApprovalSettings(gid interface{}, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error)
		UpdateGroupMergeRequestApprovalSettings(gid interface{}, opt *UpdateMergeRequestApprovalSettingsOptions, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error)
		GetProjectMergeRequestApprovalSettings(pid interface{}, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error)
		UpdateProjectMergeRequestApprovalSettings(pid interface{}, opt *UpdateMergeRequestApprovalSettingsOptions, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error)
	}

	// MergeRequestApprovalSettingsService handles communication with the merge
	// requests approval settings related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/merge_request_approval_settings/
	MergeRequestApprovalSettingsService struct {
		client *Client
	}
)

var _ MergeRequestApprovalSettingsServiceInterface = (*MergeRequestApprovalSettingsService)(nil)

// MergeRequestApprovalSettings represents the merge request approval settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/
type MergeRequestApprovalSettings struct {
	AllowAuthorApproval                         MergeRequestApprovalSetting `json:"allow_author_approval"`
	AllowCommitterApproval                      MergeRequestApprovalSetting `json:"allow_committer_approval"`
	AllowOverridesToApproverListPerMergeRequest MergeRequestApprovalSetting `json:"allow_overrides_to_approver_list_per_merge_request"`
	RetainApprovalsOnPush                       MergeRequestApprovalSetting `json:"retain_approvals_on_push"`
	SelectiveCodeOwnerRemovals                  MergeRequestApprovalSetting `json:"selective_code_owner_removals"`
	RequirePasswordToApprove                    MergeRequestApprovalSetting `json:"require_password_to_approve"`
	RequireReauthenticationToApprove            MergeRequestApprovalSetting `json:"require_reauthentication_to_approve"`
}

// MergeRequestApprovalSetting represents an individual merge request approval
// setting.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/
type MergeRequestApprovalSetting struct {
	Value         bool   `json:"value"`
	Locked        bool   `json:"locked"`
	InheritedFrom string `json:"inherited_from"`
}

// GetGroupMergeRequestApprovalSettings gets the merge request approval settings
// of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/#get-group-mr-approval-settings
func (s *MergeRequestApprovalSettingsService) GetGroupMergeRequestApprovalSettings(gid interface{}, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/merge_request_approval_setting", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	settings := new(MergeRequestApprovalSettings)
	resp, err := s.client.Do(req, settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

// UpdateMergeRequestApprovalSettingsOptions represents the available
// UpdateGroupMergeRequestApprovalSettings() and UpdateProjectMergeRequestApprovalSettings()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/#update-group-mr-approval-settings
// https://docs.gitlab.com/api/merge_request_approval_settings/#update-project-mr-approval-settings
type UpdateMergeRequestApprovalSettingsOptions struct {
	AllowAuthorApproval                         *bool `url:"allow_author_approval,omitempty" json:"allow_author_approval,omitempty"`
	AllowCommitterApproval                      *bool `url:"allow_committer_approval,omitempty" json:"allow_committer_approval,omitempty"`
	AllowOverridesToApproverListPerMergeRequest *bool `url:"allow_overrides_to_approver_list_per_merge_request,omitempty" json:"allow_overrides_to_approver_list_per_merge_request,omitempty"`
	RetainApprovalsOnPush                       *bool `url:"retain_approvals_on_push,omitempty" json:"retain_approvals_on_push,omitempty"`
	SelectiveCodeOwnerRemovals                  *bool `url:"selective_code_owner_removals,omitempty" json:"selective_code_owner_removals,omitempty"`
	RequireReauthenticationToApprove            *bool `url:"require_reauthentication_to_approve,omitempty" json:"require_reauthentication_to_approve,omitempty"`
}

// UpdateGroupMergeRequestApprovalSettings updates the merge request approval
// settings of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/#update-group-mr-approval-settings
func (s *MergeRequestApprovalSettingsService) UpdateGroupMergeRequestApprovalSettings(gid interface{}, opt *UpdateMergeRequestApprovalSettingsOptions, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/merge_request_approval_setting", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	settings := new(MergeRequestApprovalSettings)
	resp, err := s.client.Do(req, settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

// GetProjectMergeRequestApprovalSettings gets the merge request approval settings
// of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/#get-project-mr-approval-settings
func (s *MergeRequestApprovalSettingsService) GetProjectMergeRequestApprovalSettings(pid interface{}, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_request_approval_setting", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	settings := new(MergeRequestApprovalSettings)
	resp, err := s.client.Do(req, settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

// UpdateProjectMergeRequestApprovalSettings updates the merge request approval
// settings of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approval_settings/#update-project-mr-approval-settings
func (s *MergeRequestApprovalSettingsService) UpdateProjectMergeRequestApprovalSettings(pid interface{}, opt *UpdateMergeRequestApprovalSettingsOptions, options ...RequestOptionFunc) (*MergeRequestApprovalSettings, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_request_approval_setting", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	settings := new(MergeRequestApprovalSettings)
	resp, err := s.client.Do(req, settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}
