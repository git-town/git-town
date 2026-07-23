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
	PagesDomainsServiceInterface interface {
		ListPagesDomains(pid any, opt *ListPagesDomainsOptions, options ...RequestOptionFunc) ([]*PagesDomain, *Response, error)
		ListAllPagesDomains(options ...RequestOptionFunc) ([]*PagesDomain, *Response, error)
		GetPagesDomain(pid any, domain string, options ...RequestOptionFunc) (*PagesDomain, *Response, error)
		CreatePagesDomain(pid any, opt *CreatePagesDomainOptions, options ...RequestOptionFunc) (*PagesDomain, *Response, error)
		UpdatePagesDomain(pid any, domain string, opt *UpdatePagesDomainOptions, options ...RequestOptionFunc) (*PagesDomain, *Response, error)
		DeletePagesDomain(pid any, domain string, options ...RequestOptionFunc) (*Response, error)
	}

	// PagesDomainsService handles communication with the pages domains
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/pages_domains/
	PagesDomainsService struct {
		client *Client
	}
)

var _ PagesDomainsServiceInterface = (*PagesDomainsService)(nil)

// PagesDomain represents a pages domain.
//
// GitLab API docs: https://docs.gitlab.com/api/pages_domains/
type PagesDomain struct {
	Domain           string                 `json:"domain"`
	AutoSslEnabled   bool                   `json:"auto_ssl_enabled"`
	URL              string                 `json:"url"`
	ProjectID        int64                  `json:"project_id"`
	Verified         bool                   `json:"verified"`
	VerificationCode string                 `json:"verification_code"`
	EnabledUntil     *time.Time             `json:"enabled_until"`
	Certificate      PagesDomainCertificate `json:"certificate"`
}

// PagesDomainCertificate represents a pages domain certificate.
//
// GitLab API docs: https://docs.gitlab.com/api/pages_domains/
type PagesDomainCertificate struct {
	Subject         string     `json:"subject"`
	Expired         bool       `json:"expired"`
	Expiration      *time.Time `json:"expiration"`
	Certificate     string     `json:"certificate"`
	CertificateText string     `json:"certificate_text"`
}

// ListPagesDomainsOptions represents the available ListPagesDomains() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#list-pages-domains
type ListPagesDomainsOptions struct {
	ListOptions
}

// ListPagesDomains gets a list of project pages domains.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#list-pages-domains
func (s *PagesDomainsService) ListPagesDomains(pid any, opt *ListPagesDomainsOptions, options ...RequestOptionFunc) ([]*PagesDomain, *Response, error) {
	return do[[]*PagesDomain](s.client,
		withPath("projects/%s/pages/domains", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListAllPagesDomains gets a list of all pages domains.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#list-all-pages-domains
func (s *PagesDomainsService) ListAllPagesDomains(options ...RequestOptionFunc) ([]*PagesDomain, *Response, error) {
	return do[[]*PagesDomain](s.client,
		withPath("pages/domains"),
		withRequestOpts(options...),
	)
}

// GetPagesDomain get a specific pages domain for a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#single-pages-domain
func (s *PagesDomainsService) GetPagesDomain(pid any, domain string, options ...RequestOptionFunc) (*PagesDomain, *Response, error) {
	return do[*PagesDomain](s.client,
		withPath("projects/%s/pages/domains/%s", ProjectID{pid}, domain),
		withRequestOpts(options...),
	)
}

// CreatePagesDomainOptions represents the available CreatePagesDomain() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#create-new-pages-domain
type CreatePagesDomainOptions struct {
	Domain         *string `url:"domain,omitempty" json:"domain,omitempty"`
	AutoSslEnabled *bool   `url:"auto_ssl_enabled,omitempty" json:"auto_ssl_enabled,omitempty"`
	Certificate    *string `url:"certificate,omitempty" json:"certificate,omitempty"`
	Key            *string `url:"key,omitempty" json:"key,omitempty"`
}

// CreatePagesDomain creates a new project pages domain.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#create-new-pages-domain
func (s *PagesDomainsService) CreatePagesDomain(pid any, opt *CreatePagesDomainOptions, options ...RequestOptionFunc) (*PagesDomain, *Response, error) {
	return do[*PagesDomain](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/pages/domains", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdatePagesDomainOptions represents the available UpdatePagesDomain() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#update-pages-domain
type UpdatePagesDomainOptions struct {
	AutoSslEnabled *bool   `url:"auto_ssl_enabled,omitempty" json:"auto_ssl_enabled,omitempty"`
	Certificate    *string `url:"certificate,omitempty" json:"certificate,omitempty"`
	Key            *string `url:"key,omitempty" json:"key,omitempty"`
}

// UpdatePagesDomain updates an existing project pages domain.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#update-pages-domain
func (s *PagesDomainsService) UpdatePagesDomain(pid any, domain string, opt *UpdatePagesDomainOptions, options ...RequestOptionFunc) (*PagesDomain, *Response, error) {
	return do[*PagesDomain](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/pages/domains/%s", ProjectID{pid}, domain),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeletePagesDomain deletes an existing project pages domain.
//
// GitLab API docs:
// https://docs.gitlab.com/api/pages_domains/#delete-pages-domain
func (s *PagesDomainsService) DeletePagesDomain(pid any, domain string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/pages/domains/%s", ProjectID{pid}, domain),
		withRequestOpts(options...),
	)
	return resp, err
}
