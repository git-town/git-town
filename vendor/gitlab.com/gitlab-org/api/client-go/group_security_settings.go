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
	"net/http"
)

type (
	GroupSecuritySettingsServiceInterface interface {
		UpdateSecretPushProtectionEnabledSetting(gid any, opt UpdateGroupSecuritySettingsOptions, options ...RequestOptionFunc) (*GroupSecuritySettings, *Response, error)
	}

	// GroupSecuritySettingsService handles communication with the Group Security Settings
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_security_settings/
	GroupSecuritySettingsService struct {
		client *Client
	}
)

var _ GroupSecuritySettingsServiceInterface = (*GroupSecuritySettingsService)(nil)

// GroupSecuritySettings represents the group security settings data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_security_settings/
type GroupSecuritySettings struct {
	SecretPushProtectionEnabled bool     `json:"secret_push_protection_enabled"`
	Errors                      []string `json:"errors"`
}

// String gets a string representation of the GroupSecuritySettings data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_security_settings/
func (s GroupSecuritySettings) String() string {
	return Stringify(s)
}

// UpdateGroupSecuritySettingsOptions represent the request options for updating
// the group security settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_security_settings/#update-secret_push_protection_enabled-setting
type UpdateGroupSecuritySettingsOptions struct {
	SecretPushProtectionEnabled *bool    `url:"secret_push_protection_enabled,omitempty" json:"secret_push_protection_enabled,omitempty"`
	ProjectsToExclude           *[]int64 `url:"projects_to_exclude,omitempty" json:"projects_to_exclude,omitempty"`
}

// UpdateSecretPushProtectionEnabledSetting updates the secret_push_protection_enabled
// setting for the all projects in a group to the provided value.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_security_settings/#update-secret_push_protection_enabled-setting
func (s *GroupSecuritySettingsService) UpdateSecretPushProtectionEnabledSetting(gid any, opt UpdateGroupSecuritySettingsOptions, options ...RequestOptionFunc) (*GroupSecuritySettings, *Response, error) {
	return do[*GroupSecuritySettings](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/security_settings", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
