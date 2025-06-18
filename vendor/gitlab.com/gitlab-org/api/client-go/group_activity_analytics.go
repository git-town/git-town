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

import "net/http"

type (
	GroupActivityAnalyticsServiceInterface interface {
		GetRecentlyCreatedIssuesCount(opt *GetRecentlyCreatedIssuesCountOptions, options ...RequestOptionFunc) (*IssuesCount, *Response, error)
		GetRecentlyCreatedMergeRequestsCount(opt *GetRecentlyCreatedMergeRequestsCountOptions, options ...RequestOptionFunc) (*MergeRequestsCount, *Response, error)
		GetRecentlyAddedMembersCount(opt *GetRecentlyAddedMembersCountOptions, options ...RequestOptionFunc) (*NewMembersCount, *Response, error)
	}

	// GroupActivityAnalyticsService handles communication with the group activity
	// analytics related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_activity_analytics/
	GroupActivityAnalyticsService struct {
		client *Client
	}
)

var _ GroupActivityAnalyticsServiceInterface = (*GroupActivityAnalyticsService)(nil)

// IssuesCount represents the total count of recently created issues in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-issues-for-group
type IssuesCount struct {
	IssuesCount int `url:"issues_count" json:"issues_count"`
}

// GetRecentlyCreatedIssuesCountOptions represents the available
// GetRecentlyCreatedIssuesCount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-issues-for-group
type GetRecentlyCreatedIssuesCountOptions struct {
	GroupPath string `url:"group_path" json:"group_path"`
}

// GetRecentlyCreatedIssuesCount gets the count of recently created issues for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-issues-for-group
func (s *GroupActivityAnalyticsService) GetRecentlyCreatedIssuesCount(opt *GetRecentlyCreatedIssuesCountOptions, options ...RequestOptionFunc) (*IssuesCount, *Response, error) {
	u := "analytics/group_activity/issues_count"
	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	res := new(IssuesCount)
	resp, err := s.client.Do(req, res)
	if err != nil {
		return nil, resp, err
	}

	return res, resp, nil
}

// MergeRequestsCount represents the total count of recently created merge requests
// in a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-merge-requests-for-group
type MergeRequestsCount struct {
	MergeRequestsCount int `url:"merge_requests_count" json:"merge_requests_count"`
}

// GetRecentlyCreatedMergeRequestsCountOptions represents the available
// GetRecentlyCreatedMergeRequestsCount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-merge-requests-for-group
type GetRecentlyCreatedMergeRequestsCountOptions struct {
	GroupPath string `url:"group_path" json:"group_path"`
}

// GetRecentlyCreatedMergeRequestsCount gets the count of recently created merge
// requests for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-recently-created-merge-requests-for-group
func (s *GroupActivityAnalyticsService) GetRecentlyCreatedMergeRequestsCount(opt *GetRecentlyCreatedMergeRequestsCountOptions, options ...RequestOptionFunc) (*MergeRequestsCount, *Response, error) {
	u := "analytics/group_activity/merge_requests_count"
	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	res := new(MergeRequestsCount)
	resp, err := s.client.Do(req, res)
	if err != nil {
		return nil, resp, err
	}

	return res, resp, nil
}

// NewMembersCount represents the total count of recently added members to a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-members-recently-added-to-group
type NewMembersCount struct {
	NewMembersCount int `url:"new_members_count" json:"new_members_count"`
}

// GetRecentlyAddedMembersCountOptions represents the available
// GetRecentlyAddedMembersCount() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-members-recently-added-to-group
type GetRecentlyAddedMembersCountOptions struct {
	GroupPath string `url:"group_path" json:"group_path"`
}

// GetRecentlyAddedMembersCount gets the count of recently added members to a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_activity_analytics/#get-count-of-members-recently-added-to-group
func (s *GroupActivityAnalyticsService) GetRecentlyAddedMembersCount(opt *GetRecentlyAddedMembersCountOptions, options ...RequestOptionFunc) (*NewMembersCount, *Response, error) {
	u := "analytics/group_activity/new_members_count"
	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	res := new(NewMembersCount)
	resp, err := s.client.Do(req, res)
	if err != nil {
		return nil, resp, err
	}

	return res, resp, nil
}
