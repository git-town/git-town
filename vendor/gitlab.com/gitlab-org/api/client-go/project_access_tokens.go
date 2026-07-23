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
)

type (
	ProjectAccessTokensServiceInterface interface {
		ListProjectAccessTokens(pid any, opt *ListProjectAccessTokensOptions, options ...RequestOptionFunc) ([]*ProjectAccessToken, *Response, error)
		GetProjectAccessToken(pid any, id int64, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error)
		CreateProjectAccessToken(pid any, opt *CreateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error)
		RotateProjectAccessToken(pid any, id int64, opt *RotateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error)
		RotateProjectAccessTokenSelf(pid any, opt *RotateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error)
		RevokeProjectAccessToken(pid any, id int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectAccessTokensService handles communication with the
	// project access tokens related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_access_tokens/
	ProjectAccessTokensService struct {
		client *Client
	}
)

var _ ProjectAccessTokensServiceInterface = (*ProjectAccessTokensService)(nil)

// ProjectAccessToken represents a GitLab project access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/
type ProjectAccessToken resourceAccessToken

func (v ProjectAccessToken) String() string {
	return Stringify(v)
}

// ListProjectAccessTokensOptions represents the available
// ListProjectAccessTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#list-all-project-access-tokens
type ListProjectAccessTokensOptions struct {
	ListOptions
	State *string `url:"state,omitempty" json:"state,omitempty"`
}

// ListProjectAccessTokens gets a list of all project access tokens in a
// project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#list-all-project-access-tokens
func (s *ProjectAccessTokensService) ListProjectAccessTokens(pid any, opt *ListProjectAccessTokensOptions, options ...RequestOptionFunc) ([]*ProjectAccessToken, *Response, error) {
	return do[[]*ProjectAccessToken](s.client,
		withPath("projects/%s/access_tokens", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProjectAccessToken gets a single project access tokens in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#get-details-on-a-project-access-token
func (s *ProjectAccessTokensService) GetProjectAccessToken(pid any, id int64, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error) {
	return do[*ProjectAccessToken](s.client,
		withPath("projects/%s/access_tokens/%d", ProjectID{pid}, id),
		withRequestOpts(options...),
	)
}

// CreateProjectAccessTokenOptions represents the available CreateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#create-a-project-access-token
type CreateProjectAccessTokenOptions struct {
	Name        *string           `url:"name,omitempty" json:"name,omitempty"`
	Description *string           `url:"description,omitempty" json:"description,omitempty"`
	Scopes      *[]string         `url:"scopes,omitempty" json:"scopes,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ExpiresAt   *ISOTime          `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// CreateProjectAccessToken creates a new project access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#create-a-project-access-token
func (s *ProjectAccessTokensService) CreateProjectAccessToken(pid any, opt *CreateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error) {
	return do[*ProjectAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/access_tokens", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RotateProjectAccessTokenOptions represents the available RotateProjectAccessToken()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#rotate-a-project-access-token
type RotateProjectAccessTokenOptions struct {
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// RotateProjectAccessToken revokes a project access token and returns a new
// project access token that expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#rotate-a-project-access-token
func (s *ProjectAccessTokensService) RotateProjectAccessToken(pid any, id int64, opt *RotateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error) {
	return do[*ProjectAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/access_tokens/%d/rotate", ProjectID{pid}, id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RotateProjectAccessTokenSelf revokes the project access token used for the request
// and returns a new project access token that expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#self-rotate
func (s *ProjectAccessTokensService) RotateProjectAccessTokenSelf(pid any, opt *RotateProjectAccessTokenOptions, options ...RequestOptionFunc) (*ProjectAccessToken, *Response, error) {
	return do[*ProjectAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/access_tokens/self/rotate", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RevokeProjectAccessToken revokes a project access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_access_tokens/#revoke-a-project-access-token
func (s *ProjectAccessTokensService) RevokeProjectAccessToken(pid any, id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/access_tokens/%d", ProjectID{pid}, id),
		withRequestOpts(options...),
	)
	return resp, err
}
