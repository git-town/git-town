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
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	WikisServiceInterface interface {
		ListWikis(pid any, opt *ListWikisOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error)
		GetWikiPage(pid any, slug string, opt *GetWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		CreateWikiPage(pid any, opt *CreateWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		EditWikiPage(pid any, slug string, opt *EditWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error)
		DeleteWikiPage(pid any, slug string, options ...RequestOptionFunc) (*Response, error)
		UploadWikiAttachment(pid any, content io.Reader, filename string, opt *UploadWikiAttachmentOptions, options ...RequestOptionFunc) (*WikiAttachment, *Response, error)
	}

	// WikisService handles communication with the wikis related methods of
	// the Gitlab API.
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

// ListWikis lists all pages of the wiki of the given project id.
// When with_content is set, it also returns the content of the pages.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#list-wiki-pages
func (s *WikisService) ListWikis(pid any, opt *ListWikisOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ws []*Wiki
	resp, err := s.client.Do(req, &ws)
	if err != nil {
		return nil, resp, err
	}

	return ws, resp, nil
}

// GetWikiPageOptions represents options to GetWikiPage
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#get-a-wiki-page
type GetWikiPageOptions struct {
	RenderHTML *bool   `url:"render_html,omitempty" json:"render_html,omitempty"`
	Version    *string `url:"version,omitempty" json:"version,omitempty"`
}

// GetWikiPage gets a wiki page for a given project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#get-a-wiki-page
func (s *WikisService) GetWikiPage(pid any, slug string, opt *GetWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis/%s", PathEscape(project), url.PathEscape(slug))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	w := new(Wiki)
	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
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

// CreateWikiPage creates a new wiki page for the given repository with
// the given title, slug, and content.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#create-a-new-wiki-page
func (s *WikisService) CreateWikiPage(pid any, opt *CreateWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	w := new(Wiki)
	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
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

// EditWikiPage Updates an existing wiki page. At least one parameter is
// required to update the wiki page.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#edit-an-existing-wiki-page
func (s *WikisService) EditWikiPage(pid any, slug string, opt *EditWikiPageOptions, options ...RequestOptionFunc) (*Wiki, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis/%s", PathEscape(project), url.PathEscape(slug))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	w := new(Wiki)
	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
}

// DeleteWikiPage deletes a wiki page with a given slug.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#delete-a-wiki-page
func (s *WikisService) DeleteWikiPage(pid any, slug string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis/%s", PathEscape(project), url.PathEscape(slug))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UploadWikiAttachmentOptions represents options to UploadWikiAttachment.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#upload-an-attachment-to-the-wiki-repository
type UploadWikiAttachmentOptions struct {
	Branch *string `url:"branch,omitempty" json:"branch,omitempty"`
}

// UploadWikiAttachment uploads a file to the attachment folder inside the wiki’s repository. The attachment folder is the uploads folder.
//
// GitLab API docs:
// https://docs.gitlab.com/api/wikis/#upload-an-attachment-to-the-wiki-repository
func (s *WikisService) UploadWikiAttachment(pid any, content io.Reader, filename string, opt *UploadWikiAttachmentOptions, options ...RequestOptionFunc) (*WikiAttachment, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/wikis/attachments", PathEscape(project))

	req, err := s.client.UploadRequest(http.MethodPost, u, content, filename, UploadFile, opt, options)
	if err != nil {
		return nil, nil, err
	}

	attachment := new(WikiAttachment)
	resp, err := s.client.Do(req, attachment)
	if err != nil {
		return nil, resp, err
	}

	return attachment, resp, nil
}
