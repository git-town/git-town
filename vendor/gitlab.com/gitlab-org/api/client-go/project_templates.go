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

type (
	ProjectTemplatesServiceInterface interface {
		// ListTemplates gets a list of project templates.
		//
		// GitLab API docs: https://docs.gitlab.com/api/project_templates/#get-all-templates-of-a-particular-type
		ListTemplates(pid any, templateType string, opt *ListProjectTemplatesOptions, options ...RequestOptionFunc) ([]*ProjectTemplate, *Response, error)
		// GetProjectTemplate gets a single project template.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_templates/#get-one-template-of-a-particular-type
		GetProjectTemplate(pid any, templateType string, templateName string, options ...RequestOptionFunc) (*ProjectTemplate, *Response, error)
	}

	// ProjectTemplatesService handles communication with the project templates
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/project_templates/
	ProjectTemplatesService struct {
		client *Client
	}
)

var _ ProjectTemplatesServiceInterface = (*ProjectTemplatesService)(nil)

// ProjectTemplate represents a GitLab ProjectTemplate.
//
// GitLab API docs: https://docs.gitlab.com/api/project_templates/
type ProjectTemplate struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Popular     bool     `json:"popular"`
	HTMLURL     string   `json:"html_url"`
	SourceURL   string   `json:"source_url"`
	Description string   `json:"description"`
	Conditions  []string `json:"conditions"`
	Permissions []string `json:"permissions"`
	Limitations []string `json:"limitations"`
	Content     string   `json:"content"`
}

func (s ProjectTemplate) String() string {
	return Stringify(s)
}

// ListProjectTemplatesOptions represents the available ListSnippets() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_templates/#get-all-templates-of-a-particular-type
type ListProjectTemplatesOptions struct {
	ListOptions
	ID   *int64  `url:"id,omitempty" json:"id,omitempty"`
	Type *string `url:"type,omitempty" json:"type,omitempty"`
}

func (s *ProjectTemplatesService) ListTemplates(pid any, templateType string, opt *ListProjectTemplatesOptions, options ...RequestOptionFunc) ([]*ProjectTemplate, *Response, error) {
	return do[[]*ProjectTemplate](s.client,
		withPath("projects/%s/templates/%s", ProjectID{pid}, NoEscape{templateType}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ProjectTemplatesService) GetProjectTemplate(pid any, templateType string, templateName string, options ...RequestOptionFunc) (*ProjectTemplate, *Response, error) {
	return do[*ProjectTemplate](s.client,
		withPath("projects/%s/templates/%s/%s", ProjectID{pid}, NoEscape{templateType}, NoEscape{templateName}),
		withRequestOpts(options...),
	)
}
