//
// Copyright 2021, Arkbriar
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
	"bytes"
	"net/http"
	"time"
)

type (
	JobsServiceInterface interface {
		// ListProjectJobs gets a list of jobs in a project.
		//
		// The scope of jobs to show, one or array of: created, pending, running,
		// failed, success, canceled, skipped; showing all jobs if none provided
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#list-project-jobs
		ListProjectJobs(pid any, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error)
		// ListPipelineJobs gets a list of jobs for specific pipeline in a
		// project. If the pipeline ID is not found, it will respond with 404.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#list-pipeline-jobs
		ListPipelineJobs(pid any, pipelineID int64, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error)
		// ListPipelineBridges gets a list of bridges for specific pipeline in a
		// project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#list-pipeline-trigger-jobs
		ListPipelineBridges(pid any, pipelineID int64, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Bridge, *Response, error)
		// GetJobTokensJob retrieves the job that generated a job token.
		//
		// GitLab API docs: https://docs.gitlab.com/api/jobs/#get-job-tokens-job
		GetJobTokensJob(opts *GetJobTokensJobOptions, options ...RequestOptionFunc) (*Job, *Response, error)
		// GetJob gets a single job of a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#get-a-single-job
		GetJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error)
		// GetJobArtifacts gets jobs artifacts of a project
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#get-job-artifacts
		GetJobArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		// DownloadArtifactsFile downloads the artifacts file from the given
		// reference name and job provided the job finished successfully.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#download-the-artifacts-archive
		DownloadArtifactsFile(pid any, refName string, opt *DownloadArtifactsFileOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		// DownloadSingleArtifactsFile downloads a file from the artifacts from the
		// given reference name and job provided the job finished successfully.
		// Only a single file is going to be extracted from the archive and streamed
		// to a client.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#download-a-single-artifact-file-by-job-id
		DownloadSingleArtifactsFile(pid any, jobID int64, artifactPath string, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		// DownloadSingleArtifactsFileByTagOrBranch downloads a single file from
		// a job's artifacts in the latest successful pipeline using the reference name.
		// The file is extracted from the archive and streamed to the client.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#download-a-single-artifact-file-from-specific-tag-or-branch
		DownloadSingleArtifactsFileByTagOrBranch(pid any, refName string, artifactPath string, opt *DownloadArtifactsFileOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		// GetTraceFile gets a trace of a specific job of a project
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#get-a-log-file
		GetTraceFile(pid any, jobID int64, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		// CancelJob cancels a single job of a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#cancel-a-job
		CancelJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error)
		// RetryJob retries a single job of a project
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#retry-a-job
		RetryJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error)
		// EraseJob erases a single job of a project, removes a job
		// artifacts and a job trace.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#erase-a-job
		EraseJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error)
		// KeepArtifacts prevents artifacts from being deleted when
		// expiration is set.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#keep-artifacts
		KeepArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error)
		// PlayJob triggers a manual action to start a job.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/jobs/#run-a-job
		PlayJob(pid any, jobID int64, opt *PlayJobOptions, options ...RequestOptionFunc) (*Job, *Response, error)
		// DeleteArtifacts deletes artifacts of a job
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#delete-job-artifacts
		DeleteArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*Response, error)
		// DeleteProjectArtifacts deletes artifacts eligible for deletion in a project
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/job_artifacts/#delete-job-artifacts
		DeleteProjectArtifacts(pid any, options ...RequestOptionFunc) (*Response, error)
	}

	// JobsService handles communication with the ci builds related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/jobs/
	JobsService struct {
		client *Client
	}
)

var _ JobsServiceInterface = (*JobsService)(nil)

// Job represents a ci build.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/
type Job struct {
	Commit            *Commit          `json:"commit"`
	Coverage          float64          `json:"coverage"`
	AllowFailure      bool             `json:"allow_failure"`
	CreatedAt         *time.Time       `json:"created_at"`
	StartedAt         *time.Time       `json:"started_at"`
	FinishedAt        *time.Time       `json:"finished_at"`
	ErasedAt          *time.Time       `json:"erased_at"`
	Duration          float64          `json:"duration"`
	QueuedDuration    float64          `json:"queued_duration"`
	ArtifactsExpireAt *time.Time       `json:"artifacts_expire_at"`
	TagList           []string         `json:"tag_list"`
	ID                int64            `json:"id"`
	Name              string           `json:"name"`
	Pipeline          JobPipeline      `json:"pipeline"`
	Ref               string           `json:"ref"`
	Artifacts         []JobArtifact    `json:"artifacts"`
	ArtifactsFile     JobArtifactsFile `json:"artifacts_file"`
	Runner            JobRunner        `json:"runner"`
	Stage             string           `json:"stage"`
	Status            string           `json:"status"`
	FailureReason     string           `json:"failure_reason"`
	Tag               bool             `json:"tag"`
	WebURL            string           `json:"web_url"`
	Project           *Project         `json:"project"`
	User              *User            `json:"user"`
}

// JobPipeline represents a ci build pipeline.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/
type JobPipeline struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	Status    string `json:"status"`
}

// JobArtifact represents a ci build artifact.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/
type JobArtifact struct {
	FileType   string `json:"file_type"`
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	FileFormat string `json:"file_format"`
}

// JobArtifactsFile represents a ci build artifacts file.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/
type JobArtifactsFile struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// JobRunner represents a ci build runner.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/
type JobRunner struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	IsShared    bool   `json:"is_shared"`
	Name        string `json:"name"`
}

// Bridge represents a pipeline bridge.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/#list-pipeline-trigger-jobs
type Bridge struct {
	Commit             *Commit       `json:"commit"`
	Coverage           float64       `json:"coverage"`
	AllowFailure       bool          `json:"allow_failure"`
	CreatedAt          *time.Time    `json:"created_at"`
	StartedAt          *time.Time    `json:"started_at"`
	FinishedAt         *time.Time    `json:"finished_at"`
	ErasedAt           *time.Time    `json:"erased_at"`
	Duration           float64       `json:"duration"`
	QueuedDuration     float64       `json:"queued_duration"`
	ID                 int64         `json:"id"`
	Name               string        `json:"name"`
	Pipeline           PipelineInfo  `json:"pipeline"`
	Ref                string        `json:"ref"`
	Stage              string        `json:"stage"`
	Status             string        `json:"status"`
	FailureReason      string        `json:"failure_reason"`
	Tag                bool          `json:"tag"`
	WebURL             string        `json:"web_url"`
	User               *User         `json:"user"`
	DownstreamPipeline *PipelineInfo `json:"downstream_pipeline"`
}

// ListJobsOptions represents the available ListProjectJobs() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/jobs/#list-project-jobs
type ListJobsOptions struct {
	ListOptions
	Scope          *[]BuildStateValue `url:"scope[],omitempty" json:"scope,omitempty"`
	IncludeRetried *bool              `url:"include_retried,omitempty" json:"include_retried,omitempty"`
}

func (s *JobsService) ListProjectJobs(pid any, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error) {
	return do[[]*Job](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/jobs", ProjectID{pid}),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *JobsService) ListPipelineJobs(pid any, pipelineID int64, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error) {
	return do[[]*Job](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/pipelines/%d/jobs", ProjectID{pid}, pipelineID),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *JobsService) ListPipelineBridges(pid any, pipelineID int64, opts *ListJobsOptions, options ...RequestOptionFunc) ([]*Bridge, *Response, error) {
	return do[[]*Bridge](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/pipelines/%d/bridges", ProjectID{pid}, pipelineID),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// GetJobTokensJobOptions represents the available GetJobTokensJob() options.
//
// GitLab API docs: https://docs.gitlab.com/api/jobs/#get-job-tokens-job
type GetJobTokensJobOptions struct {
	JobToken *string `url:"job_token,omitempty" json:"job_token,omitempty"`
}

func (s *JobsService) GetJobTokensJob(opts *GetJobTokensJobOptions, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodGet),
		withPath("job"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *JobsService) GetJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/jobs/%d", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

func (s *JobsService) GetJobArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/jobs/%d/artifacts", ProjectID{pid}, jobID),
		withRequestOpts(options...),
	)

	return bytes.NewReader(b.Bytes()), resp, err
}

// DownloadArtifactsFileOptions represents the available DownloadArtifactsFile()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/job_artifacts/#download-the-artifacts-archive
type DownloadArtifactsFileOptions struct {
	Job *string `url:"job" json:"job"`
}

func (s *JobsService) DownloadArtifactsFile(pid any, refName string, opt *DownloadArtifactsFileOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/jobs/artifacts/%s/download", ProjectID{pid}, NoEscape{refName}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)

	return bytes.NewReader(b.Bytes()), resp, err
}

func (s *JobsService) DownloadSingleArtifactsFile(pid any, jobID int64, artifactPath string, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/jobs/%d/artifacts/%s", ProjectID{pid}, jobID, NoEscape{artifactPath}),
		withRequestOpts(options...),
	)

	return bytes.NewReader(b.Bytes()), resp, err
}

func (s *JobsService) DownloadSingleArtifactsFileByTagOrBranch(pid any, refName string, artifactPath string, opt *DownloadArtifactsFileOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/jobs/artifacts/%s/raw/%s", ProjectID{pid}, refName, NoEscape{artifactPath}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)

	return bytes.NewReader(b.Bytes()), resp, err
}

func (s *JobsService) GetTraceFile(pid any, jobID int64, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/jobs/%d/trace", ProjectID{pid}, jobID),
		withRequestOpts(options...),
	)

	return bytes.NewReader(b.Bytes()), resp, err
}

func (s *JobsService) CancelJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/jobs/%d/cancel", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

func (s *JobsService) RetryJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/jobs/%d/retry", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

func (s *JobsService) EraseJob(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/jobs/%d/erase", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

func (s *JobsService) KeepArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/jobs/%d/artifacts/keep", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
}

// PlayJobOptions represents the available PlayJob() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/jobs/#run-a-job
type PlayJobOptions struct {
	JobVariablesAttributes *[]*JobVariableOptions `url:"job_variables_attributes,omitempty" json:"job_variables_attributes,omitempty"`
}

// JobVariableOptions represents a single job variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/jobs/#run-a-job
type JobVariableOptions struct {
	Key          *string            `url:"key,omitempty" json:"key,omitempty"`
	Value        *string            `url:"value,omitempty" json:"value,omitempty"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

func (s *JobsService) PlayJob(pid any, jobID int64, opt *PlayJobOptions, options ...RequestOptionFunc) (*Job, *Response, error) {
	return do[*Job](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/jobs/%d/play", ProjectID{pid}, jobID),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *JobsService) DeleteArtifacts(pid any, jobID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/jobs/%d/artifacts", ProjectID{pid}, jobID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *JobsService) DeleteProjectArtifacts(pid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/artifacts", ProjectID{pid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}
