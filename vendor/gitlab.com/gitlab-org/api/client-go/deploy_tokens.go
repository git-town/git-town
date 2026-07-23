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
	"time"
)

type (
	DeployTokensServiceInterface interface {
		// ListAllDeployTokens gets a list of all deploy tokens.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#list-all-deploy-tokens
		// ListAllDeployTokens gets a list of all deploy tokens.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#list-all-deploy-tokens
		ListAllDeployTokens(options ...RequestOptionFunc) ([]*DeployToken, *Response, error)

		// ListProjectDeployTokens gets a list of a project's deploy tokens.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#list-project-deploy-tokens

		// ListProjectDeployTokens gets a list of a project's deploy tokens.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#list-project-deploy-tokens
		ListProjectDeployTokens(pid any, opt *ListProjectDeployTokensOptions, options ...RequestOptionFunc) ([]*DeployToken, *Response, error)

		// GetProjectDeployToken gets a single deploy token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#get-a-project-deploy-token
		GetProjectDeployToken(pid any, deployToken int64, options ...RequestOptionFunc) (*DeployToken, *Response, error)

		// CreateProjectDeployToken creates a new deploy token for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#create-a-project-deploy-token
		CreateProjectDeployToken(pid any, opt *CreateProjectDeployTokenOptions, options ...RequestOptionFunc) (*DeployToken, *Response, error)

		// DeleteProjectDeployToken removes a deploy token from the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#delete-a-project-deploy-token
		DeleteProjectDeployToken(pid any, deployToken int64, options ...RequestOptionFunc) (*Response, error)

		// ListGroupDeployTokens gets a list of a group’s deploy tokens.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#list-group-deploy-tokens
		ListGroupDeployTokens(gid any, opt *ListGroupDeployTokensOptions, options ...RequestOptionFunc) ([]*DeployToken, *Response, error)

		// GetGroupDeployToken gets a single deploy token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#get-a-group-deploy-token
		GetGroupDeployToken(gid any, deployToken int64, options ...RequestOptionFunc) (*DeployToken, *Response, error)

		// CreateGroupDeployToken creates a new deploy token for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#create-a-group-deploy-token
		CreateGroupDeployToken(gid any, opt *CreateGroupDeployTokenOptions, options ...RequestOptionFunc) (*DeployToken, *Response, error)

		// DeleteGroupDeployToken removes a deploy token from the group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_tokens/#delete-a-group-deploy-token
		DeleteGroupDeployToken(gid any, deployToken int64, options ...RequestOptionFunc) (*Response, error)
	}

	// DeployTokensService handles communication with the deploy tokens related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/deploy_tokens/
	DeployTokensService struct {
		client *Client
	}
)

// DeployToken represents a GitLab deploy token.
type DeployToken struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	ExpiresAt *time.Time `json:"expires_at"`
	Revoked   bool       `json:"revoked"`
	Expired   bool       `json:"expired"`
	Token     string     `json:"token,omitempty"`
	Scopes    []string   `json:"scopes"`
}

func (k DeployToken) String() string {
	return Stringify(k)
}

func (s *DeployTokensService) ListAllDeployTokens(options ...RequestOptionFunc) ([]*DeployToken, *Response, error) {
	return do[[]*DeployToken](s.client,
		withPath("deploy_tokens"),
		withRequestOpts(options...),
	)
}

// ListProjectDeployTokensOptions represents the available ListProjectDeployTokens()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_tokens/#list-project-deploy-tokens
type ListProjectDeployTokensOptions struct {
	ListOptions
}

func (s *DeployTokensService) ListProjectDeployTokens(pid any, opt *ListProjectDeployTokensOptions, options ...RequestOptionFunc) ([]*DeployToken, *Response, error) {
	return do[[]*DeployToken](s.client,
		withPath("projects/%s/deploy_tokens", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployTokensService) GetProjectDeployToken(pid any, deployToken int64, options ...RequestOptionFunc) (*DeployToken, *Response, error) {
	return do[*DeployToken](s.client,
		withPath("projects/%s/deploy_tokens/%d", ProjectID{pid}, deployToken),
		withRequestOpts(options...),
	)
}

// CreateProjectDeployTokenOptions represents the available CreateProjectDeployToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_tokens/#create-a-project-deploy-token
type CreateProjectDeployTokenOptions struct {
	Name      *string    `url:"name,omitempty" json:"name,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	Username  *string    `url:"username,omitempty" json:"username,omitempty"`
	Scopes    *[]string  `url:"scopes,omitempty" json:"scopes,omitempty"`
}

func (s *DeployTokensService) CreateProjectDeployToken(pid any, opt *CreateProjectDeployTokenOptions, options ...RequestOptionFunc) (*DeployToken, *Response, error) {
	return do[*DeployToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/deploy_tokens", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployTokensService) DeleteProjectDeployToken(pid any, deployToken int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/deploy_tokens/%d", ProjectID{pid}, deployToken),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListGroupDeployTokensOptions represents the available ListGroupDeployTokens()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_tokens/#list-group-deploy-tokens
type ListGroupDeployTokensOptions struct {
	ListOptions
}

func (s *DeployTokensService) ListGroupDeployTokens(gid any, opt *ListGroupDeployTokensOptions, options ...RequestOptionFunc) ([]*DeployToken, *Response, error) {
	return do[[]*DeployToken](s.client,
		withPath("groups/%s/deploy_tokens", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployTokensService) GetGroupDeployToken(gid any, deployToken int64, options ...RequestOptionFunc) (*DeployToken, *Response, error) {
	return do[*DeployToken](s.client,
		withPath("groups/%s/deploy_tokens/%d", GroupID{gid}, deployToken),
		withRequestOpts(options...),
	)
}

// CreateGroupDeployTokenOptions represents the available CreateGroupDeployToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_tokens/#create-a-group-deploy-token
type CreateGroupDeployTokenOptions struct {
	Name      *string    `url:"name,omitempty" json:"name,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	Username  *string    `url:"username,omitempty" json:"username,omitempty"`
	Scopes    *[]string  `url:"scopes,omitempty" json:"scopes,omitempty"`
}

func (s *DeployTokensService) CreateGroupDeployToken(gid any, opt *CreateGroupDeployTokenOptions, options ...RequestOptionFunc) (*DeployToken, *Response, error) {
	return do[*DeployToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/deploy_tokens", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployTokensService) DeleteGroupDeployToken(gid any, deployToken int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/deploy_tokens/%d", GroupID{gid}, deployToken),
		withRequestOpts(options...),
	)
	return resp, err
}
