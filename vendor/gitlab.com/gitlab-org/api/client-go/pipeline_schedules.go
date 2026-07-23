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
	"net/http"
	"time"
)

type (
	PipelineSchedulesServiceInterface interface {
		ListPipelineSchedules(pid any, opt *ListPipelineSchedulesOptions, options ...RequestOptionFunc) ([]*PipelineSchedule, *Response, error)
		GetPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error)
		ListPipelinesTriggeredBySchedule(pid any, schedule int64, opt *ListPipelinesTriggeredByScheduleOptions, options ...RequestOptionFunc) ([]*Pipeline, *Response, error)
		CreatePipelineSchedule(pid any, opt *CreatePipelineScheduleOptions, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error)
		EditPipelineSchedule(pid any, schedule int64, opt *EditPipelineScheduleOptions, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error)
		TakeOwnershipOfPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error)
		DeletePipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*Response, error)
		RunPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*Response, error)
		CreatePipelineScheduleVariable(pid any, schedule int64, opt *CreatePipelineScheduleVariableOptions, options ...RequestOptionFunc) (*PipelineVariable, *Response, error)
		EditPipelineScheduleVariable(pid any, schedule int64, key string, opt *EditPipelineScheduleVariableOptions, options ...RequestOptionFunc) (*PipelineVariable, *Response, error)
		DeletePipelineScheduleVariable(pid any, schedule int64, key string, options ...RequestOptionFunc) (*PipelineVariable, *Response, error)
	}

	// PipelineSchedulesService handles communication with the pipeline
	// schedules related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/pipeline_schedules/
	PipelineSchedulesService struct {
		client *Client
	}
)

var _ PipelineSchedulesServiceInterface = (*PipelineSchedulesService)(nil)

// PipelineSchedule represents a pipeline schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/
type PipelineSchedule struct {
	ID           int64               `json:"id"`
	Description  string              `json:"description"`
	Ref          string              `json:"ref"`
	Cron         string              `json:"cron"`
	CronTimezone string              `json:"cron_timezone"`
	NextRunAt    *time.Time          `json:"next_run_at"`
	Active       bool                `json:"active"`
	CreatedAt    *time.Time          `json:"created_at"`
	UpdatedAt    *time.Time          `json:"updated_at"`
	Owner        *User               `json:"owner"`
	LastPipeline *LastPipeline       `json:"last_pipeline"`
	Variables    []*PipelineVariable `json:"variables"`
	Inputs       []*PipelineInput    `json:"inputs"`
}

// LastPipeline represents the last pipeline ran by schedule
// this will be returned only for individual schedule get operation
type LastPipeline struct {
	ID     int64  `json:"id"`
	SHA    string `json:"sha"`
	Ref    string `json:"ref"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// ListPipelineSchedulesOptions represents the available ListPipelineTriggers() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#get-all-pipeline-schedules
type ListPipelineSchedulesOptions struct {
	ListOptions
	Scope *PipelineScheduleScopeValue `url:"scope,omitempty"`
}

// ListPipelineSchedules gets a list of project triggers.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#get-all-pipeline-schedules
func (s *PipelineSchedulesService) ListPipelineSchedules(pid any, opt *ListPipelineSchedulesOptions, options ...RequestOptionFunc) ([]*PipelineSchedule, *Response, error) {
	return do[[]*PipelineSchedule](s.client,
		withPath("projects/%s/pipeline_schedules", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetPipelineSchedule gets a pipeline schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#get-a-single-pipeline-schedule
func (s *PipelineSchedulesService) GetPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error) {
	return do[*PipelineSchedule](s.client,
		withPath("projects/%s/pipeline_schedules/%d", ProjectID{pid}, schedule),
		withRequestOpts(options...),
	)
}

// ListPipelinesTriggeredByScheduleOptions represents the available
// ListPipelinesTriggeredBySchedule() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#get-all-pipelines-triggered-by-a-pipeline-schedule
type ListPipelinesTriggeredByScheduleOptions struct {
	ListOptions
}

// ListPipelinesTriggeredBySchedule gets all pipelines triggered by a pipeline
// schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#get-all-pipelines-triggered-by-a-pipeline-schedule
func (s *PipelineSchedulesService) ListPipelinesTriggeredBySchedule(pid any, schedule int64, opt *ListPipelinesTriggeredByScheduleOptions, options ...RequestOptionFunc) ([]*Pipeline, *Response, error) {
	return do[[]*Pipeline](s.client,
		withPath("projects/%s/pipeline_schedules/%d/pipelines", ProjectID{pid}, schedule),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreatePipelineScheduleOptions represents the available
// CreatePipelineSchedule() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#create-a-new-pipeline-schedule
type CreatePipelineScheduleOptions struct {
	Description  *string          `url:"description" json:"description"`
	Ref          *string          `url:"ref" json:"ref"`
	Cron         *string          `url:"cron" json:"cron"`
	CronTimezone *string          `url:"cron_timezone,omitempty" json:"cron_timezone,omitempty"`
	Active       *bool            `url:"active,omitempty" json:"active,omitempty"`
	Inputs       []*PipelineInput `url:"inputs,omitempty" json:"inputs,omitempty"`
}

// CreatePipelineSchedule creates a pipeline schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#create-a-new-pipeline-schedule
func (s *PipelineSchedulesService) CreatePipelineSchedule(pid any, opt *CreatePipelineScheduleOptions, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error) {
	return do[*PipelineSchedule](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipeline_schedules", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditPipelineScheduleOptions represents the available
// EditPipelineSchedule() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#edit-a-pipeline-schedule
type EditPipelineScheduleOptions struct {
	Description  *string          `url:"description,omitempty" json:"description,omitempty"`
	Ref          *string          `url:"ref,omitempty" json:"ref,omitempty"`
	Cron         *string          `url:"cron,omitempty" json:"cron,omitempty"`
	CronTimezone *string          `url:"cron_timezone,omitempty" json:"cron_timezone,omitempty"`
	Active       *bool            `url:"active,omitempty" json:"active,omitempty"`
	Inputs       []*PipelineInput `url:"inputs,omitempty" json:"inputs,omitempty"`
}

// EditPipelineSchedule edits a pipeline schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#edit-a-pipeline-schedule
func (s *PipelineSchedulesService) EditPipelineSchedule(pid any, schedule int64, opt *EditPipelineScheduleOptions, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error) {
	return do[*PipelineSchedule](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/pipeline_schedules/%d", ProjectID{pid}, schedule),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// TakeOwnershipOfPipelineSchedule sets the owner of the specified
// pipeline schedule to the user issuing the request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#take-ownership-of-a-pipeline-schedule
func (s *PipelineSchedulesService) TakeOwnershipOfPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*PipelineSchedule, *Response, error) {
	return do[*PipelineSchedule](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipeline_schedules/%d/take_ownership", ProjectID{pid}, schedule),
		withRequestOpts(options...),
	)
}

// DeletePipelineSchedule deletes a pipeline schedule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#delete-a-pipeline-schedule
func (s *PipelineSchedulesService) DeletePipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/pipeline_schedules/%d", ProjectID{pid}, schedule),
		withRequestOpts(options...),
	)
	return resp, err
}

// RunPipelineSchedule triggers a new scheduled pipeline to run immediately.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#run-a-scheduled-pipeline-immediately
func (s *PipelineSchedulesService) RunPipelineSchedule(pid any, schedule int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipeline_schedules/%d/play", ProjectID{pid}, schedule),
		withRequestOpts(options...),
	)
	return resp, err
}

// CreatePipelineScheduleVariableOptions represents the available
// CreatePipelineScheduleVariable() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#create-a-new-pipeline-schedule
type CreatePipelineScheduleVariableOptions struct {
	Key          *string            `url:"key" json:"key"`
	Value        *string            `url:"value" json:"value"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

// CreatePipelineScheduleVariable creates a pipeline schedule variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#create-a-new-pipeline-schedule
func (s *PipelineSchedulesService) CreatePipelineScheduleVariable(pid any, schedule int64, opt *CreatePipelineScheduleVariableOptions, options ...RequestOptionFunc) (*PipelineVariable, *Response, error) {
	return do[*PipelineVariable](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pipeline_schedules/%d/variables", ProjectID{pid}, schedule),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditPipelineScheduleVariableOptions represents the available
// EditPipelineScheduleVariable() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#edit-a-pipeline-schedule-variable
type EditPipelineScheduleVariableOptions struct {
	Value        *string            `url:"value" json:"value"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

// EditPipelineScheduleVariable creates a pipeline schedule variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#edit-a-pipeline-schedule-variable
func (s *PipelineSchedulesService) EditPipelineScheduleVariable(pid any, schedule int64, key string, opt *EditPipelineScheduleVariableOptions, options ...RequestOptionFunc) (*PipelineVariable, *Response, error) {
	return do[*PipelineVariable](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/pipeline_schedules/%d/variables/%s", ProjectID{pid}, schedule, NoEscape{key}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeletePipelineScheduleVariable creates a pipeline schedule variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_schedules/#delete-a-pipeline-schedule-variable
func (s *PipelineSchedulesService) DeletePipelineScheduleVariable(pid any, schedule int64, key string, options ...RequestOptionFunc) (*PipelineVariable, *Response, error) {
	return do[*PipelineVariable](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/pipeline_schedules/%d/variables/%s", ProjectID{pid}, schedule, NoEscape{key}),
		withRequestOpts(options...),
	)
}
