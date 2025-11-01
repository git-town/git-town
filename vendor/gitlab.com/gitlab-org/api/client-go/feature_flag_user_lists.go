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
	FeatureFlagUserListsServiceInterface interface {
		ListFeatureFlagUserLists(pid any, opt *ListFeatureFlagUserListsOptions, options ...RequestOptionFunc) ([]*FeatureFlagUserList, *Response, error)
		CreateFeatureFlagUserList(pid any, opt *CreateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		GetFeatureFlagUserList(pid any, iid int, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		UpdateFeatureFlagUserList(pid any, iid int, opt *UpdateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		DeleteFeatureFlagUserList(pid any, iid int, options ...RequestOptionFunc) (*Response, error)
	}

	// FeatureFlagUserListsService handles communication with the feature flag
	// user list related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/feature_flag_user_lists/
	FeatureFlagUserListsService struct {
		client *Client
	}
)

var _ FeatureFlagUserListsServiceInterface = (*FeatureFlagUserListsService)(nil)

// FeatureFlagUserList represents a project feature flag user list.
//
// GitLab API docs: https://docs.gitlab.com/api/feature_flag_user_lists/
type FeatureFlagUserList struct {
	Name      string     `url:"name" json:"name"`
	UserXIDs  string     `url:"user_xids" json:"user_xids"`
	ID        int        `url:"id" json:"id"`
	IID       int        `url:"iid" json:"iid"`
	ProjectID int        `url:"project_id" json:"project_id"`
	CreatedAt *time.Time `url:"created_at" json:"created_at"`
	UpdatedAt *time.Time `url:"updated_at" json:"updated_at"`
}

// ListFeatureFlagUserListsOptions represents the available
// ListFeatureFlagUserLists() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#list-all-feature-flag-user-lists-for-a-project
type ListFeatureFlagUserListsOptions struct {
	ListOptions
	Search string `url:"search,omitempty" json:"search,omitempty"`
}

// ListFeatureFlagUserLists gets all feature flag user lists for the requested
// project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#list-all-feature-flag-user-lists-for-a-project
func (s *FeatureFlagUserListsService) ListFeatureFlagUserLists(pid any, opt *ListFeatureFlagUserListsOptions, options ...RequestOptionFunc) ([]*FeatureFlagUserList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/feature_flags_user_lists", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var lists []*FeatureFlagUserList
	resp, err := s.client.Do(req, &lists)
	if err != nil {
		return nil, resp, err
	}

	return lists, resp, nil
}

// CreateFeatureFlagUserListOptions represents the available
// CreateFeatureFlagUserList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#create-a-feature-flag-user-list
type CreateFeatureFlagUserListOptions struct {
	Name     string `url:"name,omitempty" json:"name,omitempty"`
	UserXIDs string `url:"user_xids,omitempty" json:"user_xids,omitempty"`
}

// CreateFeatureFlagUserList creates a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#create-a-feature-flag-user-list
func (s *FeatureFlagUserListsService) CreateFeatureFlagUserList(pid any, opt *CreateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/feature_flags_user_lists", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	list := new(FeatureFlagUserList)
	resp, err := s.client.Do(req, list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetFeatureFlagUserList gets a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#get-a-feature-flag-user-list
func (s *FeatureFlagUserListsService) GetFeatureFlagUserList(pid any, iid int, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/feature_flags_user_lists/%d", PathEscape(project), iid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	list := new(FeatureFlagUserList)
	resp, err := s.client.Do(req, list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// UpdateFeatureFlagUserListOptions represents the available
// UpdateFeatureFlagUserList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#update-a-feature-flag-user-list
type UpdateFeatureFlagUserListOptions struct {
	Name     string `url:"name,omitempty" json:"name,omitempty"`
	UserXIDs string `url:"user_xids,omitempty" json:"user_xids,omitempty"`
}

// UpdateFeatureFlagUserList updates a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#update-a-feature-flag-user-list
func (s *FeatureFlagUserListsService) UpdateFeatureFlagUserList(pid any, iid int, opt *UpdateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/feature_flags_user_lists/%d", PathEscape(project), iid)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	list := new(FeatureFlagUserList)
	resp, err := s.client.Do(req, list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// DeleteFeatureFlagUserList deletes a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#delete-feature-flag-user-list
func (s *FeatureFlagUserListsService) DeleteFeatureFlagUserList(pid any, iid int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/feature_flags_user_lists/%d", PathEscape(project), iid)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
