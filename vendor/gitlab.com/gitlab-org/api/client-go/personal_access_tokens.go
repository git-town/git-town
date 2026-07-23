//
// Copyright 2022, Ryan Glab <ryan.j.glab@gmail.com>
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
	PersonalAccessTokensServiceInterface interface {
		ListPersonalAccessTokens(opt *ListPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*PersonalAccessToken, *Response, error)
		GetSinglePersonalAccessTokenByID(token int64, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		GetSinglePersonalAccessToken(options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		RotatePersonalAccessToken(token int64, opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		RotatePersonalAccessTokenByID(token int64, opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		RotatePersonalAccessTokenSelf(opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		// Deprecated: to be removed in 2.0; use RevokePersonalAccessTokenByID instead
		RevokePersonalAccessToken(token int64, options ...RequestOptionFunc) (*Response, error)
		RevokePersonalAccessTokenByID(token int64, options ...RequestOptionFunc) (*Response, error)
		RevokePersonalAccessTokenSelf(options ...RequestOptionFunc) (*Response, error)
	}

	// PersonalAccessTokensService handles communication with the personal access
	// tokens related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/personal_access_tokens/
	PersonalAccessTokensService struct {
		client *Client
	}
)

var _ PersonalAccessTokensServiceInterface = (*PersonalAccessTokensService)(nil)

// PersonalAccessToken represents a personal access token.
//
// GitLab API docs: https://docs.gitlab.com/api/personal_access_tokens/
type PersonalAccessToken struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Revoked     bool       `json:"revoked"`
	CreatedAt   *time.Time `json:"created_at"`
	Description string     `json:"description"`
	Scopes      []string   `json:"scopes"`
	UserID      int64      `json:"user_id"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Active      bool       `json:"active"`
	ExpiresAt   *ISOTime   `json:"expires_at"`
	Token       string     `json:"token,omitempty"`
}

// ResourceAccessToken represents a generic access token used for both
// project and group access tokens. It's only used as an alias type, which
// is why it's not exported.
type resourceAccessToken struct {
	PersonalAccessToken
	AccessLevel AccessLevelValue `json:"access_level"`
}

func (p PersonalAccessToken) String() string {
	return Stringify(p)
}

// ListPersonalAccessTokensOptions represents the available
// ListPersonalAccessTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#list-all-personal-access-tokens
type ListPersonalAccessTokensOptions struct {
	ListOptions
	CreatedAfter   *ISOTime `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore  *ISOTime `url:"created_before,omitempty" json:"created_before,omitempty"`
	ExpiresAfter   *ISOTime `url:"expires_after,omitempty" json:"expires_after,omitempty"`
	ExpiresBefore  *ISOTime `url:"expires_before,omitempty" json:"expires_before,omitempty"`
	LastUsedAfter  *ISOTime `url:"last_used_after,omitempty" json:"last_used_after,omitempty"`
	LastUsedBefore *ISOTime `url:"last_used_before,omitempty" json:"last_used_before,omitempty"`
	Revoked        *bool    `url:"revoked,omitempty" json:"revoked,omitempty"`
	Search         *string  `url:"search,omitempty" json:"search,omitempty"`
	Sort           *string  `url:"sort,omitempty" json:"sort,omitempty"`
	State          *string  `url:"state,omitempty" json:"state,omitempty"`
	UserID         *int64   `url:"user_id,omitempty" json:"user_id,omitempty"`
}

// ListPersonalAccessTokens gets a list of all personal access tokens.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#list-all-personal-access-tokens
func (s *PersonalAccessTokensService) ListPersonalAccessTokens(opt *ListPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*PersonalAccessToken, *Response, error) {
	return do[[]*PersonalAccessToken](s.client,
		withPath("personal_access_tokens"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetSinglePersonalAccessTokenByID get a single personal access token by its ID.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#get-details-on-a-personal-access-token
func (s *PersonalAccessTokensService) GetSinglePersonalAccessTokenByID(token int64, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withPath("personal_access_tokens/%d", token),
		withRequestOpts(options...),
	)
}

// GetSinglePersonalAccessToken get a single personal access token by using
// passing the token in a header.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#self-inform
func (s *PersonalAccessTokensService) GetSinglePersonalAccessToken(options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withPath("personal_access_tokens/self"),
		withRequestOpts(options...),
	)
}

// RotatePersonalAccessTokenOptions represents the available RotatePersonalAccessToken()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#rotate-a-personal-access-token
type RotatePersonalAccessTokenOptions struct {
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// RotatePersonalAccessToken is a backwards-compat shim for RotatePersonalAccessTokenByID.
func (s *PersonalAccessTokensService) RotatePersonalAccessToken(token int64, opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return s.RotatePersonalAccessTokenByID(token, opt, options...)
}

// RotatePersonalAccessTokenByID revokes a token and returns a new token that
// expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#rotate-a-personal-access-token
func (s *PersonalAccessTokensService) RotatePersonalAccessTokenByID(token int64, opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("personal_access_tokens/%d/rotate", token),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RotatePersonalAccessTokenSelf revokes the currently authenticated token
// and returns a new token that expires in one week per default.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#self-rotate
func (s *PersonalAccessTokensService) RotatePersonalAccessTokenSelf(opt *RotatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("personal_access_tokens/self/rotate"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RevokePersonalAccessToken is a backwards-compat shim for RevokePersonalAccessTokenByID.
// Deprecated: to be removed in 2.0; use RevokePersonalAccessTokenByID instead
func (s *PersonalAccessTokensService) RevokePersonalAccessToken(token int64, options ...RequestOptionFunc) (*Response, error) {
	return s.RevokePersonalAccessTokenByID(token, options...)
}

// RevokePersonalAccessTokenByID revokes a personal access token by its ID.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#revoke-a-personal-access-token
func (s *PersonalAccessTokensService) RevokePersonalAccessTokenByID(token int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("personal_access_tokens/%d", token),
		withRequestOpts(options...),
	)
	return resp, err
}

// RevokePersonalAccessTokenSelf revokes the currently authenticated
// personal access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/personal_access_tokens/#self-revoke
func (s *PersonalAccessTokensService) RevokePersonalAccessTokenSelf(options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("personal_access_tokens/self"),
		withRequestOpts(options...),
	)
	return resp, err
}
