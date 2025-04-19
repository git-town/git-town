//
// Copyright 2024, Sander van Harmelen
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
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	ProjectMarkdownUploadsServiceInterface interface {
		UploadProjectMarkdown(pid interface{}, content io.Reader, filename string, options ...RequestOptionFunc) (*ProjectMarkdownUploadedFile, *Response, error)
		ListProjectMarkdownUploads(pid interface{}, options ...RequestOptionFunc) ([]*ProjectMarkdownUpload, *Response, error)
		DownloadProjectMarkdownUploadByID(pid interface{}, uploadID int, options ...RequestOptionFunc) ([]byte, *Response, error)
		DownloadProjectMarkdownUploadBySecretAndFilename(pid interface{}, secret string, filename string, options ...RequestOptionFunc) ([]byte, *Response, error)
		DeleteProjectMarkdownUploadByID(pid interface{}, uploadID int, options ...RequestOptionFunc) (*Response, error)
		DeleteProjectMarkdownUploadBySecretAndFilename(pid interface{}, secret string, filename string, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectMarkdownUploadsService handles communication with the project markdown uploads
	// related methods of the GitLab API.
	//
	// Gitlab API docs: https://docs.gitlab.com/ee/api/project_markdown_uploads.html
	ProjectMarkdownUploadsService struct {
		client *Client
	}
)

var _ ProjectMarkdownUploadsServiceInterface = (*ProjectMarkdownUploadsService)(nil)

// ProjectMarkdownUploadedFile represents a single project markdown uploaded file.
//
// Gitlab API docs: https://docs.gitlab.com/ee/api/project_markdown_uploads.html
type ProjectMarkdownUploadedFile struct {
	ID       int    `json:"id"`
	Alt      string `json:"alt"`
	URL      string `json:"url"`
	FullPath string `json:"full_path"`
	Markdown string `json:"markdown"`
}

// ProjectMarkdownUpload represents a single project markdown upload.
//
// Gitlab API docs: https://docs.gitlab.com/ee/api/project_markdown_uploads.html
type ProjectMarkdownUpload struct {
	ID         int        `json:"id"`
	Size       int        `json:"size"`
	Filename   string     `json:"filename"`
	CreatedAt  *time.Time `json:"created_at"`
	UploadedBy *User      `json:"uploaded_by"`
}

// Gets a string representation of a ProjectMarkdownUpload.
//
// GitLab API docs: https://docs.gitlab.com/ee/api/project_markdown_uploads.html
func (m ProjectMarkdownUpload) String() string {
	return Stringify(m)
}

// UploadProjectMarkdown uploads a markdown file to a project.
//
// GitLab docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#upload-a-file
func (s *ProjectMarkdownUploadsService) UploadProjectMarkdown(pid interface{}, content io.Reader, filename string, options ...RequestOptionFunc) (*ProjectMarkdownUploadedFile, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads", PathEscape(project))

	req, err := s.client.UploadRequest(
		http.MethodPost,
		u,
		content,
		filename,
		UploadFile,
		nil,
		options,
	)
	if err != nil {
		return nil, nil, err
	}

	f := new(ProjectMarkdownUploadedFile)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

// ListProjectMarkdownUploads gets all markdown uploads for a project.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#list-uploads
func (s *ProjectMarkdownUploadsService) ListProjectMarkdownUploads(pid interface{}, options ...RequestOptionFunc) ([]*ProjectMarkdownUpload, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var uploads []*ProjectMarkdownUpload
	resp, err := s.client.Do(req, &uploads)
	if err != nil {
		return nil, resp, err
	}

	return uploads, resp, err
}

// DownloadProjectMarkdownUploadByID downloads a specific upload by ID.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#download-an-uploaded-file-by-id
func (s *ProjectMarkdownUploadsService) DownloadProjectMarkdownUploadByID(pid interface{}, uploadID int, options ...RequestOptionFunc) ([]byte, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads/%d", PathEscape(project), uploadID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var f bytes.Buffer
	resp, err := s.client.Do(req, &f)
	if err != nil {
		return nil, resp, err
	}

	return f.Bytes(), resp, err
}

// DownloadProjectMarkdownUploadBySecretAndFilename downloads a specific upload
// by secret and filename.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#download-an-uploaded-file-by-secret-and-filename
func (s *ProjectMarkdownUploadsService) DownloadProjectMarkdownUploadBySecretAndFilename(pid interface{}, secret string, filename string, options ...RequestOptionFunc) ([]byte, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads/%s/%s", PathEscape(project), PathEscape(secret), PathEscape(filename))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var f bytes.Buffer
	resp, err := s.client.Do(req, &f)
	if err != nil {
		return nil, resp, err
	}

	return f.Bytes(), resp, err
}

// DeleteProjectMarkdownUploadByID deletes an upload by ID.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#delete-an-uploaded-file-by-id
func (s *ProjectMarkdownUploadsService) DeleteProjectMarkdownUploadByID(pid interface{}, uploadID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads/%d", PathEscape(project), uploadID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteProjectMarkdownUploadBySecretAndFilename deletes an upload
// by secret and filename.
//
// GitLab API Docs:
// https://docs.gitlab.com/ee/api/project_markdown_uploads.html#delete-an-uploaded-file-by-secret-and-filename
func (s *ProjectMarkdownUploadsService) DeleteProjectMarkdownUploadBySecretAndFilename(pid interface{}, secret string, filename string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/uploads/%s/%s",
		PathEscape(project), PathEscape(secret), PathEscape(filename))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
