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
	SystemHooksServiceInterface interface {
		// ListHooks gets a list of system hooks.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/system_hooks/#list-system-hooks
		ListHooks(options ...RequestOptionFunc) ([]*Hook, *Response, error)
		// GetHook gets a single system hook.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/system_hooks/#get-system-hook
		GetHook(hook int64, options ...RequestOptionFunc) (*Hook, *Response, error)
		// AddHook adds a new system hook.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/system_hooks/#add-new-system-hook
		AddHook(opt *AddHookOptions, options ...RequestOptionFunc) (*Hook, *Response, error)
		// TestHook tests a system hook.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/system_hooks/#test-system-hook
		TestHook(hook int64, options ...RequestOptionFunc) (*HookEvent, *Response, error)
		// DeleteHook deletes a system hook. This is an idempotent API function and
		// returns 200 OK even if the hook is not available. If the hook is deleted it
		// is also returned as JSON.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/system_hooks/#delete-system-hook
		DeleteHook(hook int64, options ...RequestOptionFunc) (*Response, error)
	}

	// SystemHooksService handles communication with the system hooks related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/system_hooks/
	SystemHooksService struct {
		client *Client
	}
)

var _ SystemHooksServiceInterface = (*SystemHooksService)(nil)

// Hook represents a GitLab system hook.
//
// GitLab API docs: https://docs.gitlab.com/api/system_hooks/
type Hook struct {
	ID                     int64      `json:"id"`
	URL                    string     `json:"url"`
	CreatedAt              *time.Time `json:"created_at"`
	PushEvents             bool       `json:"push_events"`
	TagPushEvents          bool       `json:"tag_push_events"`
	MergeRequestsEvents    bool       `json:"merge_requests_events"`
	RepositoryUpdateEvents bool       `json:"repository_update_events"`
	EnableSSLVerification  bool       `json:"enable_ssl_verification"`
}

func (h Hook) String() string {
	return Stringify(h)
}

func (s *SystemHooksService) ListHooks(options ...RequestOptionFunc) ([]*Hook, *Response, error) {
	return do[[]*Hook](s.client,
		withPath("hooks"),
		withRequestOpts(options...),
	)
}

func (s *SystemHooksService) GetHook(hook int64, options ...RequestOptionFunc) (*Hook, *Response, error) {
	return do[*Hook](s.client,
		withPath("hooks/%d", hook),
		withRequestOpts(options...),
	)
}

// AddHookOptions represents the available AddHook() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/system_hooks/#add-new-system-hook
type AddHookOptions struct {
	URL                    *string `url:"url,omitempty" json:"url,omitempty"`
	Token                  *string `url:"token,omitempty" json:"token,omitempty"`
	PushEvents             *bool   `url:"push_events,omitempty" json:"push_events,omitempty"`
	TagPushEvents          *bool   `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	MergeRequestsEvents    *bool   `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	RepositoryUpdateEvents *bool   `url:"repository_update_events,omitempty" json:"repository_update_events,omitempty"`
	EnableSSLVerification  *bool   `url:"enable_ssl_verification,omitempty" json:"enable_ssl_verification,omitempty"`
}

func (s *SystemHooksService) AddHook(opt *AddHookOptions, options ...RequestOptionFunc) (*Hook, *Response, error) {
	return do[*Hook](s.client,
		withMethod(http.MethodPost),
		withPath("hooks"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// HookEvent represents an event trigger by a GitLab system hook.
//
// GitLab API docs: https://docs.gitlab.com/api/system_hooks/
type HookEvent struct {
	EventName  string `json:"event_name"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	ProjectID  int64  `json:"project_id"`
	OwnerName  string `json:"owner_name"`
	OwnerEmail string `json:"owner_email"`
}

func (h HookEvent) String() string {
	return Stringify(h)
}

func (s *SystemHooksService) TestHook(hook int64, options ...RequestOptionFunc) (*HookEvent, *Response, error) {
	return do[*HookEvent](s.client,
		withPath("hooks/%d", hook),
		withRequestOpts(options...),
	)
}

func (s *SystemHooksService) DeleteHook(hook int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("hooks/%d", hook),
		withRequestOpts(options...),
	)
	return resp, err
}
