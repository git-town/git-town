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
	AdminCompliancePolicySettingsServiceInterface interface {
		// GetCompliancePolicySettings gets the current security policy settings for the GitLab instance.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/compliance_policy_settings/#get-security-policy-settings
		GetCompliancePolicySettings(options ...RequestOptionFunc) (*AdminCompliancePolicySettings, *Response, error)

		// UpdateCompliancePolicySettings updates the security policy settings for the GitLab instance.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/compliance_policy_settings/#update-security-policy-settings
		UpdateCompliancePolicySettings(opt *UpdateAdminCompliancePolicySettingsOptions, options ...RequestOptionFunc) (*AdminCompliancePolicySettings, *Response, error)
	}

	// AdminCompliancePolicySettingsService handles communication with the
	// admin compliance policy settings related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/compliance_policy_settings/
	AdminCompliancePolicySettingsService struct {
		client *Client
	}
)

var _ AdminCompliancePolicySettingsServiceInterface = (*AdminCompliancePolicySettingsService)(nil)

// AdminCompliancePolicySettings represents the GitLab admin compliance policy settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/compliance_policy_settings/
type AdminCompliancePolicySettings struct {
	CSPNamespaceID *int64 `json:"csp_namespace_id"`
}

func (s AdminCompliancePolicySettings) String() string {
	return Stringify(s)
}

func (s *AdminCompliancePolicySettingsService) GetCompliancePolicySettings(options ...RequestOptionFunc) (*AdminCompliancePolicySettings, *Response, error) {
	return do[*AdminCompliancePolicySettings](s.client,
		withMethod(http.MethodGet),
		withPath("admin/security/compliance_policy_settings"),
		withRequestOpts(options...),
	)
}

// UpdateAdminCompliancePolicySettingsOptions represents the available
// UpdateCompliancePolicySettings() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/compliance_policy_settings/#update-security-policy-settings
type UpdateAdminCompliancePolicySettingsOptions struct {
	CSPNamespaceID *int64 `url:"csp_namespace_id,omitempty" json:"csp_namespace_id,omitempty"`
}

func (s *AdminCompliancePolicySettingsService) UpdateCompliancePolicySettings(opt *UpdateAdminCompliancePolicySettingsOptions, options ...RequestOptionFunc) (*AdminCompliancePolicySettings, *Response, error) {
	return do[*AdminCompliancePolicySettings](s.client,
		withMethod(http.MethodPut),
		withPath("admin/security/compliance_policy_settings"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
