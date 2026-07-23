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

// LicenseTemplate represents a license template.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/licenses/
type LicenseTemplate struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Featured    bool     `json:"featured"`
	HTMLURL     string   `json:"html_url"`
	SourceURL   string   `json:"source_url"`
	Description string   `json:"description"`
	Conditions  []string `json:"conditions"`
	Permissions []string `json:"permissions"`
	Limitations []string `json:"limitations"`
	Content     string   `json:"content"`
}

type (
	LicenseTemplatesServiceInterface interface {
		ListLicenseTemplates(opt *ListLicenseTemplatesOptions, options ...RequestOptionFunc) ([]*LicenseTemplate, *Response, error)
		GetLicenseTemplate(template string, opt *GetLicenseTemplateOptions, options ...RequestOptionFunc) (*LicenseTemplate, *Response, error)
	}

	// LicenseTemplatesService handles communication with the license templates
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/templates/licenses/
	LicenseTemplatesService struct {
		client *Client
	}
)

var _ LicenseTemplatesServiceInterface = (*LicenseTemplatesService)(nil)

// ListLicenseTemplatesOptions represents the available
// ListLicenseTemplates() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/licenses/#list-license-templates
type ListLicenseTemplatesOptions struct {
	ListOptions
	Popular *bool `url:"popular,omitempty" json:"popular,omitempty"`
}

// ListLicenseTemplates get all license templates.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/licenses/#list-license-templates
func (s *LicenseTemplatesService) ListLicenseTemplates(opt *ListLicenseTemplatesOptions, options ...RequestOptionFunc) ([]*LicenseTemplate, *Response, error) {
	return do[[]*LicenseTemplate](s.client,
		withPath("templates/licenses"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetLicenseTemplateOptions represents the available
// GetLicenseTemplate() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/licenses/#single-license-template
type GetLicenseTemplateOptions struct {
	Project  *string `url:"project,omitempty" json:"project,omitempty"`
	Fullname *string `url:"fullname,omitempty" json:"fullname,omitempty"`
}

// GetLicenseTemplate get a single license template. You can pass parameters
// to replace the license placeholder.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/licenses/#single-license-template
func (s *LicenseTemplatesService) GetLicenseTemplate(template string, opt *GetLicenseTemplateOptions, options ...RequestOptionFunc) (*LicenseTemplate, *Response, error) {
	return do[*LicenseTemplate](s.client,
		withPath("templates/licenses/%s", NoEscape{template}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
