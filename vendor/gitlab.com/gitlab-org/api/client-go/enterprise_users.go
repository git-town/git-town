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

type (
	EnterpriseUsersServiceInterface interface {
		ListEnterpriseUsers(gid any, opt *ListEnterpriseUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
		GetEnterpriseUser(gid any, uid int, options ...RequestOptionFunc) (*User, *Response, error)
		Disable2FAForEnterpriseUser(gid any, uid int, options ...RequestOptionFunc) (*Response, error)
	}

	// EnterpriseUsersService handles communication with the enterprise users
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_enterprise_users/
	EnterpriseUsersService struct {
		client *Client
	}
)

var _ EnterpriseUsersServiceInterface = (*EnterpriseUsersService)(nil)

// ListEnterpriseUsersOptions represents the available
// ListEnterpriseUsers() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_enterprise_users/#list-all-enterprise-users
type ListEnterpriseUsersOptions struct {
	ListOptions
	Username      string     `url:"username,omitempty" json:"username,omitempty"`
	Search        string     `url:"search,omitempty" json:"search,omitempty"`
	Active        bool       `url:"active,omitempty" json:"active,omitempty"`
	Blocked       bool       `url:"blocked,omitempty" json:"blocked,omitempty"`
	CreatedAfter  *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
	TwoFactor     string     `url:"two_factor,omitempty" json:"two_factor,omitempty"`
}

// ListEnterpriseUsers lists all enterprise users for a given top-level group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_enterprise_users/#list-all-enterprise-users
func (s *EnterpriseUsersService) ListEnterpriseUsers(gid any, opt *ListEnterpriseUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/enterprise_users", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// GetEnterpriseUser gets details on a specified enterprise user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_enterprise_users/#get-details-on-an-enterprise-user
func (s *EnterpriseUsersService) GetEnterpriseUser(gid any, uid int, options ...RequestOptionFunc) (*User, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/enterprise_users/%d", PathEscape(group), uid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// Disable2FAForEnterpriseUser disables two-factor authentication (2FA) for a
// specified enterprise user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_enterprise_users/#disable-two-factor-authentication-for-an-enterprise-user
func (s *EnterpriseUsersService) Disable2FAForEnterpriseUser(gid any, uid int, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/enterprise_users/%d/disable_two_factor", PathEscape(group), uid)

	req, err := s.client.NewRequest(http.MethodPatch, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
