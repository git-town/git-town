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
	FeatureFlagUserListsServiceInterface interface {
		ListFeatureFlagUserLists(pid any, opt *ListFeatureFlagUserListsOptions, options ...RequestOptionFunc) ([]*FeatureFlagUserList, *Response, error)
		CreateFeatureFlagUserList(pid any, opt *CreateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		GetFeatureFlagUserList(pid any, iid int64, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		UpdateFeatureFlagUserList(pid any, iid int64, opt *UpdateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error)
		DeleteFeatureFlagUserList(pid any, iid int64, options ...RequestOptionFunc) (*Response, error)
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
	ID        int64      `url:"id" json:"id"`
	IID       int64      `url:"iid" json:"iid"`
	ProjectID int64      `url:"project_id" json:"project_id"`
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
	return do[[]*FeatureFlagUserList](s.client,
		withPath("projects/%s/feature_flags_user_lists", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
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
	return do[*FeatureFlagUserList](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/feature_flags_user_lists", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetFeatureFlagUserList gets a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#get-a-feature-flag-user-list
func (s *FeatureFlagUserListsService) GetFeatureFlagUserList(pid any, iid int64, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error) {
	return do[*FeatureFlagUserList](s.client,
		withPath("projects/%s/feature_flags_user_lists/%d", ProjectID{pid}, iid),
		withRequestOpts(options...),
	)
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
func (s *FeatureFlagUserListsService) UpdateFeatureFlagUserList(pid any, iid int64, opt *UpdateFeatureFlagUserListOptions, options ...RequestOptionFunc) (*FeatureFlagUserList, *Response, error) {
	return do[*FeatureFlagUserList](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/feature_flags_user_lists/%d", ProjectID{pid}, iid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteFeatureFlagUserList deletes a feature flag user list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flag_user_lists/#delete-feature-flag-user-list
func (s *FeatureFlagUserListsService) DeleteFeatureFlagUserList(pid any, iid int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/feature_flags_user_lists/%d", ProjectID{pid}, iid),
		withRequestOpts(options...),
	)
	return resp, err
}
