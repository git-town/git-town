//
// Copyright 2022, Ryan Glab <ryan.j.glab@gmail.com>
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
	// ErrorTrackingServiceInterface defines all the API methods for the ErrorTrackingService
	ErrorTrackingServiceInterface interface {
		// GetErrorTrackingSettings gets error tracking settings.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/error_tracking/#get-error-tracking-settings
		GetErrorTrackingSettings(pid any, options ...RequestOptionFunc) (*ErrorTrackingSettings, *Response, error)

		// EnableDisableErrorTracking allows you to enable or disable the error tracking
		// settings for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/error_tracking/#enable-or-disable-the-error-tracking-project-settings
		EnableDisableErrorTracking(pid any, opt *EnableDisableErrorTrackingOptions, options ...RequestOptionFunc) (*ErrorTrackingSettings, *Response, error)

		// ListClientKeys lists error tracking project client keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/error_tracking/#list-project-client-keys
		ListClientKeys(pid any, opt *ListClientKeysOptions, options ...RequestOptionFunc) ([]*ErrorTrackingClientKey, *Response, error)

		// CreateClientKey creates a new client key for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/error_tracking/#create-a-client-key
		CreateClientKey(pid any, options ...RequestOptionFunc) (*ErrorTrackingClientKey, *Response, error)

		// DeleteClientKey removes a client key from the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/error_tracking/#delete-a-client-key
		DeleteClientKey(pid any, keyID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ErrorTrackingService handles communication with the error tracking
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/error_tracking/
	ErrorTrackingService struct {
		client *Client
	}
)

var _ ErrorTrackingServiceInterface = (*ErrorTrackingService)(nil)

// ErrorTrackingClientKey represents an error tracking client key.
//
// GitLab docs:
// https://docs.gitlab.com/api/error_tracking/#error-tracking-client-keys
type ErrorTrackingClientKey struct {
	ID        int64  `json:"id"`
	Active    bool   `json:"active"`
	PublicKey string `json:"public_key"`
	SentryDsn string `json:"sentry_dsn"`
}

func (p ErrorTrackingClientKey) String() string {
	return Stringify(p)
}

// ErrorTrackingSettings represents error tracking settings for a GitLab project.
//
// GitLab API docs: https://docs.gitlab.com/api/error_tracking/#error-tracking-project-settings
type ErrorTrackingSettings struct {
	Active            bool   `json:"active"`
	ProjectName       string `json:"project_name"`
	SentryExternalURL string `json:"sentry_external_url"`
	APIURL            string `json:"api_url"`
	Integrated        bool   `json:"integrated"`
}

func (p ErrorTrackingSettings) String() string {
	return Stringify(p)
}

func (s *ErrorTrackingService) GetErrorTrackingSettings(pid any, options ...RequestOptionFunc) (*ErrorTrackingSettings, *Response, error) {
	return do[*ErrorTrackingSettings](s.client,
		withPath("projects/%s/error_tracking/settings", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

// EnableDisableErrorTrackingOptions represents the available
// EnableDisableErrorTracking() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/error_tracking/#enable-or-disable-the-error-tracking-project-settings
type EnableDisableErrorTrackingOptions struct {
	Active     *bool `url:"active,omitempty" json:"active,omitempty"`
	Integrated *bool `url:"integrated,omitempty" json:"integrated,omitempty"`
}

func (s *ErrorTrackingService) EnableDisableErrorTracking(pid any, opt *EnableDisableErrorTrackingOptions, options ...RequestOptionFunc) (*ErrorTrackingSettings, *Response, error) {
	return do[*ErrorTrackingSettings](s.client,
		withMethod(http.MethodPatch),
		withPath("projects/%s/error_tracking/settings", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListClientKeysOptions represents the available ListClientKeys() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/error_tracking/#list-project-client-keys
type ListClientKeysOptions struct {
	ListOptions
}

func (s *ErrorTrackingService) ListClientKeys(pid any, opt *ListClientKeysOptions, options ...RequestOptionFunc) ([]*ErrorTrackingClientKey, *Response, error) {
	return do[[]*ErrorTrackingClientKey](s.client,
		withPath("projects/%s/error_tracking/client_keys", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ErrorTrackingService) CreateClientKey(pid any, options ...RequestOptionFunc) (*ErrorTrackingClientKey, *Response, error) {
	return do[*ErrorTrackingClientKey](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/error_tracking/client_keys", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

// DeleteClientKey removes a client key from the project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/error_tracking/#delete-a-client-key
func (s *ErrorTrackingService) DeleteClientKey(pid any, keyID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/error_tracking/client_keys/%d", ProjectID{pid}, keyID),
		withRequestOpts(options...),
	)
	return resp, err
}
