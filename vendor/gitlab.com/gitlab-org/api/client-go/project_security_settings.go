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
	"time"
)

type (
	ProjectSecuritySettingsServiceInterface interface {
		ListProjectSecuritySettings(pid interface{}, options ...RequestOptionFunc) (*ProjectSecuritySettings, *Response, error)
		UpdateSecretPushProtectionEnabledSetting(pid interface{}, opt UpdateProjectSecuritySettingsOptions, options ...RequestOptionFunc) (*ProjectSecuritySettings, *Response, error)
	}

	// ProjectSecuritySettingsService handles communication with the Project Security Settings
	// related methods of the GitLab API.
	//
	// Gitlab API docs:
	// https://docs.gitlab.com/ee/api/project_security_settings.html
	ProjectSecuritySettingsService struct {
		client *Client
	}
)

var _ ProjectSecuritySettingsServiceInterface = (*ProjectSecuritySettingsService)(nil)

// ProjectSecuritySettings represents the project security settings data.
//
// Gitlab API docs:
// https://docs.gitlab.com/ee/api/project_security_settings.html
type ProjectSecuritySettings struct {
	ProjectID                           int64      `json:"project_id"`
	CreatedAt                           *time.Time `json:"created_at"`
	UpdatedAt                           *time.Time `json:"updated_at"`
	AutoFixContainerScanning            bool       `json:"auto_fix_container_scanning"`
	AutoFixDAST                         bool       `json:"auto_fix_dast"`
	AutoFixDependencyScanning           bool       `json:"auto_fix_dependency_scanning"`
	AutoFixSAST                         bool       `json:"auto_fix_sast"`
	ContinuousVulnerabilityScansEnabled bool       `json:"continuous_vulnerability_scans_enabled"`
	ContainerScanningForRegistryEnabled bool       `json:"container_scanning_for_registry_enabled"`
	SecretPushProtectionEnabled         bool       `json:"secret_push_protection_enabled"`
}

// Gets a string representation of the ProjectSecuritySettings data.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/project_security_settings.html
func (s ProjectSecuritySettings) String() string {
	return Stringify(s)
}

// ListProjectSecuritySettings lists all of a project's security settings.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_security_settings.html#list-project-security-settings
func (s *ProjectSecuritySettingsService) ListProjectSecuritySettings(pid interface{}, options ...RequestOptionFunc) (*ProjectSecuritySettings, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/security_settings", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}
	settings := new(ProjectSecuritySettings)
	resp, err := s.client.Do(req, &settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, err
}

// UpdateProjectSecuritySettingsOptions represent the request options for updating
// the project security settings.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/project_security_settings.html#update-secret_push_protection_enabled-setting
type UpdateProjectSecuritySettingsOptions struct {
	SecretPushProtectionEnabled *bool `url:"secret_push_protection_enabled,omitempty" json:"secret_push_protection_enabled,omitempty"`
}

// UpdateSecretPushProtectionEnabledSetting updates the secret_push_protection_enabled
// setting for the all projects in a project to the provided value.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_security_settings.html#update-secret_push_protection_enabled-setting
func (s *ProjectSecuritySettingsService) UpdateSecretPushProtectionEnabledSetting(pid interface{}, opt UpdateProjectSecuritySettingsOptions, options ...RequestOptionFunc) (*ProjectSecuritySettings, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/security_settings", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}
	settings := new(ProjectSecuritySettings)
	resp, err := s.client.Do(req, &settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, err
}
