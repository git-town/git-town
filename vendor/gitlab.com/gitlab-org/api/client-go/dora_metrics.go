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

type (
	// DORAMetricsServiceInterface defines all the API methods for the DORAMetricsService
	DORAMetricsServiceInterface interface {
		// GetProjectDORAMetrics gets the DORA metrics for a project.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/dora/metrics/#get-project-level-dora-metrics
		GetProjectDORAMetrics(pid any, opt GetDORAMetricsOptions, options ...RequestOptionFunc) ([]DORAMetric, *Response, error)

		// GetGroupDORAMetrics gets the DORA metrics for a group.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/dora/metrics/#get-group-level-dora-metrics
		GetGroupDORAMetrics(gid any, opt GetDORAMetricsOptions, options ...RequestOptionFunc) ([]DORAMetric, *Response, error)
	}

	// DORAMetricsService handles communication with the DORA metrics related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/dora/metrics/
	DORAMetricsService struct {
		client *Client
	}
)

var _ DORAMetricsServiceInterface = (*DORAMetricsService)(nil)

// DORAMetric represents a single DORA metric data point.
//
// GitLab API docs: https://docs.gitlab.com/api/dora/metrics/
type DORAMetric struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// String gets a string representation of a DORAMetric data point
//
// GitLab API docs: https://docs.gitlab.com/api/dora/metrics/
func (m DORAMetric) String() string {
	return Stringify(m)
}

// GetDORAMetricsOptions represent the request body options for getting
// DORA metrics
//
// GitLab API docs: https://docs.gitlab.com/api/dora/metrics/
type GetDORAMetricsOptions struct {
	Metric           *DORAMetricType     `url:"metric,omitempty" json:"metric,omitempty"`
	EndDate          *ISOTime            `url:"end_date,omitempty" json:"end_date,omitempty"`
	EnvironmentTiers *[]string           `url:"environment_tiers,comma,omitempty" json:"environment_tiers,omitempty"`
	Interval         *DORAMetricInterval `url:"interval,omitempty" json:"interval,omitempty"`
	StartDate        *ISOTime            `url:"start_date,omitempty" json:"start_date,omitempty"`
}

func (s *DORAMetricsService) GetProjectDORAMetrics(pid any, opt GetDORAMetricsOptions, options ...RequestOptionFunc) ([]DORAMetric, *Response, error) {
	return do[[]DORAMetric](s.client,
		withPath("projects/%s/dora/metrics", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DORAMetricsService) GetGroupDORAMetrics(gid any, opt GetDORAMetricsOptions, options ...RequestOptionFunc) ([]DORAMetric, *Response, error) {
	return do[[]DORAMetric](s.client,
		withPath("groups/%s/dora/metrics", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
