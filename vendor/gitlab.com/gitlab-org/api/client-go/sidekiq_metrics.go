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
	SidekiqServiceInterface interface {
		GetQueueMetrics(options ...RequestOptionFunc) (*QueueMetrics, *Response, error)
		GetProcessMetrics(options ...RequestOptionFunc) (*ProcessMetrics, *Response, error)
		GetJobStats(options ...RequestOptionFunc) (*JobStats, *Response, error)
		GetCompoundMetrics(options ...RequestOptionFunc) (*CompoundMetrics, *Response, error)
	}

	// SidekiqService handles communication with the sidekiq service
	//
	// GitLab API docs: https://docs.gitlab.com/api/sidekiq_metrics/
	SidekiqService struct {
		client *Client
	}
)

var _ SidekiqServiceInterface = (*SidekiqService)(nil)

// QueueMetrics represents the GitLab sidekiq queue metrics.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-queue-metrics
type QueueMetrics struct {
	Queues map[string]QueueMetricsQueue `json:"queues"`
}

// QueueMetricsQueue represents the GitLab sidekiq queue metrics queue.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-queue-metrics
type QueueMetricsQueue struct {
	Backlog int64 `json:"backlog"`
	Latency int64 `json:"latency"`
}

// GetQueueMetrics lists information about all the registered queues,
// their backlog and their latency.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-queue-metrics
func (s *SidekiqService) GetQueueMetrics(options ...RequestOptionFunc) (*QueueMetrics, *Response, error) {
	return do[*QueueMetrics](s.client,
		withPath("/sidekiq/queue_metrics"),
		withRequestOpts(options...),
	)
}

// ProcessMetrics represents the GitLab sidekiq process metrics.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-process-metrics
type ProcessMetrics struct {
	Processes []ProcessMetricsProcess `json:"processes"`
}

// ProcessMetricsProcess represents the GitLab sidekiq process metrics process.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-process-metrics
type ProcessMetricsProcess struct {
	Hostname    string     `json:"hostname"`
	Pid         int64      `json:"pid"`
	Tag         string     `json:"tag"`
	StartedAt   *time.Time `json:"started_at"`
	Queues      []string   `json:"queues"`
	Labels      []string   `json:"labels"`
	Concurrency int64      `json:"concurrency"`
	Busy        int64      `json:"busy"`
}

// GetProcessMetrics lists information about all the Sidekiq workers registered
// to process your queues.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-process-metrics
func (s *SidekiqService) GetProcessMetrics(options ...RequestOptionFunc) (*ProcessMetrics, *Response, error) {
	return do[*ProcessMetrics](s.client,
		withPath("/sidekiq/process_metrics"),
		withRequestOpts(options...),
	)
}

// JobStats represents the GitLab sidekiq job stats.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-job-statistics
type JobStats struct {
	Jobs JobStatsJobs `json:"jobs"`
}

// JobStatsJobs represents the GitLab sidekiq job stats jobs.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-job-statistics
type JobStatsJobs struct {
	Processed int64 `json:"processed"`
	Failed    int64 `json:"failed"`
	Enqueued  int64 `json:"enqueued"`
}

// GetJobStats list information about the jobs that Sidekiq has performed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-the-current-job-statistics
func (s *SidekiqService) GetJobStats(options ...RequestOptionFunc) (*JobStats, *Response, error) {
	return do[*JobStats](s.client,
		withPath("/sidekiq/job_stats"),
		withRequestOpts(options...),
	)
}

// CompoundMetrics represents the GitLab sidekiq compounded stats.
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-a-compound-response-of-all-the-previously-mentioned-metrics
type CompoundMetrics struct {
	QueueMetrics
	ProcessMetrics
	JobStats
}

// GetCompoundMetrics lists all the currently available information about Sidekiq.
// Get a compound response of all the previously mentioned metrics
//
// GitLab API docs:
// https://docs.gitlab.com/api/sidekiq_metrics/#get-a-compound-response-of-all-the-previously-mentioned-metrics
func (s *SidekiqService) GetCompoundMetrics(options ...RequestOptionFunc) (*CompoundMetrics, *Response, error) {
	return do[*CompoundMetrics](s.client,
		withPath("/sidekiq/compound_metrics"),
		withRequestOpts(options...),
	)
}
