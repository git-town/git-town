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
		// GetApplicationStatistics gets details on the current application statistics.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/statistics/#get-details-on-current-application-statistics
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
	Forks         int64 `url:"forks" json:"forks"`
	Issues        int64 `url:"issues" json:"issues"`
	MergeRequests int64 `url:"merge_requests" json:"merge_requests"`
	Notes         int64 `url:"notes" json:"notes"`
	Snippets      int64 `url:"snippets" json:"snippets"`
	SSHKeys       int64 `url:"ssh_keys" json:"ssh_keys"`
	Milestones    int64 `url:"milestones" json:"milestones"`
	Users         int64 `url:"users" json:"users"`
	Groups        int64 `url:"groups" json:"groups"`
	Projects      int64 `url:"projects" json:"projects"`
	ActiveUsers   int64 `url:"active_users" json:"active_users"`
}

func (s *ApplicationStatisticsService) GetApplicationStatistics(options ...RequestOptionFunc) (*ApplicationStatistics, *Response, error) {
	return do[*ApplicationStatistics](s.client,
		withMethod(http.MethodGet),
		withPath("application/statistics"),
		withRequestOpts(options...),
	)
}
