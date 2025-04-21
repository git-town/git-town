//
// Copyright 2021, Eric Stevens
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
	"time"
)

// GroupHook represents a GitLab group hook.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/
type GroupHook struct {
	ID                        int                 `json:"id"`
	URL                       string              `json:"url"`
	Name                      string              `json:"name"`
	Description               string              `json:"description"`
	CreatedAt                 *time.Time          `json:"created_at"`
	PushEvents                bool                `json:"push_events"`
	TagPushEvents             bool                `json:"tag_push_events"`
	MergeRequestsEvents       bool                `json:"merge_requests_events"`
	RepositoryUpdateEvents    bool                `json:"repository_update_events"`
	EnableSSLVerification     bool                `json:"enable_ssl_verification"`
	AlertStatus               string              `json:"alert_status"`
	PushEventsBranchFilter    string              `json:"push_events_branch_filter"`
	BranchFilterStrategy      string              `json:"branch_filter_strategy"`
	CustomWebhookTemplate     string              `json:"custom_webhook_template"`
	CustomHeaders             []*HookCustomHeader `url:"custom_headers,omitempty" json:"custom_headers,omitempty"`
	GroupID                   int                 `json:"group_id"`
	IssuesEvents              bool                `json:"issues_events"`
	ConfidentialIssuesEvents  bool                `json:"confidential_issues_events"`
	NoteEvents                bool                `json:"note_events"`
	ConfidentialNoteEvents    bool                `json:"confidential_note_events"`
	PipelineEvents            bool                `json:"pipeline_events"`
	WikiPageEvents            bool                `json:"wiki_page_events"`
	JobEvents                 bool                `json:"job_events"`
	DeploymentEvents          bool                `json:"deployment_events"`
	FeatureFlagEvents         bool                `json:"feature_flag_events"`
	ReleasesEvents            bool                `json:"releases_events"`
	SubGroupEvents            bool                `json:"subgroup_events"`
	EmojiEvents               bool                `json:"emoji_events"`
	ResourceAccessTokenEvents bool                `json:"resource_access_token_events"`
	MemberEvents              bool                `json:"member_events"`
}

// ListGroupHooksOptions represents the available ListGroupHooks() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#list-group-hooks
type ListGroupHooksOptions ListOptions

// ListGroupHooks gets a list of group hooks.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#list-group-hooks
func (s *GroupsService) ListGroupHooks(gid interface{}, opt *ListGroupHooksOptions, options ...RequestOptionFunc) ([]*GroupHook, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}
	var gh []*GroupHook
	resp, err := s.client.Do(req, &gh)
	if err != nil {
		return nil, resp, err
	}

	return gh, resp, nil
}

// GetGroupHook gets a specific hook for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#get-a-group-hook
func (s *GroupsService) GetGroupHook(gid interface{}, hook int, options ...RequestOptionFunc) (*GroupHook, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d", PathEscape(group), hook)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	gh := new(GroupHook)
	resp, err := s.client.Do(req, gh)
	if err != nil {
		return nil, resp, err
	}

	return gh, resp, nil
}

// ResendGroupHookEvent resends a specific hook event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#resend-group-hook-event
func (s *GroupsService) ResendGroupHookEvent(gid interface{}, hook int, hookEventID int, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/events/%d/resend", PathEscape(group), hook, hookEventID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddGroupHookOptions represents the available AddGroupHook() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#add-a-group-hook
type AddGroupHookOptions struct {
	URL                       *string              `url:"url,omitempty" json:"url,omitempty"`
	Name                      *string              `url:"name,omitempty" json:"name,omitempty"`
	Description               *string              `url:"description,omitempty" json:"description,omitempty"`
	PushEvents                *bool                `url:"push_events,omitempty"  json:"push_events,omitempty"`
	PushEventsBranchFilter    *string              `url:"push_events_branch_filter,omitempty"  json:"push_events_branch_filter,omitempty"`
	BranchFilterStrategy      *string              `url:"branch_filter_strategy,omitempty"  json:"branch_filter_strategy,omitempty"`
	IssuesEvents              *bool                `url:"issues_events,omitempty"  json:"issues_events,omitempty"`
	ConfidentialIssuesEvents  *bool                `url:"confidential_issues_events,omitempty"  json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents       *bool                `url:"merge_requests_events,omitempty"  json:"merge_requests_events,omitempty"`
	TagPushEvents             *bool                `url:"tag_push_events,omitempty"  json:"tag_push_events,omitempty"`
	NoteEvents                *bool                `url:"note_events,omitempty"  json:"note_events,omitempty"`
	ConfidentialNoteEvents    *bool                `url:"confidential_note_events,omitempty"  json:"confidential_note_events,omitempty"`
	JobEvents                 *bool                `url:"job_events,omitempty"  json:"job_events,omitempty"`
	PipelineEvents            *bool                `url:"pipeline_events,omitempty"  json:"pipeline_events,omitempty"`
	WikiPageEvents            *bool                `url:"wiki_page_events,omitempty"  json:"wiki_page_events,omitempty"`
	DeploymentEvents          *bool                `url:"deployment_events,omitempty" json:"deployment_events,omitempty"`
	FeatureFlagEvents         *bool                `url:"feature_flag_events,omitempty" json:"feature_flag_events,omitempty"`
	ReleasesEvents            *bool                `url:"releases_events,omitempty" json:"releases_events,omitempty"`
	SubGroupEvents            *bool                `url:"subgroup_events,omitempty" json:"subgroup_events,omitempty"`
	MemberEvents              *bool                `url:"member_events,omitempty" json:"member_events,omitempty"`
	EnableSSLVerification     *bool                `url:"enable_ssl_verification,omitempty"  json:"enable_ssl_verification,omitempty"`
	Token                     *string              `url:"token,omitempty" json:"token,omitempty"`
	ResourceAccessTokenEvents *bool                `url:"resource_access_token_events,omitempty" json:"resource_access_token_events,omitempty"`
	CustomWebhookTemplate     *string              `url:"custom_webhook_template,omitempty" json:"custom_webhook_template,omitempty"`
	CustomHeaders             *[]*HookCustomHeader `url:"custom_headers,omitempty" json:"custom_headers,omitempty"`
}

// AddGroupHook creates a new group scoped webhook.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#add-a-group-hook
func (s *GroupsService) AddGroupHook(gid interface{}, opt *AddGroupHookOptions, options ...RequestOptionFunc) (*GroupHook, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gh := new(GroupHook)
	resp, err := s.client.Do(req, gh)
	if err != nil {
		return nil, resp, err
	}

	return gh, resp, nil
}

// EditGroupHookOptions represents the available EditGroupHook() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#edit-group-hook
type EditGroupHookOptions struct {
	URL                                   *string              `url:"url,omitempty" json:"url,omitempty"`
	Name                                  *string              `url:"name,omitempty" json:"name,omitempty"`
	Description                           *string              `url:"description,omitempty" json:"description,omitempty"`
	PushEvents                            *bool                `url:"push_events,omitempty" json:"push_events,omitempty"`
	PushEventsBranchFilter                *string              `url:"push_events_branch_filter,omitempty"  json:"push_events_branch_filter,omitempty"`
	BranchFilterStrategy                  *string              `url:"branch_filter_strategy,omitempty"  json:"branch_filter_strategy,omitempty"`
	IssuesEvents                          *bool                `url:"issues_events,omitempty" json:"issues_events,omitempty"`
	ConfidentialIssuesEvents              *bool                `url:"confidential_issues_events,omitempty" json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents                   *bool                `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	TagPushEvents                         *bool                `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	NoteEvents                            *bool                `url:"note_events,omitempty" json:"note_events,omitempty"`
	ConfidentialNoteEvents                *bool                `url:"confidential_note_events,omitempty" json:"confidential_note_events,omitempty"`
	JobEvents                             *bool                `url:"job_events,omitempty" json:"job_events,omitempty"`
	PipelineEvents                        *bool                `url:"pipeline_events,omitempty" json:"pipeline_events,omitempty"`
	WikiPageEvents                        *bool                `url:"wiki_page_events,omitempty" json:"wiki_page_events,omitempty"`
	DeploymentEvents                      *bool                `url:"deployment_events,omitempty" json:"deployment_events,omitempty"`
	FeatureFlagEvents                     *bool                `url:"feature_flag_events,omitempty" json:"feature_flag_events,omitempty"`
	ReleasesEvents                        *bool                `url:"releases_events,omitempty" json:"releases_events,omitempty"`
	SubGroupEvents                        *bool                `url:"subgroup_events,omitempty" json:"subgroup_events,omitempty"`
	MemberEvents                          *bool                `url:"member_events,omitempty" json:"member_events,omitempty"`
	EnableSSLVerification                 *bool                `url:"enable_ssl_verification,omitempty" json:"enable_ssl_verification,omitempty"`
	ServiceAccessTokensExpirationEnforced *bool                `url:"service_access_tokens_expiration_enforced,omitempty" json:"service_access_tokens_expiration_enforced,omitempty"`
	Token                                 *string              `url:"token,omitempty" json:"token,omitempty"`
	ResourceAccessTokenEvents             *bool                `url:"resource_access_token_events,omitempty" json:"resource_access_token_events,omitempty"`
	CustomWebhookTemplate                 *string              `url:"custom_webhook_template,omitempty" json:"custom_webhook_template,omitempty"`
	CustomHeaders                         *[]*HookCustomHeader `url:"custom_headers,omitempty" json:"custom_headers,omitempty"`
}

// EditGroupHook edits a hook for a specified group.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/group_webhooks/#edit-group-hook
func (s *GroupsService) EditGroupHook(gid interface{}, hook int, opt *EditGroupHookOptions, options ...RequestOptionFunc) (*GroupHook, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d", PathEscape(group), hook)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	gh := new(GroupHook)
	resp, err := s.client.Do(req, gh)
	if err != nil {
		return nil, resp, err
	}

	return gh, resp, nil
}

// DeleteGroupHook removes a hook from a group. This is an idempotent
// method and can be called multiple times.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#delete-a-group-hook
func (s *GroupsService) DeleteGroupHook(gid interface{}, hook int, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d", PathEscape(group), hook)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// TriggerTestGroupHook triggers a test hook for a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#trigger-a-test-group-hook
func (s *GroupsService) TriggerTestGroupHook(pid interface{}, hook int, trigger GroupHookTrigger, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/test/%s", PathEscape(group), hook, trigger)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetGroupCustomHeader creates or updates a group custom webhook header.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#set-a-custom-header
func (s *GroupsService) SetGroupCustomHeader(gid interface{}, hook int, key string, opt *SetHookCustomHeaderOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/custom_headers/%s", PathEscape(group), hook, key)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteGroupCustomHeader deletes a group custom webhook header.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#delete-a-custom-header
func (s *GroupsService) DeleteGroupCustomHeader(gid interface{}, hook int, key string, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/custom_headers/%s", PathEscape(group), hook, key)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetHookURLVariableOptions represents the available SetGroupHookURLVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#set-a-url-variable
type SetHookURLVariableOptions struct {
	Value *string `json:"value,omitempty"`
}

// SetGroupHookURLVariable sets a group hook URL variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#set-a-url-variable
func (s *GroupsService) SetGroupHookURLVariable(gid interface{}, hook int, key string, opt *SetHookURLVariableOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/url_variables/%s", PathEscape(group), hook, key)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteGroupHookURLVariable sets a group hook URL variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_webhooks/#delete-a-url-variable
func (s *GroupsService) DeleteGroupHookURLVariable(gid interface{}, hook int, key string, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/hooks/%d/url_variables/%s", PathEscape(group), hook, key)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
