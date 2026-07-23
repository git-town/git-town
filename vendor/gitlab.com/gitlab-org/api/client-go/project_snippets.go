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
	"bytes"
	"net/http"
)

type (
	ProjectSnippetsServiceInterface interface {
		ListSnippets(pid any, opt *ListProjectSnippetsOptions, options ...RequestOptionFunc) ([]*Snippet, *Response, error)
		GetSnippet(pid any, snippet int64, options ...RequestOptionFunc) (*Snippet, *Response, error)
		CreateSnippet(pid any, opt *CreateProjectSnippetOptions, options ...RequestOptionFunc) (*Snippet, *Response, error)
		UpdateSnippet(pid any, snippet int64, opt *UpdateProjectSnippetOptions, options ...RequestOptionFunc) (*Snippet, *Response, error)
		DeleteSnippet(pid any, snippet int64, options ...RequestOptionFunc) (*Response, error)
		SnippetContent(pid any, snippet int64, options ...RequestOptionFunc) ([]byte, *Response, error)
	}

	// ProjectSnippetsService handles communication with the project snippets
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/project_snippets/
	ProjectSnippetsService struct {
		client *Client
	}
)

var _ ProjectSnippetsServiceInterface = (*ProjectSnippetsService)(nil)

// ListProjectSnippetsOptions represents the available ListSnippets() options.
//
// GitLab API docs: https://docs.gitlab.com/api/project_snippets/#list-snippets
type ListProjectSnippetsOptions struct {
	ListOptions
}

// ListSnippets gets a list of project snippets.
//
// GitLab API docs: https://docs.gitlab.com/api/project_snippets/#list-snippets
func (s *ProjectSnippetsService) ListSnippets(pid any, opt *ListProjectSnippetsOptions, options ...RequestOptionFunc) ([]*Snippet, *Response, error) {
	return do[[]*Snippet](s.client,
		withPath("projects/%s/snippets", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetSnippet gets a single project snippet
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#single-snippet
func (s *ProjectSnippetsService) GetSnippet(pid any, snippet int64, options ...RequestOptionFunc) (*Snippet, *Response, error) {
	return do[*Snippet](s.client,
		withPath("projects/%s/snippets/%d", ProjectID{pid}, snippet),
		withRequestOpts(options...),
	)
}

// CreateProjectSnippetOptions represents the available CreateSnippet() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#create-new-snippet
type CreateProjectSnippetOptions struct {
	Title       *string                      `url:"title,omitempty" json:"title,omitempty"`
	Description *string                      `url:"description,omitempty" json:"description,omitempty"`
	Visibility  *VisibilityValue             `url:"visibility,omitempty" json:"visibility,omitempty"`
	Files       *[]*CreateSnippetFileOptions `url:"files,omitempty" json:"files,omitempty"`

	// Deprecated: use Files instead
	FileName *string `url:"file_name,omitempty" json:"file_name,omitempty"`
	// Deprecated: use Files instead
	Content *string `url:"content,omitempty" json:"content,omitempty"`
}

// CreateSnippet creates a new project snippet. The user must have permission
// to create new snippets.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#create-new-snippet
func (s *ProjectSnippetsService) CreateSnippet(pid any, opt *CreateProjectSnippetOptions, options ...RequestOptionFunc) (*Snippet, *Response, error) {
	return do[*Snippet](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/snippets", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateProjectSnippetOptions represents the available UpdateSnippet() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#update-snippet
type UpdateProjectSnippetOptions struct {
	Title       *string                      `url:"title,omitempty" json:"title,omitempty"`
	Description *string                      `url:"description,omitempty" json:"description,omitempty"`
	Visibility  *VisibilityValue             `url:"visibility,omitempty" json:"visibility,omitempty"`
	Files       *[]*UpdateSnippetFileOptions `url:"files,omitempty" json:"files,omitempty"`

	// Deprecated: use Files instead
	FileName *string `url:"file_name,omitempty" json:"file_name,omitempty"`
	// Deprecated: use Files instead
	Content *string `url:"content,omitempty" json:"content,omitempty"`
}

// UpdateSnippet updates an existing project snippet. The user must have
// permission to change an existing snippet.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#update-snippet
func (s *ProjectSnippetsService) UpdateSnippet(pid any, snippet int64, opt *UpdateProjectSnippetOptions, options ...RequestOptionFunc) (*Snippet, *Response, error) {
	return do[*Snippet](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/snippets/%d", ProjectID{pid}, snippet),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteSnippet deletes an existing project snippet. This is an idempotent
// function and deleting a non-existent snippet still returns a 200 OK status
// code.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#delete-snippet
func (s *ProjectSnippetsService) DeleteSnippet(pid any, snippet int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/snippets/%d", ProjectID{pid}, snippet),
		withRequestOpts(options...),
	)
	return resp, err
}

// SnippetContent returns the raw project snippet as plain text.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_snippets/#snippet-content
func (s *ProjectSnippetsService) SnippetContent(pid any, snippet int64, options ...RequestOptionFunc) ([]byte, *Response, error) {
	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/snippets/%d/raw", ProjectID{pid}, snippet),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return buf.Bytes(), resp, nil
}
