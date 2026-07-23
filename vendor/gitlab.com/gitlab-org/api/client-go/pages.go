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
	PagesServiceInterface interface {
		// UnpublishPages unpublishes pages. The user must have admin privileges.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/pages/#unpublish-pages
		UnpublishPages(gid any, options ...RequestOptionFunc) (*Response, error)
		// GetPages lists Pages settings for a project. The user must have at least
		// maintainer privileges.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/pages/#get-pages-settings-for-a-project
		GetPages(gid any, options ...RequestOptionFunc) (*Pages, *Response, error)
		// UpdatePages updates Pages settings for a project. The user must have
		// administrator privileges.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/pages/#update-pages-settings-for-a-project
		UpdatePages(pid any, opt UpdatePagesOptions, options ...RequestOptionFunc) (*Pages, *Response, error)
	}

	// PagesService handles communication with the pages related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/pages/
	PagesService struct {
		client *Client
	}
)

var _ PagesServiceInterface = (*PagesService)(nil)

// Pages represents the Pages of a project.
//
// GitLab API docs: https://docs.gitlab.com/api/pages/
type Pages struct {
	URL                   string             `json:"url"`
	IsUniqueDomainEnabled bool               `json:"is_unique_domain_enabled"`
	ForceHTTPS            bool               `json:"force_https"`
	Deployments           []*PagesDeployment `json:"deployments"`
	PrimaryDomain         string             `json:"primary_domain"`
}

// PagesDeployment represents a Pages deployment.
//
// GitLab API docs: https://docs.gitlab.com/api/pages/
type PagesDeployment struct {
	CreatedAt     time.Time `json:"created_at"`
	URL           string    `json:"url"`
	PathPrefix    string    `json:"path_prefix"`
	RootDirectory string    `json:"root_directory"`
}

func (s *PagesService) UnpublishPages(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/pages", ProjectID{gid}),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *PagesService) GetPages(gid any, options ...RequestOptionFunc) (*Pages, *Response, error) {
	return do[*Pages](s.client,
		withPath("projects/%s/pages", ProjectID{gid}),
		withRequestOpts(options...),
	)
}

// UpdatePagesOptions represents the available UpdatePages() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages/#update-pages-settings-for-a-project
type UpdatePagesOptions struct {
	PagesUniqueDomainEnabled *bool   `url:"pages_unique_domain_enabled,omitempty" json:"pages_unique_domain_enabled,omitempty"`
	PagesHTTPSOnly           *bool   `url:"pages_https_only,omitempty" json:"pages_https_only,omitempty"`
	PagesPrimaryDomain       *string `url:"pages_primary_domain,omitempty" json:"pages_primary_domain,omitempty"`
}

func (s *PagesService) UpdatePages(pid any, opt UpdatePagesOptions, options ...RequestOptionFunc) (*Pages, *Response, error) {
	return do[*Pages](s.client,
		withMethod(http.MethodPatch),
		withPath("projects/%s/pages", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
