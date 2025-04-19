//
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
	GroupSecuritySettingsServiceInterface interface {
		UpdateSecretPushProtectionEnabledSetting(gid interface{}, opt UpdateGroupSecuritySettingsOptions, options ...RequestOptionFunc) (*GroupSecuritySettings, *Response, error)
	}

	// GroupSecuritySettingsService handles communication with the Group Security Settings
	// related methods of the GitLab API.
	//
	// Gitlab API docs:
	// https://docs.gitlab.com/api/group_security_settings/
	GroupSecuritySettingsService struct {
		client *Client
	}
)

var _ GroupSecuritySettingsServiceInterface = (*GroupSecuritySettingsService)(nil)

// GroupSecuritySettings represents the group security settings data.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/group_security_settings/
type GroupSecuritySettings struct {
	SecretPushProtectionEnabled bool     `json:"secret_push_protection_enabled"`
	Errors                      []string `json:"errors"`
}

// Gets a string representation of the GroupSecuritySettings data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_security_settings/
func (s GroupSecuritySettings) String() string {
	return Stringify(s)
}

// GetGroupSecuritySettingsOptions represent the request options for updating
// the group security settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_security_settings/#update-secret_push_protection_enabled-setting
type UpdateGroupSecuritySettingsOptions struct {
	SecretPushProtectionEnabled *bool  `url:"secret_push_protection_enabled,omitempty" json:"secret_push_protection_enabled,omitempty"`
	ProjectsToExclude           *[]int `url:"projects_to_exclude,omitempty" json:"projects_to_exclude,omitempty"`
}

// UpdateSecretPushProtectionEnabledSetting updates the secret_push_protection_enabled
// setting for the all projects in a group to the provided value.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_security_settings/#update-secret_push_protection_enabled-setting
func (s *GroupSecuritySettingsService) UpdateSecretPushProtectionEnabledSetting(gid interface{}, opt UpdateGroupSecuritySettingsOptions, options ...RequestOptionFunc) (*GroupSecuritySettings, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/security_settings", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}
	settings := new(GroupSecuritySettings)
	resp, err := s.client.Do(req, &settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, err
}
