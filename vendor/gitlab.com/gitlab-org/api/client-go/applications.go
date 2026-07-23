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
)

type (
	ApplicationsServiceInterface interface {
		// CreateApplication creates a new application owned by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/applications/#create-an-application
		CreateApplication(opt *CreateApplicationOptions, options ...RequestOptionFunc) (*Application, *Response, error)

		// ListApplications get a list of administrables applications by the authenticated user
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/applications/#list-all-applications
		ListApplications(opt *ListApplicationsOptions, options ...RequestOptionFunc) ([]*Application, *Response, error)

		// DeleteApplication removes a specific application.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/applications/#delete-an-application
		DeleteApplication(application int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ApplicationsService handles communication with administrables applications
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/applications/
	ApplicationsService struct {
		client *Client
	}
)

var _ ApplicationsServiceInterface = (*ApplicationsService)(nil)

// Application represents a GitLab application
type Application struct {
	ID              int64  `json:"id"`
	ApplicationID   string `json:"application_id"`
	ApplicationName string `json:"application_name"`
	Secret          string `json:"secret"`
	CallbackURL     string `json:"callback_url"`
	Confidential    bool   `json:"confidential"`
}

// CreateApplicationOptions represents the available CreateApplication() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/applications/#create-an-application
type CreateApplicationOptions struct {
	Name         *string `url:"name,omitempty" json:"name,omitempty"`
	RedirectURI  *string `url:"redirect_uri,omitempty" json:"redirect_uri,omitempty"`
	Scopes       *string `url:"scopes,omitempty" json:"scopes,omitempty"`
	Confidential *bool   `url:"confidential,omitempty" json:"confidential,omitempty"`
}

func (s *ApplicationsService) CreateApplication(opt *CreateApplicationOptions, options ...RequestOptionFunc) (*Application, *Response, error) {
	return do[*Application](s.client,
		withMethod(http.MethodPost),
		withPath("applications"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListApplicationsOptions represents the available
// ListApplications() options.
type ListApplicationsOptions struct {
	ListOptions
}

func (s *ApplicationsService) ListApplications(opt *ListApplicationsOptions, options ...RequestOptionFunc) ([]*Application, *Response, error) {
	return do[[]*Application](s.client,
		withMethod(http.MethodGet),
		withPath("applications"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ApplicationsService) DeleteApplication(application int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("applications/%d", application),
		withRequestOpts(options...),
	)
	return resp, err
}
