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
	"fmt"
	"net/http"
)

// timeStatsService handles communication with the time tracking related
// methods of the GitLab API.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#time-tracking
type timeStatsService struct {
	client *Client
}

// TimeStats represents the time estimates and time spent for an issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#time-tracking
type TimeStats struct {
	HumanTimeEstimate   string `json:"human_time_estimate"`
	HumanTotalTimeSpent string `json:"human_total_time_spent"`
	TimeEstimate        int    `json:"time_estimate"`
	TotalTimeSpent      int    `json:"total_time_spent"`
}

func (t TimeStats) String() string {
	return Stringify(t)
}

// SetTimeEstimateOptions represents the available SetTimeEstimate()
// options.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#set-a-time-estimate-for-an-issue
type SetTimeEstimateOptions struct {
	Duration *string `url:"duration,omitempty" json:"duration,omitempty"`
}

// setTimeEstimate sets the time estimate for a single project issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#set-a-time-estimate-for-an-issue
func (s *timeStatsService) setTimeEstimate(pid any, entity string, issue int, opt *SetTimeEstimateOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/time_estimate", PathEscape(project), entity, issue)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(TimeStats)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// resetTimeEstimate resets the time estimate for a single project issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#reset-the-time-estimate-for-an-issue
func (s *timeStatsService) resetTimeEstimate(pid any, entity string, issue int, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/reset_time_estimate", PathEscape(project), entity, issue)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(TimeStats)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// AddSpentTimeOptions represents the available AddSpentTime() options.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#add-spent-time-for-an-issue
type AddSpentTimeOptions struct {
	Duration *string `url:"duration,omitempty" json:"duration,omitempty"`
	Summary  *string `url:"summary,omitempty" json:"summary,omitempty"`
}

// addSpentTime adds spent time for a single project issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#add-spent-time-for-an-issue
func (s *timeStatsService) addSpentTime(pid any, entity string, issue int, opt *AddSpentTimeOptions, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/add_spent_time", PathEscape(project), entity, issue)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(TimeStats)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// resetSpentTime resets the spent time for a single project issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#reset-spent-time-for-an-issue
func (s *timeStatsService) resetSpentTime(pid any, entity string, issue int, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/reset_spent_time", PathEscape(project), entity, issue)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(TimeStats)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// getTimeSpent gets the spent time for a single project issue.
//
// GitLab docs: https://docs.gitlab.com/api/issues/#get-time-tracking-stats
func (s *timeStatsService) getTimeSpent(pid any, entity string, issue int, options ...RequestOptionFunc) (*TimeStats, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/%s/%d/time_stats", PathEscape(project), entity, issue)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(TimeStats)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}
