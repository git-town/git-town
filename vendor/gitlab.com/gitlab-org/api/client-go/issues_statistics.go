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
	"time"
)

type (
	IssuesStatisticsServiceInterface interface {
		// GetIssuesStatistics gets issues statistics on all issues the authenticated
		// user has access to.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/issues_statistics/#get-issues-statistics
		GetIssuesStatistics(opt *GetIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error)
		// GetGroupIssuesStatistics gets issues count statistics for given group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/issues_statistics/#get-group-issues-statistics
		GetGroupIssuesStatistics(gid any, opt *GetGroupIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error)
		// GetProjectIssuesStatistics gets issues count statistics for given project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/issues_statistics/#get-project-issues-statistics
		GetProjectIssuesStatistics(pid any, opt *GetProjectIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error)
	}

	// IssuesStatisticsService handles communication with the issues statistics
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/issues_statistics/
	IssuesStatisticsService struct {
		client *Client
	}
)

var _ IssuesStatisticsServiceInterface = (*IssuesStatisticsService)(nil)

// IssuesStatistics represents a GitLab issues statistic.
//
// GitLab API docs: https://docs.gitlab.com/api/issues_statistics/
type IssuesStatistics struct {
	Statistics IssuesStatisticsStatistics `json:"statistics"`
}

func (n IssuesStatistics) String() string {
	return Stringify(n)
}

// IssuesStatisticsStatistics represents a GitLab issues statistic statistics.
//
// GitLab API docs: https://docs.gitlab.com/api/issues_statistics/
type IssuesStatisticsStatistics struct {
	Counts IssuesStatisticsCounts `json:"counts"`
}

// IssuesStatisticsCounts represents a GitLab issues statistic counts.
//
// GitLab API docs: https://docs.gitlab.com/api/issues_statistics/
type IssuesStatisticsCounts struct {
	All    int64 `json:"all"`
	Closed int64 `json:"closed"`
	Opened int64 `json:"opened"`
}

// GetIssuesStatisticsOptions represents the available GetIssuesStatistics() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issues_statistics/#get-issues-statistics
type GetIssuesStatisticsOptions struct {
	Labels           *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	Milestone        *string       `url:"milestone,omitempty" json:"milestone,omitempty"`
	Scope            *string       `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID         *int64        `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername   *string       `url:"author_username,omitempty" json:"author_username,omitempty"`
	AssigneeID       *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	AssigneeUsername *[]string     `url:"assignee_username,omitempty" json:"assignee_username,omitempty"`
	MyReactionEmoji  *string       `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	IIDs             *[]int64      `url:"iids[],omitempty" json:"iids,omitempty"`
	Search           *string       `url:"search,omitempty" json:"search,omitempty"`
	In               *string       `url:"in,omitempty" json:"in,omitempty"`
	CreatedAfter     *time.Time    `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore    *time.Time    `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter     *time.Time    `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore    *time.Time    `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Confidential     *bool         `url:"confidential,omitempty" json:"confidential,omitempty"`
}

func (s *IssuesStatisticsService) GetIssuesStatistics(opt *GetIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error) {
	return do[*IssuesStatistics](s.client,
		withPath("issues_statistics"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupIssuesStatisticsOptions represents the available GetGroupIssuesStatistics()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issues_statistics/#get-group-issues-statistics
type GetGroupIssuesStatisticsOptions struct {
	Labels           *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	IIDs             *[]int64      `url:"iids[],omitempty" json:"iids,omitempty"`
	Milestone        *string       `url:"milestone,omitempty" json:"milestone,omitempty"`
	Scope            *string       `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID         *int64        `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername   *string       `url:"author_username,omitempty" json:"author_username,omitempty"`
	AssigneeID       *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	AssigneeUsername *[]string     `url:"assignee_username,omitempty" json:"assignee_username,omitempty"`
	MyReactionEmoji  *string       `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	Search           *string       `url:"search,omitempty" json:"search,omitempty"`
	CreatedAfter     *time.Time    `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore    *time.Time    `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter     *time.Time    `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore    *time.Time    `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Confidential     *bool         `url:"confidential,omitempty" json:"confidential,omitempty"`
}

func (s *IssuesStatisticsService) GetGroupIssuesStatistics(gid any, opt *GetGroupIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error) {
	return do[*IssuesStatistics](s.client,
		withPath("groups/%s/issues_statistics", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProjectIssuesStatisticsOptions represents the available
// GetProjectIssuesStatistics() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/issues_statistics/#get-project-issues-statistics
type GetProjectIssuesStatisticsOptions struct {
	IIDs             *[]int64      `url:"iids[],omitempty" json:"iids,omitempty"`
	Labels           *LabelOptions `url:"labels,omitempty" json:"labels,omitempty"`
	Milestone        *string       `url:"milestone,omitempty" json:"milestone,omitempty"`
	Scope            *string       `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID         *int64        `url:"author_id,omitempty" json:"author_id,omitempty"`
	AuthorUsername   *string       `url:"author_username,omitempty" json:"author_username,omitempty"`
	AssigneeID       *int64        `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	AssigneeUsername *[]string     `url:"assignee_username,omitempty" json:"assignee_username,omitempty"`
	MyReactionEmoji  *string       `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	Search           *string       `url:"search,omitempty" json:"search,omitempty"`
	CreatedAfter     *time.Time    `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore    *time.Time    `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter     *time.Time    `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore    *time.Time    `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Confidential     *bool         `url:"confidential,omitempty" json:"confidential,omitempty"`
}

func (s *IssuesStatisticsService) GetProjectIssuesStatistics(pid any, opt *GetProjectIssuesStatisticsOptions, options ...RequestOptionFunc) (*IssuesStatistics, *Response, error) {
	return do[*IssuesStatistics](s.client,
		withPath("projects/%s/issues_statistics", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
