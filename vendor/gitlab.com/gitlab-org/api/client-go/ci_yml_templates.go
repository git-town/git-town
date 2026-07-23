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
	"net/http"
)

type (
	CIYMLTemplatesServiceInterface interface {
		// ListAllTemplates get all GitLab CI YML templates.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/#list-gitlab-ci-yaml-templates
		ListAllTemplates(opt *ListCIYMLTemplatesOptions, options ...RequestOptionFunc) ([]*CIYMLTemplateListItem, *Response, error)

		// GetTemplate get a single GitLab CI YML template.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/#single-gitlab-ci-yaml-template
		GetTemplate(key string, options ...RequestOptionFunc) (*CIYMLTemplate, *Response, error)
	}

	// CIYMLTemplatesService handles communication with the gitlab
	// CI YML templates related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/
	CIYMLTemplatesService struct {
		client *Client
	}
)

var _ CIYMLTemplatesServiceInterface = (*CIYMLTemplatesService)(nil)

// CIYMLTemplate represents a GitLab CI YML template.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/
type CIYMLTemplate struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// CIYMLTemplateListItem represents a GitLab CI YML template from the list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/
type CIYMLTemplateListItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ListCIYMLTemplatesOptions represents the available ListAllTemplates() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/gitlab_ci_ymls/#list-gitlab-ci-yaml-templates
type ListCIYMLTemplatesOptions struct {
	ListOptions
}

func (s *CIYMLTemplatesService) ListAllTemplates(opt *ListCIYMLTemplatesOptions, options ...RequestOptionFunc) ([]*CIYMLTemplateListItem, *Response, error) {
	return do[[]*CIYMLTemplateListItem](s.client,
		withMethod(http.MethodGet),
		withPath("templates/gitlab_ci_ymls"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *CIYMLTemplatesService) GetTemplate(key string, options ...RequestOptionFunc) (*CIYMLTemplate, *Response, error) {
	return do[*CIYMLTemplate](s.client,
		withMethod(http.MethodGet),
		withPath("templates/gitlab_ci_ymls/%s", key),
		withRequestOpts(options...),
	)
}
