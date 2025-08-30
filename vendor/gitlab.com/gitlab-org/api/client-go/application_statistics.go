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
	ApplicationStatisticsServiceInterface interface {
		GetApplicationStatistics(options ...RequestOptionFunc) (*ApplicationStatistics, *Response, error)
	}

	// ApplicationStatisticsService handles communication with the application
	// statistics related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/statistics/
	ApplicationStatisticsService struct {
		client *Client
	}
)

var _ ApplicationStatisticsServiceInterface = (*ApplicationStatisticsService)(nil)

// ApplicationStatistics represents application statistics.
//
// GitLab API docs: https://docs.gitlab.com/api/statistics/
type ApplicationStatistics struct {
	Forks         int `url:"forks" json:"forks"`
	Issues        int `url:"issues" json:"issues"`
	MergeRequests int `url:"merge_requests" json:"merge_requests"`
	Notes         int `url:"notes" json:"notes"`
	Snippets      int `url:"snippets" json:"snippets"`
	SSHKeys       int `url:"ssh_keys" json:"ssh_keys"`
	Milestones    int `url:"milestones" json:"milestones"`
	Users         int `url:"users" json:"users"`
	Groups        int `url:"groups" json:"groups"`
	Projects      int `url:"projects" json:"projects"`
	ActiveUsers   int `url:"active_users" json:"active_users"`
}

// GetApplicationStatistics gets details on the current application statistics.
//
// GitLab API docs:
// https://docs.gitlab.com/api/statistics/#get-details-on-current-application-statistics
func (s *ApplicationStatisticsService) GetApplicationStatistics(options ...RequestOptionFunc) (*ApplicationStatistics, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "application/statistics", nil, options)
	if err != nil {
		return nil, nil, err
	}

	statistics := new(ApplicationStatistics)
	resp, err := s.client.Do(req, statistics)
	if err != nil {
		return nil, resp, err
	}
	return statistics, resp, nil
}
