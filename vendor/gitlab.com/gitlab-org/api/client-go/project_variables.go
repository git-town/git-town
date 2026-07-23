//
// Copyright 2021, Patrick Webster
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
	"net/url"
)

type (
	ProjectVariablesServiceInterface interface {
		// ListVariables gets a list of all variables in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_level_variables/#list-project-variables
		ListVariables(pid any, opt *ListProjectVariablesOptions, options ...RequestOptionFunc) ([]*ProjectVariable, *Response, error)
		// GetVariable gets a variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_level_variables/#get-a-single-variable
		GetVariable(pid any, key string, opt *GetProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error)
		// CreateVariable creates a new project variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_level_variables/#create-a-variable
		CreateVariable(pid any, opt *CreateProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error)
		// UpdateVariable updates a project's variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_level_variables/#update-a-variable
		UpdateVariable(pid any, key string, opt *UpdateProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error)
		// RemoveVariable removes a project's variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_level_variables/#delete-a-variable
		RemoveVariable(pid any, key string, opt *RemoveProjectVariableOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectVariablesService handles communication with the
	// project variables related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_level_variables/
	ProjectVariablesService struct {
		client *Client
	}
)

var _ ProjectVariablesServiceInterface = (*ProjectVariablesService)(nil)

// ProjectVariable represents a GitLab Project Variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/
type ProjectVariable struct {
	Key              string            `json:"key"`
	Value            string            `json:"value"`
	VariableType     VariableTypeValue `json:"variable_type"`
	Protected        bool              `json:"protected"`
	Masked           bool              `json:"masked"`
	Hidden           bool              `json:"hidden"`
	Raw              bool              `json:"raw"`
	EnvironmentScope string            `json:"environment_scope"`
	Description      string            `json:"description"`
}

func (v ProjectVariable) String() string {
	return Stringify(v)
}

// VariableFilter filters available for project variable related functions
type VariableFilter struct {
	EnvironmentScope string `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
}

// ListProjectVariablesOptions represents the available options for listing variables
// in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/#list-project-variables
type ListProjectVariablesOptions struct {
	ListOptions
}

func (s *ProjectVariablesService) ListVariables(pid any, opt *ListProjectVariablesOptions, options ...RequestOptionFunc) ([]*ProjectVariable, *Response, error) {
	return do[[]*ProjectVariable](s.client,
		withPath("projects/%s/variables", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProjectVariableOptions represents the available GetVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/#get-a-single-variable
type GetProjectVariableOptions struct {
	Filter *VariableFilter `url:"filter,omitempty" json:"filter,omitempty"`
}

func (s *ProjectVariablesService) GetVariable(pid any, key string, opt *GetProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error) {
	return do[*ProjectVariable](s.client,
		withPath("projects/%s/variables/%s", ProjectID{pid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateProjectVariableOptions represents the available CreateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/#create-a-variable
type CreateProjectVariableOptions struct {
	Key              *string            `url:"key,omitempty" json:"key,omitempty"`
	Value            *string            `url:"value,omitempty" json:"value,omitempty"`
	Description      *string            `url:"description,omitempty" json:"description,omitempty"`
	EnvironmentScope *string            `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
	Masked           *bool              `url:"masked,omitempty" json:"masked,omitempty"`
	MaskedAndHidden  *bool              `url:"masked_and_hidden,omitempty" json:"masked_and_hidden,omitempty"`
	Protected        *bool              `url:"protected,omitempty" json:"protected,omitempty"`
	Raw              *bool              `url:"raw,omitempty" json:"raw,omitempty"`
	VariableType     *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

func (s *ProjectVariablesService) CreateVariable(pid any, opt *CreateProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error) {
	return do[*ProjectVariable](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/variables", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateProjectVariableOptions represents the available UpdateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/#update-a-variable
type UpdateProjectVariableOptions struct {
	Value            *string            `url:"value,omitempty" json:"value,omitempty"`
	Description      *string            `url:"description,omitempty" json:"description,omitempty"`
	EnvironmentScope *string            `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
	Filter           *VariableFilter    `url:"filter,omitempty" json:"filter,omitempty"`
	Masked           *bool              `url:"masked,omitempty" json:"masked,omitempty"`
	Protected        *bool              `url:"protected,omitempty" json:"protected,omitempty"`
	Raw              *bool              `url:"raw,omitempty" json:"raw,omitempty"`
	VariableType     *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

func (s *ProjectVariablesService) UpdateVariable(pid any, key string, opt *UpdateProjectVariableOptions, options ...RequestOptionFunc) (*ProjectVariable, *Response, error) {
	return do[*ProjectVariable](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/variables/%s", ProjectID{pid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RemoveProjectVariableOptions represents the available RemoveVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_level_variables/#delete-a-variable
type RemoveProjectVariableOptions struct {
	Filter *VariableFilter `url:"filter,omitempty" json:"filter,omitempty"`
}

func (s *ProjectVariablesService) RemoveVariable(pid any, key string, opt *RemoveProjectVariableOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/variables/%s", ProjectID{pid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
