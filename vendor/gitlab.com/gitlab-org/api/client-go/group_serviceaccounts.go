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
	"fmt"
	"net/http"
	"time"
)

// GroupServiceAccount represents a GitLab service account user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#create-a-service-account-user
type GroupServiceAccount struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

// ListServiceAccountsOptions represents the available ListServiceAccounts() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#list-all-service-account-users
type ListServiceAccountsOptions struct {
	ListOptions
	OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
}

// ListServiceAccounts gets a list of service accounts.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#list-all-service-account-users
func (s *GroupsService) ListServiceAccounts(gid any, opt *ListServiceAccountsOptions, options ...RequestOptionFunc) ([]*GroupServiceAccount, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var sa []*GroupServiceAccount
	resp, err := s.client.Do(req, &sa)
	if err != nil {
		return nil, resp, err
	}

	return sa, resp, nil
}

// CreateServiceAccountOptions represents the available CreateServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#create-a-service-account-user
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
// https://docs.gitlab.com/api/group_service_accounts/#create-a-service-account-user
func (s *GroupsService) CreateServiceAccount(gid any, opt *CreateServiceAccountOptions, options ...RequestOptionFunc) (*GroupServiceAccount, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	sa := new(GroupServiceAccount)
	resp, err := s.client.Do(req, sa)
	if err != nil {
		return nil, resp, err
	}

	return sa, resp, nil
}

// UpdateServiceAccountOptions represents the available UpdateServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#update-a-service-account-user
type UpdateServiceAccountOptions struct {
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
	Username *string `url:"username,omitempty" json:"username,omitempty"`
}

// UpdateServiceAccount updates a service account user.
//
// This API endpoint works on top-level groups only. It does not work on subgroups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#update-a-service-account-user
func (s *GroupsService) UpdateServiceAccount(gid any, serviceAccount int, opt *UpdateServiceAccountOptions, options ...RequestOptionFunc) (*GroupServiceAccount, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d", PathEscape(group), serviceAccount)

	req, err := s.client.NewRequest(http.MethodPatch, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	sa := new(GroupServiceAccount)
	resp, err := s.client.Do(req, sa)
	if err != nil {
		return nil, resp, err
	}

	return sa, resp, nil
}

// DeleteServiceAccountOptions represents the available DeleteServiceAccount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#delete-a-service-account-user
type DeleteServiceAccountOptions struct {
	HardDelete *bool `url:"hard_delete,omitempty" json:"hard_delete,omitempty"`
}

// DeleteServiceAccount Deletes a service account user.
//
// This API endpoint works on top-level groups only. It does not work on subgroups.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#delete-a-service-account-user
func (s *GroupsService) DeleteServiceAccount(gid any, serviceAccount int, opt *DeleteServiceAccountOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d", PathEscape(group), serviceAccount)

	req, err := s.client.NewRequest(http.MethodDelete, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListServiceAccountPersonalAccessTokensOptions represents the available
// ListServiceAccountPersonalAccessTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#list-all-personal-access-tokens-for-a-service-account-user
type ListServiceAccountPersonalAccessTokensOptions struct {
	ListOptions
	CreatedAfter   *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore  *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
	ExpiresAfter   *ISOTime   `url:"expires_after,omitempty" json:"expires_after,omitempty"`
	ExpiresBefore  *ISOTime   `url:"expires_before,omitempty" json:"expires_before,omitempty"`
	LastUsedAfter  *time.Time `url:"last_used_after,omitempty" json:"last_used_after,omitempty"`
	LastUsedBefore *time.Time `url:"last_used_before,omitempty" json:"last_used_before,omitempty"`
	Revoked        *bool      `url:"revoked,omitempty" json:"revoked,omitempty"`
	UserID         *int       `url:"user_id,omitempty" json:"user_id,omitempty"`
	Search         *string    `url:"search,omitempty" json:"search,omitempty"`
	Sort           *string    `url:"sort,omitempty" json:"sort,omitempty"`
	State          *string    `url:"state,omitempty" json:"state,omitempty"`
}

// ListServiceAccountPersonalAccessTokens gets a list of personal access tokens for a
// service account user for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#list-all-personal-access-tokens-for-a-service-account-user
func (s *GroupsService) ListServiceAccountPersonalAccessTokens(gid any, serviceAccount int, opt *ListServiceAccountPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*PersonalAccessToken, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d/personal_access_tokens", PathEscape(group), serviceAccount)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var pats []*PersonalAccessToken
	resp, err := s.client.Do(req, &pats)
	if err != nil {
		return nil, resp, err
	}

	return pats, resp, nil
}

// CreateServiceAccountPersonalAccessTokenOptions represents the available
// CreateServiceAccountPersonalAccessToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#create-a-personal-access-token-for-a-service-account-user
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
// https://docs.gitlab.com/api/group_service_accounts/#create-a-personal-access-token-for-a-service-account-user
func (s *GroupsService) CreateServiceAccountPersonalAccessToken(gid any, serviceAccount int, opt *CreateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d/personal_access_tokens", PathEscape(group), serviceAccount)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	pat := new(PersonalAccessToken)
	resp, err := s.client.Do(req, pat)
	if err != nil {
		return nil, resp, err
	}

	return pat, resp, nil
}

// RevokeServiceAccountPersonalAccessToken revokes a personal access token for an
// existing service account user in a given top-level group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#revoke-a-personal-access-token-for-a-service-account-user
func (s *GroupsService) RevokeServiceAccountPersonalAccessToken(gid any, serviceAccount, token int, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d/personal_access_tokens/%d", PathEscape(group), serviceAccount, token)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RotateServiceAccountPersonalAccessTokenOptions represents the available RotateServiceAccountPersonalAccessToken()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#rotate-a-personal-access-token-for-a-service-account-user
type RotateServiceAccountPersonalAccessTokenOptions struct {
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// RotateServiceAccountPersonalAccessToken rotates a Personal Access Token for a
// service account user for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_service_accounts/#rotate-a-personal-access-token-for-a-service-account-user
func (s *GroupsService) RotateServiceAccountPersonalAccessToken(gid any, serviceAccount, token int, opt *RotateServiceAccountPersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/service_accounts/%d/personal_access_tokens/%d/rotate", PathEscape(group), serviceAccount, token)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	pat := new(PersonalAccessToken)
	resp, err := s.client.Do(req, pat)
	if err != nil {
		return nil, resp, err
	}

	return pat, resp, nil
}
