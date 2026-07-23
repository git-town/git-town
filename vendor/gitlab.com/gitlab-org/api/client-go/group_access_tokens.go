//
// Copyright 2022, Masahiro Yoshida
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
	// GroupAccessTokensServiceInterface defines all the API methods for the GroupAccessTokensService
	GroupAccessTokensServiceInterface interface {
		ListGroupAccessTokens(gid any, opt *ListGroupAccessTokensOptions, options ...RequestOptionFunc) ([]*GroupAccessToken, *Response, error)
		GetGroupAccessToken(gid any, id int64, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error)
		CreateGroupAccessToken(gid any, opt *CreateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error)
		RotateGroupAccessToken(gid any, id int64, opt *RotateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error)
		RotateGroupAccessTokenSelf(gid any, opt *RotateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error)
		RevokeGroupAccessToken(gid any, id int64, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupAccessTokensService handles communication with the
	// groups access tokens related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_access_tokens/
	GroupAccessTokensService struct {
		client *Client
	}
)

var _ GroupAccessTokensServiceInterface = (*GroupAccessTokensService)(nil)

// GroupAccessToken represents a GitLab group access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/
type GroupAccessToken resourceAccessToken

func (v GroupAccessToken) String() string {
	return Stringify(v)
}

// ListGroupAccessTokensOptions represents the available options for
// listing access tokens in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#list-all-group-access-tokens
type ListGroupAccessTokensOptions struct {
	ListOptions
	CreatedAfter   *ISOTime          `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore  *ISOTime          `url:"created_before,omitempty" json:"created_before,omitempty"`
	LastUsedAfter  *ISOTime          `url:"last_used_after,omitempty" json:"last_used_after,omitempty"`
	LastUsedBefore *ISOTime          `url:"last_used_before,omitempty" json:"last_used_before,omitempty"`
	Revoked        *bool             `url:"revoked,omitempty" json:"revoked,omitempty"`
	Search         *string           `url:"search,omitempty" json:"search,omitempty"`
	State          *AccessTokenState `url:"state,omitempty" json:"state,omitempty"`
	ExpiresAfter   *ISOTime          `url:"expires_after,omitempty" json:"expires_after,omitempty"`
	ExpiresBefore  *ISOTime          `url:"expires_before,omitempty" json:"expires_before,omitempty"`
	Sort           *AccessTokenSort  `url:"sort,omitempty" json:"sort,omitempty"`
}

// ListGroupAccessTokens gets a list of all group access tokens in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#list-all-group-access-tokens
func (s *GroupAccessTokensService) ListGroupAccessTokens(gid any, opt *ListGroupAccessTokensOptions, options ...RequestOptionFunc) ([]*GroupAccessToken, *Response, error) {
	return do[[]*GroupAccessToken](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/access_tokens", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupAccessToken gets a single group access tokens in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#get-details-on-a-group-access-token
func (s *GroupAccessTokensService) GetGroupAccessToken(gid any, id int64, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error) {
	return do[*GroupAccessToken](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/access_tokens/%d", GroupID{gid}, id),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// CreateGroupAccessTokenOptions represents the available CreateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#create-a-group-access-token
type CreateGroupAccessTokenOptions struct {
	Name        *string           `url:"name,omitempty" json:"name,omitempty"`
	Description *string           `url:"description,omitempty" json:"description,omitempty"`
	Scopes      *[]string         `url:"scopes,omitempty" json:"scopes,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	ExpiresAt   *ISOTime          `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// CreateGroupAccessToken creates a new group access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#create-a-group-access-token
func (s *GroupAccessTokensService) CreateGroupAccessToken(gid any, opt *CreateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error) {
	return do[*GroupAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/access_tokens", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RotateGroupAccessTokenOptions represents the available RotateGroupAccessToken()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#rotate-a-group-access-token
type RotateGroupAccessTokenOptions struct {
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// RotateGroupAccessToken revokes a group access token and returns a new group
// access token that expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#rotate-a-group-access-token
func (s *GroupAccessTokensService) RotateGroupAccessToken(gid any, id int64, opt *RotateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error) {
	return do[*GroupAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/access_tokens/%d/rotate", GroupID{gid}, id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RotateGroupAccessTokenSelf revokes the group access token used for the request
// and returns a new group access token that expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#self-rotate
func (s *GroupAccessTokensService) RotateGroupAccessTokenSelf(gid any, opt *RotateGroupAccessTokenOptions, options ...RequestOptionFunc) (*GroupAccessToken, *Response, error) {
	return do[*GroupAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/access_tokens/self/rotate", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RevokeGroupAccessToken revokes a group access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_access_tokens/#revoke-a-group-access-token
func (s *GroupAccessTokensService) RevokeGroupAccessToken(gid any, id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/access_tokens/%d", GroupID{gid}, id),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}
