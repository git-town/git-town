//
// Copyright 2021, Igor Varavko
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
	"encoding/json"
	"net/http"
	"time"
)

type PipelineSource string

// PipelineSource is the source of a pipeline.
// GitLab API docs: https://docs.gitlab.com/ci/jobs/job_rules/#ci_pipeline_source-predefined-variable
const (
	PipelineSourceAPI                         PipelineSource = "api"
	PipelineSourceChat                        PipelineSource = "chat"
	PipelineSourceExternal                    PipelineSource = "external"
	PipelineSourceExternalPullRequestEvent    PipelineSource = "external_pull_request_event"
	PipelineSourceMergeRequestEvent           PipelineSource = "merge_request_event"
	PipelineSourceOndemandDastScan            PipelineSource = "ondemand_dast_scan"
	PipelineSourceOndemandDastValidation      PipelineSource = "ondemand_dast_validation"
	PipelineSourceParentPipeline              PipelineSource = "parent_pipeline"
	PipelineSourcePipeline                    PipelineSource = "pipeline"
	PipelineSourcePush                        PipelineSource = "push"
	PipelineSourceSchedule                    PipelineSource = "schedule"
	PipelineSourceSecurityOrchestrationPolicy PipelineSource = "security_orchestration_policy"
	PipelineSourceTrigger                     PipelineSource = "trigger"
	PipelineSourceWeb                         PipelineSource = "web"
	PipelineSourceWebIDE                      PipelineSource = "webide"
)

type (
	PipelinesServiceInterface interface {
		ListProjectPipelines(pid any, opt *ListProjectPipelinesOptions, options ...RequestOptionFunc) ([]*PipelineInfo, *Response, error)
		GetPipeline(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error)
		GetPipelineVariables(pid any, pipeline int64, options ...RequestOptionFunc) ([]*PipelineVariable, *Response, error)
		GetPipelineTestReport(pid any, pipeline int64, options ...RequestOptionFunc) (*PipelineTestReport, *Response, error)
		GetPipelineTestReportSummary(pid any, pipeline int64, options ...RequestOptionFunc) (*PipelineTestReportSummary, *Response, error)
		GetLatestPipeline(pid any, opt *GetLatestPipelineOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error)
		CreatePipeline(pid any, opt *CreatePipelineOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error)
		RetryPipelineBuild(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error)
		CancelPipelineBuild(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error)
		DeletePipeline(pid any, pipeline int64, options ...RequestOptionFunc) (*Response, error)
		UpdatePipelineMetadata(pid any, pipeline int64, opt *UpdatePipelineMetadataOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error)
	}

	// PipelinesService handles communication with the repositories related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/pipelines/
	PipelinesService struct {
		client *Client
	}
)

var _ PipelinesServiceInterface = (*PipelinesService)(nil)

// PipelineVariable represents a pipeline variable.
//
// GitLab API docs: https://docs.gitlab.com/api/pipelines/
type PipelineVariable struct {
	Key          string            `json:"key"`
	Value        string            `json:"value"`
	VariableType VariableTypeValue `json:"variable_type"`
}

// PipelineInput represents a pipeline input.
//
// GitLab API docs: https://docs.gitlab.com/api/pipelines/
type PipelineInput struct {
	Name    string `json:"name"`
	Value   any    `json:"value"`
	Destroy *bool  `json:"destroy,omitempty"`
}

// Pipeline represents a GitLab pipeline.
//
// GitLab API docs: https://docs.gitlab.com/api/pipelines/
type Pipeline struct {
	ID             int64           `json:"id"`
	IID            int64           `json:"iid"`
	ProjectID      int64           `json:"project_id"`
	Status         string          `json:"status"`
	Source         PipelineSource  `json:"source"`
	Ref            string          `json:"ref"`
	Name           string          `json:"name"`
	SHA            string          `json:"sha"`
	BeforeSHA      string          `json:"before_sha"`
	Tag            bool            `json:"tag"`
	YamlErrors     string          `json:"yaml_errors"`
	User           *BasicUser      `json:"user"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	CreatedAt      *time.Time      `json:"created_at"`
	StartedAt      *time.Time      `json:"started_at"`
	FinishedAt     *time.Time      `json:"finished_at"`
	CommittedAt    *time.Time      `json:"committed_at"`
	Duration       int64           `json:"duration"`
	QueuedDuration int64           `json:"queued_duration"`
	Coverage       string          `json:"coverage"`
	WebURL         string          `json:"web_url"`
	DetailedStatus *DetailedStatus `json:"detailed_status"`
}

func (p Pipeline) String() string {
	return Stringify(p)
}

// DetailedStatus contains detailed information about the status of a pipeline.
type DetailedStatus struct {
	Icon         string                     `json:"icon"`
	Text         string                     `json:"text"`
	Label        string                     `json:"label"`
	Group        string                     `json:"group"`
	Tooltip      string                     `json:"tooltip"`
	HasDetails   bool                       `json:"has_details"`
	DetailsPath  string                     `json:"details_path"`
	Illustration DetailedStatusIllustration `json:"illustration"`
	Favicon      string                     `json:"favicon"`
}

func (s DetailedStatus) String() string {
	return Stringify(s)
}

// DetailedStatusIllustration contains detailed information about the status illustration of a pipeline.
type DetailedStatusIllustration struct {
	Image string `json:"image"`
}

func (i DetailedStatusIllustration) String() string {
	return Stringify(i)
}

// PipelineTestReport contains a detailed report of a test run.
type PipelineTestReport struct {
	TotalTime    float64               `json:"total_time"`
	TotalCount   int64                 `json:"total_count"`
	SuccessCount int64                 `json:"success_count"`
	FailedCount  int64                 `json:"failed_count"`
	SkippedCount int64                 `json:"skipped_count"`
	ErrorCount   int64                 `json:"error_count"`
	TestSuites   []*PipelineTestSuites `json:"test_suites"`
}

// PipelineTestSuites contains test suites results.
type PipelineTestSuites struct {
	Name         string               `json:"name"`
	TotalTime    float64              `json:"total_time"`
	TotalCount   int64                `json:"total_count"`
	SuccessCount int64                `json:"success_count"`
	FailedCount  int64                `json:"failed_count"`
	SkippedCount int64                `json:"skipped_count"`
	ErrorCount   int64                `json:"error_count"`
	TestCases    []*PipelineTestCases `json:"test_cases"`
}

// PipelineTestCases contains test cases details.
type PipelineTestCases struct {
	Status         string          `json:"status"`
	Name           string          `json:"name"`
	Classname      string          `json:"classname"`
	File           string          `json:"file"`
	ExecutionTime  float64         `json:"execution_time"`
	SystemOutput   any             `json:"system_output"`
	StackTrace     string          `json:"stack_trace"`
	AttachmentURL  string          `json:"attachment_url"`
	RecentFailures *RecentFailures `json:"recent_failures"`
}

// PipelineTestReportSummary contains a summary report of a test run
type PipelineTestReportSummary struct {
	Total      PipelineTotalSummary       `json:"total"`
	TestSuites []PipelineTestSuiteSummary `json:"test_suites"`
}

// PipelineTotalSummary contains a total summary of a test run
type PipelineTotalSummary struct {
	// Documentation examples only show whole numbers, but the test specs for GitLab show decimals, so `float64` is the better attribute here.
	Time       float64 `json:"time"`
	Count      int64   `json:"count"`
	Success    int64   `json:"success"`
	Failed     int64   `json:"failed"`
	Skipped    int64   `json:"skipped"`
	Error      int64   `json:"error"`
	SuiteError *string `json:"suite_error"`
}

// PipelineTestSuiteSummary contains a test suite summary of a test run
type PipelineTestSuiteSummary struct {
	Name         string  `json:"name"`
	TotalTime    float64 `json:"total_time"`
	TotalCount   int64   `json:"total_count"`
	SuccessCount int64   `json:"success_count"`
	FailedCount  int64   `json:"failed_count"`
	SkippedCount int64   `json:"skipped_count"`
	ErrorCount   int64   `json:"error_count"`
	BuildIDs     []int64 `json:"build_ids"`
	SuiteError   *string `json:"suite_error"`
}

// RecentFailures contains failures count for the project's default branch.
type RecentFailures struct {
	Count      int64  `json:"count"`
	BaseBranch string `json:"base_branch"`
}

func (p PipelineTestReport) String() string {
	return Stringify(p)
}

// PipelineInfo shows the basic entities of a pipeline, mostly used as fields
// on other assets, like Commit.
type PipelineInfo struct {
	ID        int64      `json:"id"`
	IID       int64      `json:"iid"`
	ProjectID int64      `json:"project_id"`
	Status    string     `json:"status"`
	Source    string     `json:"source"`
	Ref       string     `json:"ref"`
	SHA       string     `json:"sha"`
	Name      string     `json:"name"`
	WebURL    string     `json:"web_url"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}

func (p PipelineInfo) String() string {
	return Stringify(p)
}

// ListProjectPipelinesOptions represents the available ListProjectPipelines()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#list-project-pipelines
type ListProjectPipelinesOptions struct {
	ListOptions
	Scope         *string          `url:"scope,omitempty" json:"scope,omitempty"`
	Status        *BuildStateValue `url:"status,omitempty" json:"status,omitempty"`
	Source        *string          `url:"source,omitempty" json:"source,omitempty"`
	Ref           *string          `url:"ref,omitempty" json:"ref,omitempty"`
	SHA           *string          `url:"sha,omitempty" json:"sha,omitempty"`
	YamlErrors    *bool            `url:"yaml_errors,omitempty" json:"yaml_errors,omitempty"`
	Name          *string          `url:"name,omitempty" json:"name,omitempty"`
	Username      *string          `url:"username,omitempty" json:"username,omitempty"`
	UpdatedAfter  *time.Time       `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore *time.Time       `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	OrderBy       *string          `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort          *string          `url:"sort,omitempty" json:"sort,omitempty"`
	CreatedAfter  *time.Time       `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore *time.Time       `url:"created_before,omitempty" json:"created_before,omitempty"`
}

// ListProjectPipelines gets a list of project pipelines.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#list-project-pipelines
func (s *PipelinesService) ListProjectPipelines(pid any, opt *ListProjectPipelinesOptions, options ...RequestOptionFunc) ([]*PipelineInfo, *Response, error) {
	return do[[]*PipelineInfo](s.client,
		withPath("projects/%s/pipelines", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetPipeline gets a single project pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-a-single-pipeline
func (s *PipelinesService) GetPipeline(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withPath("projects/%s/pipelines/%d", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// GetPipelineVariables gets the variables of a single project pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-variables-of-a-pipeline
func (s *PipelinesService) GetPipelineVariables(pid any, pipeline int64, options ...RequestOptionFunc) ([]*PipelineVariable, *Response, error) {
	return do[[]*PipelineVariable](s.client,
		withPath("projects/%s/pipelines/%d/variables", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// GetPipelineTestReport gets the test report of a single project pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-a-pipelines-test-report
func (s *PipelinesService) GetPipelineTestReport(pid any, pipeline int64, options ...RequestOptionFunc) (*PipelineTestReport, *Response, error) {
	return do[*PipelineTestReport](s.client,
		withPath("projects/%s/pipelines/%d/test_report", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// GetPipelineTestReportSummary gets the test report summary of a single project pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-a-test-report-summary-for-a-pipeline
func (s *PipelinesService) GetPipelineTestReportSummary(pid any, pipeline int64, options ...RequestOptionFunc) (*PipelineTestReportSummary, *Response, error) {
	return do[*PipelineTestReportSummary](s.client,
		withPath("projects/%s/pipelines/%d/test_report_summary", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// GetLatestPipelineOptions represents the available GetLatestPipeline() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-the-latest-pipeline
type GetLatestPipelineOptions struct {
	Ref *string `url:"ref,omitempty" json:"ref,omitempty"`
}

// GetLatestPipeline gets the latest pipeline for a specific ref in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#get-the-latest-pipeline
func (s *PipelinesService) GetLatestPipeline(pid any, opt *GetLatestPipelineOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withPath("projects/%s/pipelines/latest", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreatePipelineOptions represents the available CreatePipeline() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#create-a-new-pipeline
type CreatePipelineOptions struct {
	Ref       *string                     `url:"ref" json:"ref"`
	Variables *[]*PipelineVariableOptions `url:"variables,omitempty" json:"variables,omitempty"`

	// Inputs contains pipeline input parameters.
	// See PipelineInputsOption for supported types and usage.
	Inputs PipelineInputsOption `url:"inputs,omitempty" json:"inputs,omitempty"`
}

// PipelineVariableOptions represents a pipeline variable option.
//
// GitLab API docs: https://docs.gitlab.com/api/pipelines/#create-a-new-pipeline
type PipelineVariableOptions struct {
	Key          *string            `url:"key,omitempty" json:"key,omitempty"`
	Value        *string            `url:"value,omitempty" json:"value,omitempty"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

// PipelineInputsOption represents pipeline input parameters with type-safe values.
// Each value must be wrapped using NewPipelineInputValue() to ensure compile-time type safety.
//
// Supported value types:
//   - string
//   - integers (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
//   - floats (float32, float64)
//   - bool
//   - []string
//
// Example:
//
//	inputs := PipelineInputsOption{
//	    "environment": NewPipelineInputValue("production"),
//	    "replicas":    NewPipelineInputValue(3),
//	    "debug":       NewPipelineInputValue(false),
//	    "regions":     NewPipelineInputValue([]string{"us-east", "eu-west"}),
//	}
//
// GitLab API docs:
// - https://docs.gitlab.com/api/pipelines/#create-a-new-pipeline
// - https://docs.gitlab.com/api/pipeline_triggers/#trigger-a-pipeline-with-a-token
type PipelineInputsOption map[string]PipelineInputValueInterface

// PipelineInputValueInterface is implemented by PipelineInputValue[T] for supported pipeline input types.
// Use NewPipelineInputValue() to create instances - do not implement this interface directly.
//
// See PipelineInputsOption for supported types and usage examples.
type PipelineInputValueInterface interface {
	pipelineInputValue()
}

type constraintSigned interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type constraintUnsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type constraintInteger interface {
	constraintSigned | constraintUnsigned
}

type constraintFloat interface {
	~float32 | ~float64
}

// PipelineInputValueType is a type constraint for valid pipeline input value types.
// This constraint ensures only supported GitLab pipeline input types can be used.
type PipelineInputValueType interface {
	~string | constraintInteger | constraintFloat | ~bool | []string
}

// PipelineInputValue wraps a pipeline input value with compile-time type safety.
// Use NewPipelineInputValue() to create instances of this type.
type PipelineInputValue[T PipelineInputValueType] struct {
	Value T
}

// MarshalJSON implements the json.Marshaler interface.
func (v PipelineInputValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

// pipelineInputValue implements PipelineInputValueInterface.
func (PipelineInputValue[T]) pipelineInputValue() {}

// NewPipelineInputValue wraps a value for use in pipeline inputs.
// Similar to Ptr(), this ensures type safety at compile time.
// Supported types: string, integers, floats, bool, []string
func NewPipelineInputValue[T PipelineInputValueType](value T) PipelineInputValue[T] {
	return PipelineInputValue[T]{
		Value: value,
	}
}

// CreatePipeline creates a new project pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#create-a-new-pipeline
func (s *PipelinesService) CreatePipeline(pid any, opt *CreatePipelineOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipeline", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RetryPipelineBuild retries failed builds in a pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#retry-jobs-in-a-pipeline
func (s *PipelinesService) RetryPipelineBuild(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipelines/%d/retry", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// CancelPipelineBuild cancels a pipeline builds.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#cancel-a-pipelines-jobs
func (s *PipelinesService) CancelPipelineBuild(pid any, pipeline int64, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipelines/%d/cancel", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
}

// DeletePipeline deletes an existing pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#delete-a-pipeline
func (s *PipelinesService) DeletePipeline(pid any, pipeline int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/pipelines/%d", ProjectID{pid}, pipeline),
		withRequestOpts(options...),
	)
	return resp, err
}

// UpdatePipelineMetadataOptions represents the available UpdatePipelineMetadata()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#update-pipeline-metadata
type UpdatePipelineMetadataOptions struct {
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

// UpdatePipelineMetadata You can update the metadata of a pipeline. The metadata
// contains the name of the pipeline.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipelines/#update-pipeline-metadata
func (s *PipelinesService) UpdatePipelineMetadata(pid any, pipeline int64, opt *UpdatePipelineMetadataOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/pipelines/%d/metadata", ProjectID{pid}, pipeline),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
