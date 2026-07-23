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
	BranchesServiceInterface interface {
		// ListBranches gets a list of repository branches from a project, sorted by name alphabetically.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/branches/#list-repository-branches
		ListBranches(pid any, opts *ListBranchesOptions, options ...RequestOptionFunc) ([]*Branch, *Response, error)

		// GetBranch gets a single project repository branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/branches/#get-single-repository-branch
		GetBranch(pid any, branch string, options ...RequestOptionFunc) (*Branch, *Response, error)

		// CreateBranch creates branch from commit SHA or existing branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/branches/#create-repository-branch
		CreateBranch(pid any, opt *CreateBranchOptions, options ...RequestOptionFunc) (*Branch, *Response, error)

		// DeleteBranch deletes an existing branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/branches/#delete-repository-branch
		DeleteBranch(pid any, branch string, options ...RequestOptionFunc) (*Response, error)

		// DeleteMergedBranches deletes all branches that are merged into the project's default branch.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/branches/#delete-merged-branches
		DeleteMergedBranches(pid any, options ...RequestOptionFunc) (*Response, error)
	}

	// BranchesService handles communication with the branch related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/branches/
	BranchesService struct {
		client *Client
	}
)

var _ BranchesServiceInterface = (*BranchesService)(nil)

// Branch represents a GitLab branch.
//
// GitLab API docs: https://docs.gitlab.com/api/branches/
type Branch struct {
	Commit             *Commit `json:"commit"`
	Name               string  `json:"name"`
	Protected          bool    `json:"protected"`
	Merged             bool    `json:"merged"`
	Default            bool    `json:"default"`
	CanPush            bool    `json:"can_push"`
	DevelopersCanPush  bool    `json:"developers_can_push"`
	DevelopersCanMerge bool    `json:"developers_can_merge"`
	WebURL             string  `json:"web_url"`
}

func (b Branch) String() string {
	return Stringify(b)
}

// ListBranchesOptions represents the available ListBranches() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/branches/#list-repository-branches
type ListBranchesOptions struct {
	ListOptions
	Search *string `url:"search,omitempty" json:"search,omitempty"`
	Regex  *string `url:"regex,omitempty" json:"regex,omitempty"`
}

func (s *BranchesService) ListBranches(pid any, opts *ListBranchesOptions, options ...RequestOptionFunc) ([]*Branch, *Response, error) {
	return do[[]*Branch](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/branches", ProjectID{pid}),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *BranchesService) GetBranch(pid any, branch string, options ...RequestOptionFunc) (*Branch, *Response, error) {
	return do[*Branch](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/repository/branches/%s", ProjectID{pid}, branch),
		withRequestOpts(options...),
	)
}

// CreateBranchOptions represents the available CreateBranch() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/branches/#create-repository-branch
type CreateBranchOptions struct {
	Branch *string `url:"branch,omitempty" json:"branch,omitempty"`
	Ref    *string `url:"ref,omitempty" json:"ref,omitempty"`
}

func (s *BranchesService) CreateBranch(pid any, opt *CreateBranchOptions, options ...RequestOptionFunc) (*Branch, *Response, error) {
	return do[*Branch](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/branches", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *BranchesService) DeleteBranch(pid any, branch string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/repository/branches/%s", ProjectID{pid}, branch),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *BranchesService) DeleteMergedBranches(pid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/repository/merged_branches", ProjectID{pid}),
		withRequestOpts(options...),
	)
	return resp, err
}
