//
// Copyright 2022, FantasyTeddy
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
	// DockerfileTemplatesServiceInterface defines all the API methods for the DockerfileTemplatesService
	DockerfileTemplatesServiceInterface interface {
		// ListTemplates get a list of available Dockerfile templates.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/templates/dockerfiles/#list-dockerfile-templates
		ListTemplates(opt *ListDockerfileTemplatesOptions, options ...RequestOptionFunc) ([]*DockerfileTemplateListItem, *Response, error)

		// GetTemplate get a single Dockerfile template.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/templates/dockerfiles/#single-dockerfile-template
		GetTemplate(key string, options ...RequestOptionFunc) (*DockerfileTemplate, *Response, error)
	}

	// DockerfileTemplatesService handles communication with the Dockerfile
	// templates related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/templates/dockerfiles/
	DockerfileTemplatesService struct {
		client *Client
	}
)

var _ DockerfileTemplatesServiceInterface = (*DockerfileTemplatesService)(nil)

// DockerfileTemplate represents a GitLab Dockerfile template.
//
// GitLab API docs: https://docs.gitlab.com/api/templates/dockerfiles/
type DockerfileTemplate struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// DockerfileTemplateListItem represents a GitLab Dockerfile template from the list.
//
// GitLab API docs: https://docs.gitlab.com/api/templates/dockerfiles/
type DockerfileTemplateListItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ListDockerfileTemplatesOptions represents the available ListAllTemplates() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/templates/dockerfiles/#list-dockerfile-templates
type ListDockerfileTemplatesOptions struct {
	ListOptions
}

func (s *DockerfileTemplatesService) ListTemplates(opt *ListDockerfileTemplatesOptions, options ...RequestOptionFunc) ([]*DockerfileTemplateListItem, *Response, error) {
	return do[[]*DockerfileTemplateListItem](s.client,
		withPath("templates/dockerfiles"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DockerfileTemplatesService) GetTemplate(key string, options ...RequestOptionFunc) (*DockerfileTemplate, *Response, error) {
	return do[*DockerfileTemplate](s.client,
		withPath("templates/dockerfiles/%s", key),
		withRequestOpts(options...),
	)
}
