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
	"errors"
	"net/http"
)

type (
	NotificationSettingsServiceInterface interface {
		// GetGlobalSettings returns current notification settings and email address.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#global-notification-settings
		GetGlobalSettings(options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
		// UpdateGlobalSettings updates current notification settings and email address.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#update-global-notification-settings
		UpdateGlobalSettings(opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
		// GetSettingsForGroup returns current group notification settings.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#group--project-level-notification-settings
		GetSettingsForGroup(gid any, options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
		// GetSettingsForProject returns current project notification settings.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#group--project-level-notification-settings
		GetSettingsForProject(pid any, options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
		// UpdateSettingsForGroup updates current group notification settings.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#update-groupproject-level-notification-settings
		UpdateSettingsForGroup(gid any, opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
		// UpdateSettingsForProject updates current project notification settings.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/notification_settings/#update-groupproject-level-notification-settings
		UpdateSettingsForProject(pid any, opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error)
	}

	// NotificationSettingsService handles communication with the notification settings
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/notification_settings/
	NotificationSettingsService struct {
		client *Client
	}
)

var _ NotificationSettingsServiceInterface = (*NotificationSettingsService)(nil)

// NotificationSettings represents the GitLab notification setting.
//
// GitLab API docs:
// https://docs.gitlab.com/api/notification_settings/#valid-notification-levels
type NotificationSettings struct {
	Level             NotificationLevelValue `json:"level"`
	NotificationEmail string                 `json:"notification_email"`
	Events            *NotificationEvents    `json:"events"`
}

// NotificationEvents represents the available notification setting events.
//
// GitLab API docs:
// https://docs.gitlab.com/api/notification_settings/#valid-notification-levels
type NotificationEvents struct {
	CloseIssue                bool `json:"close_issue"`
	CloseMergeRequest         bool `json:"close_merge_request"`
	FailedPipeline            bool `json:"failed_pipeline"`
	FixedPipeline             bool `json:"fixed_pipeline"`
	IssueDue                  bool `json:"issue_due"`
	MergeWhenPipelineSucceeds bool `json:"merge_when_pipeline_succeeds"`
	MergeMergeRequest         bool `json:"merge_merge_request"`
	MovedProject              bool `json:"moved_project"`
	NewIssue                  bool `json:"new_issue"`
	NewMergeRequest           bool `json:"new_merge_request"`
	NewEpic                   bool `json:"new_epic"`
	NewNote                   bool `json:"new_note"`
	PushToMergeRequest        bool `json:"push_to_merge_request"`
	ReassignIssue             bool `json:"reassign_issue"`
	ReassignMergeRequest      bool `json:"reassign_merge_request"`
	ReopenIssue               bool `json:"reopen_issue"`
	ReopenMergeRequest        bool `json:"reopen_merge_request"`
	SuccessPipeline           bool `json:"success_pipeline"`
}

func (ns NotificationSettings) String() string {
	return Stringify(ns)
}

func (s *NotificationSettingsService) GetGlobalSettings(options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	return do[*NotificationSettings](s.client,
		withPath("notification_settings"),
		withRequestOpts(options...),
	)
}

// NotificationSettingsOptions represents the available options that can be passed
// to the API when updating the notification settings.
type NotificationSettingsOptions struct {
	Level                     *NotificationLevelValue `url:"level,omitempty" json:"level,omitempty"`
	NotificationEmail         *string                 `url:"notification_email,omitempty" json:"notification_email,omitempty"`
	CloseIssue                *bool                   `url:"close_issue,omitempty" json:"close_issue,omitempty"`
	CloseMergeRequest         *bool                   `url:"close_merge_request,omitempty" json:"close_merge_request,omitempty"`
	FailedPipeline            *bool                   `url:"failed_pipeline,omitempty" json:"failed_pipeline,omitempty"`
	FixedPipeline             *bool                   `url:"fixed_pipeline,omitempty" json:"fixed_pipeline,omitempty"`
	IssueDue                  *bool                   `url:"issue_due,omitempty" json:"issue_due,omitempty"`
	MergeMergeRequest         *bool                   `url:"merge_merge_request,omitempty" json:"merge_merge_request,omitempty"`
	MergeWhenPipelineSucceeds *bool                   `url:"merge_when_pipeline_succeeds,omitempty" json:"merge_when_pipeline_succeeds,omitempty"`
	MovedProject              *bool                   `url:"moved_project,omitempty" json:"moved_project,omitempty"`
	NewEpic                   *bool                   `url:"new_epic,omitempty" json:"new_epic,omitempty"`
	NewIssue                  *bool                   `url:"new_issue,omitempty" json:"new_issue,omitempty"`
	NewMergeRequest           *bool                   `url:"new_merge_request,omitempty" json:"new_merge_request,omitempty"`
	NewNote                   *bool                   `url:"new_note,omitempty" json:"new_note,omitempty"`
	PushToMergeRequest        *bool                   `url:"push_to_merge_request,omitempty" json:"push_to_merge_request,omitempty"`
	ReassignIssue             *bool                   `url:"reassign_issue,omitempty" json:"reassign_issue,omitempty"`
	ReassignMergeRequest      *bool                   `url:"reassign_merge_request,omitempty" json:"reassign_merge_request,omitempty"`
	ReopenIssue               *bool                   `url:"reopen_issue,omitempty" json:"reopen_issue,omitempty"`
	ReopenMergeRequest        *bool                   `url:"reopen_merge_request,omitempty" json:"reopen_merge_request,omitempty"`
	SuccessPipeline           *bool                   `url:"success_pipeline,omitempty" json:"success_pipeline,omitempty"`
}

func (s *NotificationSettingsService) UpdateGlobalSettings(opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	if opt.Level != nil && *opt.Level == GlobalNotificationLevel {
		return nil, nil, errors.New(
			"notification level 'global' is not valid for global notification settings")
	}

	return do[*NotificationSettings](s.client,
		withMethod(http.MethodPut),
		withPath("notification_settings"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *NotificationSettingsService) GetSettingsForGroup(gid any, options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	return do[*NotificationSettings](s.client,
		withPath("groups/%s/notification_settings", GroupID{gid}),
		withRequestOpts(options...),
	)
}

func (s *NotificationSettingsService) GetSettingsForProject(pid any, options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	return do[*NotificationSettings](s.client,
		withPath("projects/%s/notification_settings", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

func (s *NotificationSettingsService) UpdateSettingsForGroup(gid any, opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	return do[*NotificationSettings](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/notification_settings", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *NotificationSettingsService) UpdateSettingsForProject(pid any, opt *NotificationSettingsOptions, options ...RequestOptionFunc) (*NotificationSettings, *Response, error) {
	return do[*NotificationSettings](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/notification_settings", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
