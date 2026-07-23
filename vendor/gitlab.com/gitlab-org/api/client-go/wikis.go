//
// Copyright 2021, Stany MARCEL
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

package gitlab

import (
	"io"
	"net/http"
)

type (
	WikisServiceInterface interface {
		// ListWikis lists all pages of the wiki of the given project id.
		// When with_content is set, it also returns the content of the pages.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#list-all-wiki-pages
		ListWikis(pid any, opt *ListWikisOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error)
		// GetWikiPage gets a wiki page for a given project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#retrieve-a-wiki-page
		GetWikiPage(pid any, slug string, opt *GetWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		// CreateWikiPage creates a new wiki page for the given repository with
		// the given title, slug, and content.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#create-a-wiki-page
		CreateWikiPage(pid any, opt *CreateWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		// EditWikiPage Updates an existing wiki page. At least one parameter is
		// required to update the wiki page.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#update-a-wiki-page
		EditWikiPage(pid any, slug string, opt *EditWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		// DeleteWikiPage deletes a wiki page with a given slug.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#delete-a-wiki-page
		DeleteWikiPage(pid any, slug string, options ...RequestOptionFunc) (*Response, error)
		// UploadWikiAttachment uploads a file to the attachment folder inside the wiki’s repository. The attachment folder is the uploads folder.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/wikis/#upload-an-attachment-to-the-wiki-repository
		UploadWikiAttachment(pid any, content io.Reader, filename string, opt *UploadWikiAttachmentOptions, options ...RequestOptionFunc) (*WikiAttachment, *Response, error)
	}

	// WikisService handles communication with the wikis related methods of
	// the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/wikis/
	WikisService struct {
		client *Client
	}
)

var _ WikisServiceInterface = (*WikisService)(nil)

// Wiki represents a GitLab wiki.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/
type Wiki struct {
	Content  string          `json:"content"`
	Encoding string          `json:"encoding"`
	Format   WikiFormatValue `json:"format"`
	Slug     string          `json:"slug"`
	Title    string          `json:"title"`
}

func (w Wiki) String() string {
	return Stringify(w)
}

// WikiAttachment represents a GitLab wiki attachment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/
type WikiAttachment struct {
	FileName string             `json:"file_name"`
	FilePath string             `json:"file_path"`
	Branch   string             `json:"branch"`
	Link     WikiAttachmentLink `json:"link"`
}

// WikiAttachmentLink represents a GitLab wiki attachment link.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/
type WikiAttachmentLink struct {
	URL      string `json:"url"`
	Markdown string `json:"markdown"`
}

func (wa WikiAttachment) String() string {
	return Stringify(wa)
}

// ListWikisOptions represents the available ListWikis options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#list-wiki-pages
type ListWikisOptions struct {
	WithContent *bool `url:"with_content,omitempty" json:"with_content,omitempty"`
}

func (s *WikisService) ListWikis(pid any, opt *ListWikisOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	return do[[]*Wiki](s.client,
		withPath("projects/%s/wikis", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetWikiPageOptions represents options to GetWikiPage
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#get-a-wiki-page
type GetWikiPageOptions struct {
	RenderHTML *bool   `url:"render_html,omitempty" json:"render_html,omitempty"`
	Version    *string `url:"version,omitempty" json:"version,omitempty"`
}

func (s *WikisService) GetWikiPage(pid any, slug string, opt *GetWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	return do[*Wiki](s.client,
		withPath("projects/%s/wikis/%s", ProjectID{pid}, slug),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateWikiPageOptions represents options to CreateWikiPage.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#create-a-new-wiki-page
type CreateWikiPageOptions struct {
	Content *string          `url:"content,omitempty" json:"content,omitempty"`
	Title   *string          `url:"title,omitempty" json:"title,omitempty"`
	Format  *WikiFormatValue `url:"format,omitempty" json:"format,omitempty"`
}

func (s *WikisService) CreateWikiPage(pid any, opt *CreateWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	return do[*Wiki](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/wikis", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditWikiPageOptions represents options to EditWikiPage.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#edit-an-existing-wiki-page
type EditWikiPageOptions struct {
	Content *string          `url:"content,omitempty" json:"content,omitempty"`
	Title   *string          `url:"title,omitempty" json:"title,omitempty"`
	Format  *WikiFormatValue `url:"format,omitempty" json:"format,omitempty"`
}

func (s *WikisService) EditWikiPage(pid any, slug string, opt *EditWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	return do[*Wiki](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/wikis/%s", ProjectID{pid}, slug),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *WikisService) DeleteWikiPage(pid any, slug string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/wikis/%s", ProjectID{pid}, slug),
		withRequestOpts(options...),
	)
	return resp, err
}

// UploadWikiAttachmentOptions represents options to UploadWikiAttachment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#upload-an-attachment-to-the-wiki-repository
type UploadWikiAttachmentOptions struct {
	Branch *string `url:"branch,omitempty" json:"branch,omitempty"`
}

func (s *WikisService) UploadWikiAttachment(pid any, content io.Reader, filename string, opt *UploadWikiAttachmentOptions, options ...RequestOptionFunc) (*WikiAttachment, *Response, error) {
	return do[*WikiAttachment](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/wikis/attachments", ProjectID{pid}),
		withUpload(content, filename, UploadFile),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
