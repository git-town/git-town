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
	PipelineTriggersServiceInterface interface {
		ListPipelineTriggers(pid any, opt *ListPipelineTriggersOptions, options ...RequestOptionFunc) ([]*PipelineTrigger, *Response, error)
		GetPipelineTrigger(pid any, trigger int64, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error)
		AddPipelineTrigger(pid any, opt *AddPipelineTriggerOptions, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error)
		EditPipelineTrigger(pid any, trigger int64, opt *EditPipelineTriggerOptions, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error)
		DeletePipelineTrigger(pid any, trigger int64, options ...RequestOptionFunc) (*Response, error)
		RunPipelineTrigger(pid any, opt *RunPipelineTriggerOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error)
	}

	// PipelineTriggersService handles Project pipeline triggers.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/pipeline_triggers/
	PipelineTriggersService struct {
		client *Client
	}
)

var _ PipelineTriggersServiceInterface = (*PipelineTriggersService)(nil)

// PipelineTrigger represents a project pipeline trigger.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/
type PipelineTrigger struct {
	ID          int64      `json:"id"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	LastUsed    *time.Time `json:"last_used"`
	Token       string     `json:"token"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Owner       *User      `json:"owner"`
}

// ListPipelineTriggersOptions represents the available ListPipelineTriggers() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#list-project-trigger-tokens
type ListPipelineTriggersOptions struct {
	ListOptions
}

// ListPipelineTriggers gets a list of project triggers.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#list-project-trigger-tokens
func (s *PipelineTriggersService) ListPipelineTriggers(pid any, opt *ListPipelineTriggersOptions, options ...RequestOptionFunc) ([]*PipelineTrigger, *Response, error) {
	return do[[]*PipelineTrigger](s.client,
		withPath("projects/%s/triggers", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetPipelineTrigger gets a specific pipeline trigger for a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#get-trigger-token-details
func (s *PipelineTriggersService) GetPipelineTrigger(pid any, trigger int64, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error) {
	return do[*PipelineTrigger](s.client,
		withPath("projects/%s/triggers/%d", ProjectID{pid}, trigger),
		withRequestOpts(options...),
	)
}

// AddPipelineTriggerOptions represents the available AddPipelineTrigger() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#create-a-trigger-token
type AddPipelineTriggerOptions struct {
	Description *string `url:"description,omitempty" json:"description,omitempty"`
}

// AddPipelineTrigger adds a pipeline trigger to a specified project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#create-a-trigger-token
func (s *PipelineTriggersService) AddPipelineTrigger(pid any, opt *AddPipelineTriggerOptions, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error) {
	return do[*PipelineTrigger](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/triggers", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditPipelineTriggerOptions represents the available EditPipelineTrigger() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#update-a-pipeline-trigger-token
type EditPipelineTriggerOptions struct {
	Description *string `url:"description,omitempty" json:"description,omitempty"`
}

// EditPipelineTrigger edits a trigger for a specified project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#update-a-pipeline-trigger-token
func (s *PipelineTriggersService) EditPipelineTrigger(pid any, trigger int64, opt *EditPipelineTriggerOptions, options ...RequestOptionFunc) (*PipelineTrigger, *Response, error) {
	return do[*PipelineTrigger](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/triggers/%d", ProjectID{pid}, trigger),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeletePipelineTrigger removes a trigger from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#remove-a-pipeline-trigger-token
func (s *PipelineTriggersService) DeletePipelineTrigger(pid any, trigger int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/triggers/%d", ProjectID{pid}, trigger),
		withRequestOpts(options...),
	)
	return resp, err
}

// RunPipelineTriggerOptions represents the available RunPipelineTrigger() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#trigger-a-pipeline-with-a-token
type RunPipelineTriggerOptions struct {
	Ref       *string           `url:"ref" json:"ref"`
	Token     *string           `url:"token" json:"token"`
	Variables map[string]string `url:"variables,omitempty" json:"variables,omitempty"`

	// Inputs contains pipeline input parameters.
	// See PipelineInputsOption for supported types and usage.
	Inputs PipelineInputsOption `url:"inputs,omitempty" json:"inputs,omitempty"`
}

// RunPipelineTrigger starts a trigger from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pipeline_triggers/#trigger-a-pipeline-with-a-token
func (s *PipelineTriggersService) RunPipelineTrigger(pid any, opt *RunPipelineTriggerOptions, options ...RequestOptionFunc) (*Pipeline, *Response, error) {
	return do[*Pipeline](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/trigger/pipeline", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
