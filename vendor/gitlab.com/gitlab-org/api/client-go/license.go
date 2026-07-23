//
// Copyright 2021, Patrick Webster
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
	LicenseServiceInterface interface {
		GetLicense(options ...RequestOptionFunc) (*License, *Response, error)
		AddLicense(opt *AddLicenseOptions, options ...RequestOptionFunc) (*License, *Response, error)
		DeleteLicense(licenseID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// LicenseService handles communication with the license
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/license/
	LicenseService struct {
		client *Client
	}
)

var _ LicenseServiceInterface = (*LicenseService)(nil)

// License represents a GitLab license.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/
type License struct {
	ID               int64           `json:"id"`
	Plan             string          `json:"plan"`
	CreatedAt        *time.Time      `json:"created_at"`
	StartsAt         *ISOTime        `json:"starts_at"`
	ExpiresAt        *ISOTime        `json:"expires_at"`
	HistoricalMax    int64           `json:"historical_max"`
	MaximumUserCount int64           `json:"maximum_user_count"`
	Expired          bool            `json:"expired"`
	Overage          int64           `json:"overage"`
	UserLimit        int64           `json:"user_limit"`
	ActiveUsers      int64           `json:"active_users"`
	Licensee         LicenseLicensee `json:"licensee"`
	// Add on codes that may occur in legacy licenses that don't have a plan yet.
	// https://gitlab.com/gitlab-org/gitlab/-/blob/master/ee/app/models/license.rb
	AddOns LicenseAddOns `json:"add_ons"`
}

func (l License) String() string {
	return Stringify(l)
}

// LicenseLicensee represents a GitLab license licensee.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/
type LicenseLicensee struct {
	Name    string `json:"Name"`
	Company string `json:"Company"`
	Email   string `json:"Email"`
}

func (l LicenseLicensee) String() string {
	return Stringify(l)
}

// LicenseAddOns represents a GitLab license add ons.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/
type LicenseAddOns struct {
	GitLabAuditorUser int64 `json:"GitLab_Auditor_User"`
	GitLabDeployBoard int64 `json:"GitLab_DeployBoard"`
	GitLabFileLocks   int64 `json:"GitLab_FileLocks"`
	GitLabGeo         int64 `json:"GitLab_Geo"`
	GitLabServiceDesk int64 `json:"GitLab_ServiceDesk"`
}

func (a LicenseAddOns) String() string {
	return Stringify(a)
}

// GetLicense retrieves information about the current license.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/#retrieve-information-about-the-current-license
func (s *LicenseService) GetLicense(options ...RequestOptionFunc) (*License, *Response, error) {
	return do[*License](s.client,
		withPath("license"),
		withRequestOpts(options...),
	)
}

// AddLicenseOptions represents the available AddLicense() options.
//
// https://docs.gitlab.com/api/license/#add-a-new-license
type AddLicenseOptions struct {
	License *string `url:"license" json:"license"`
}

// AddLicense adds a new license.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/#add-a-new-license
func (s *LicenseService) AddLicense(opt *AddLicenseOptions, options ...RequestOptionFunc) (*License, *Response, error) {
	return do[*License](s.client,
		withMethod(http.MethodPost),
		withPath("license"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteLicense deletes an existing license.
//
// GitLab API docs:
// https://docs.gitlab.com/api/license/#delete-a-license
func (s *LicenseService) DeleteLicense(licenseID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("license/%d", licenseID),
		withRequestOpts(options...),
	)
	return resp, err
}
