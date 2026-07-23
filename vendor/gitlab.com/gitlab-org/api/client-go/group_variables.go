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
	// GroupVariablesServiceInterface defines methods for the GroupVariablesService.
	GroupVariablesServiceInterface interface {
		// ListVariables gets a list of all variables for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_level_variables/#list-group-variables
		ListVariables(gid any, opt *ListGroupVariablesOptions, options ...RequestOptionFunc) ([]*GroupVariable, *Response, error)
		// GetVariable gets a variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_level_variables/#show-variable-details
		GetVariable(gid any, key string, opt *GetGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error)
		// CreateVariable creates a new group variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_level_variables/#create-variable
		CreateVariable(gid any, opt *CreateGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error)
		// UpdateVariable updates an existing group variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_level_variables/#update-variable
		UpdateVariable(gid any, key string, opt *UpdateGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error)
		// RemoveVariable removes a group's variable.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_level_variables/#remove-variable
		RemoveVariable(gid any, key string, opt *RemoveGroupVariableOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupVariablesService handles communication with the
	// group variables related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_level_variables/
	GroupVariablesService struct {
		client *Client
	}
)

var _ GroupVariablesServiceInterface = (*GroupVariablesService)(nil)

// GroupVariable represents a GitLab group Variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/
type GroupVariable struct {
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

func (v GroupVariable) String() string {
	return Stringify(v)
}

// ListGroupVariablesOptions represents the available options for listing variables
// for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/#list-group-variables
type ListGroupVariablesOptions struct {
	ListOptions
}

func (s *GroupVariablesService) ListVariables(gid any, opt *ListGroupVariablesOptions, options ...RequestOptionFunc) ([]*GroupVariable, *Response, error) {
	return do[[]*GroupVariable](s.client,
		withPath("groups/%s/variables", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupVariableOptions represents the available GetVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/#show-variable-details
type GetGroupVariableOptions struct {
	Filter *VariableFilter `url:"filter,omitempty" json:"filter,omitempty"`
}

func (s *GroupVariablesService) GetVariable(gid any, key string, opt *GetGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error) {
	return do[*GroupVariable](s.client,
		withPath("groups/%s/variables/%s", GroupID{gid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateGroupVariableOptions represents the available CreateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/#create-variable
type CreateGroupVariableOptions struct {
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

func (s *GroupVariablesService) CreateVariable(gid any, opt *CreateGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error) {
	return do[*GroupVariable](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/variables", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateGroupVariableOptions represents the available UpdateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/#update-variable
type UpdateGroupVariableOptions struct {
	Value            *string            `url:"value,omitempty" json:"value,omitempty"`
	Description      *string            `url:"description,omitempty" json:"description,omitempty"`
	EnvironmentScope *string            `url:"environment_scope,omitempty" json:"environment_scope,omitempty"`
	Filter           *VariableFilter    `url:"filter,omitempty" json:"filter,omitempty"`
	Masked           *bool              `url:"masked,omitempty" json:"masked,omitempty"`
	Protected        *bool              `url:"protected,omitempty" json:"protected,omitempty"`
	Raw              *bool              `url:"raw,omitempty" json:"raw,omitempty"`
	VariableType     *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

func (s *GroupVariablesService) UpdateVariable(gid any, key string, opt *UpdateGroupVariableOptions, options ...RequestOptionFunc) (*GroupVariable, *Response, error) {
	return do[*GroupVariable](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/variables/%s", GroupID{gid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RemoveGroupVariableOptions represents the available RemoveVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_level_variables/#remove-variable
type RemoveGroupVariableOptions struct {
	Filter *VariableFilter `url:"filter,omitempty" json:"filter,omitempty"`
}

func (s *GroupVariablesService) RemoveVariable(gid any, key string, opt *RemoveGroupVariableOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/variables/%s", GroupID{gid}, url.PathEscape(key)),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
