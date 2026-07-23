//
// Copyright 2023, James Hong
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

// GroupServiceAccount represents a GitLab service account user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#create-a-group-service-account
type GroupServiceAccount struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

// ListServiceAccountsOptions represents the available ListServiceAccounts() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#list-all-group-service-accounts
type ListServiceAccountsOptions struct {
	ListOptions
	OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
}

// ListServiceAccounts gets a list of service accounts.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#list-all-group-service-accounts
func (s *GroupsService) ListServiceAccounts(gid any, opt *ListServiceAccountsOptions, options ...RequestOptionFunc) ([]*GroupServiceAccount, *Response, error) {
	return do[[]*GroupServiceAccount](s.client,
		withPath("groups/%s/service_accounts", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateServiceAccountOptions represents the available CreateServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#create-a-group-service-account
type CreateServiceAccountOptions struct {
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
	Username *string `url:"username,omitempty" json:"username,omitempty"`
	Email    *string `url:"email,omitempty" json:"email,omitempty"`
}

// CreateServiceAccount creates a service account user.
//
// This API endpoint works on top-level groups only. It does not work on subgroups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#create-a-group-service-account
func (s *GroupsService) CreateServiceAccount(gid any, opt *CreateServiceAccountOptions, options ...RequestOptionFunc) (*GroupServiceAccount, *Response, error) {
	return do[*GroupServiceAccount](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/service_accounts", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateServiceAccountOptions represents the available UpdateServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#update-a-group-service-account
type UpdateServiceAccountOptions struct {
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
	Username *string `url:"username,omitempty" json:"username,omitempty"`
	Email    *string `url:"email,omitempty" json:"email,omitempty"`
}

// UpdateServiceAccount updates a service account user.
//
// This API endpoint works on top-level groups only. It does not work on subgroups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#update-a-group-service-account
func (s *GroupsService) UpdateServiceAccount(gid any, serviceAccount int64, opt *UpdateServiceAccountOptions, options ...RequestOptionFunc) (*GroupServiceAccount, *Response, error) {
	return do[*GroupServiceAccount](s.client,
		withMethod(http.MethodPatch),
		withPath("groups/%s/service_accounts/%d", GroupID{gid}, serviceAccount),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteServiceAccountOptions represents the available DeleteServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#delete-a-group-service-account
type DeleteServiceAccountOptions struct {
	HardDelete *bool `url:"hard_delete,omitempty" json:"hard_delete,omitempty"`
}

// DeleteServiceAccount Deletes a service account user.
//
// This API endpoint works on top-level groups only. It does not work on subgroups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#delete-a-group-service-account
func (s *GroupsService) DeleteServiceAccount(gid any, serviceAccount int64, opt *DeleteServiceAccountOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/service_accounts/%d", GroupID{gid}, serviceAccount),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListServiceAccountPersonalAccessTokensOptions represents the available
// ListServiceAccountPersonalAccessTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#list-all-personal-access-tokens-for-a-group-service-account
type ListServiceAccountPersonalAccessTokensOptions struct {
	ListOptions
	CreatedAfter   *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore  *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
	ExpiresAfter   *ISOTime   `url:"expires_after,omitempty" json:"expires_after,omitempty"`
	ExpiresBefore  *ISOTime   `url:"expires_before,omitempty" json:"expires_before,omitempty"`
	LastUsedAfter  *time.Time `url:"last_used_after,omitempty" json:"last_used_after,omitempty"`
	LastUsedBefore *time.Time `url:"last_used_before,omitempty" json:"last_used_before,omitempty"`
	Revoked        *bool      `url:"revoked,omitempty" json:"revoked,omitempty"`
	UserID         *int64     `url:"user_id,omitempty" json:"user_id,omitempty"`
	Search         *string    `url:"search,omitempty" json:"search,omitempty"`
	Sort           *string    `url:"sort,omitempty" json:"sort,omitempty"`
	State          *string    `url:"state,omitempty" json:"state,omitempty"`
}

// ListServiceAccountPersonalAccessTokens gets a list of personal access tokens for a
// service account user for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#list-all-personal-access-tokens-for-a-group-service-account
func (s *GroupsService) ListServiceAccountPersonalAccessTokens(gid any, serviceAccount int64, opt *ListServiceAccountPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*PersonalAccessToken, *Response, error) {
	return do[[]*PersonalAccessToken](s.client,
		withPath("groups/%s/service_accounts/%d/personal_access_tokens", GroupID{gid}, serviceAccount),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateServiceAccountPersonalAccessTokenOptions represents the available
// CreateServiceAccountPersonalAccessToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#create-a-personal-access-token-for-a-group-service-account
type CreateServiceAccountPersonalAccessTokenOptions struct {
	Name        *string   `url:"name,omitempty" json:"name,omitempty"`
	Description *string   `url:"description,omitempty" json:"description,omitempty"`
	Scopes      *[]string `url:"scopes,omitempty" json:"scopes,omitempty"`
	ExpiresAt   *ISOTime  `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// CreateServiceAccountPersonalAccessToken add a new Personal Access Token for a
// service account user for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#create-a-personal-access-token-for-a-group-service-account
func (s *GroupsService) CreateServiceAccountPersonalAccessToken(gid any, serviceAccount int64, opt *CreateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/service_accounts/%d/personal_access_tokens", GroupID{gid}, serviceAccount),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RevokeServiceAccountPersonalAccessToken revokes a personal access token for an
// existing service account user in a given top-level group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#revoke-a-personal-access-token-for-a-group-service-account
func (s *GroupsService) RevokeServiceAccountPersonalAccessToken(gid any, serviceAccount, token int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/service_accounts/%d/personal_access_tokens/%d", GroupID{gid}, serviceAccount, token),
		withRequestOpts(options...),
	)
	return resp, err
}

// RotateServiceAccountPersonalAccessTokenOptions represents the available RotateServiceAccountPersonalAccessToken()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#rotate-a-personal-access-token-for-a-group-service-account
type RotateServiceAccountPersonalAccessTokenOptions struct {
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// RotateServiceAccountPersonalAccessToken rotates a Personal Access Token for a
// service account user for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/service_accounts/#rotate-a-personal-access-token-for-a-group-service-account
func (s *GroupsService) RotateServiceAccountPersonalAccessToken(gid any, serviceAccount, token int64, opt *RotateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/service_accounts/%d/personal_access_tokens/%d/rotate", GroupID{gid}, serviceAccount, token),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
