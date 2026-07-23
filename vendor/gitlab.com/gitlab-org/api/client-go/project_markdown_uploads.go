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
	"io"
	"net/http"
)

type (
	ProjectMarkdownUploadsServiceInterface interface {
		// UploadProjectMarkdown uploads a markdown file to a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#upload-a-file
		UploadProjectMarkdown(pid any, content io.Reader, filename string, options ...RequestOptionFunc) (*ProjectMarkdownUploadedFile, *Response, error)
		// ListProjectMarkdownUploads gets all markdown uploads for a project.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#list-uploads
		ListProjectMarkdownUploads(pid any, options ...RequestOptionFunc) ([]*ProjectMarkdownUpload, *Response, error)
		// DownloadProjectMarkdownUploadByID downloads a specific upload by ID.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#download-an-uploaded-file-by-id
		DownloadProjectMarkdownUploadByID(pid any, uploadID int64, options ...RequestOptionFunc) ([]byte, *Response, error)
		// DownloadProjectMarkdownUploadBySecretAndFilename downloads a specific upload
		// by secret and filename.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#download-an-uploaded-file-by-secret-and-filename
		DownloadProjectMarkdownUploadBySecretAndFilename(pid any, secret string, filename string, options ...RequestOptionFunc) ([]byte, *Response, error)
		// DeleteProjectMarkdownUploadByID deletes an upload by ID.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#delete-an-uploaded-file-by-id
		DeleteProjectMarkdownUploadByID(pid any, uploadID int64, options ...RequestOptionFunc) (*Response, error)
		// DeleteProjectMarkdownUploadBySecretAndFilename deletes an upload
		// by secret and filename.
		//
		// GitLab API Docs:
		// https://docs.gitlab.com/api/project_markdown_uploads/#delete-an-uploaded-file-by-secret-and-filename
		DeleteProjectMarkdownUploadBySecretAndFilename(pid any, secret string, filename string, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectMarkdownUploadsService handles communication with the project
	// markdown uploads related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_markdown_uploads/
	ProjectMarkdownUploadsService struct {
		client *Client
	}
)

var _ ProjectMarkdownUploadsServiceInterface = (*ProjectMarkdownUploadsService)(nil)

// Type aliases for backward compatibility
type (
	ProjectMarkdownUpload       = MarkdownUpload
	ProjectMarkdownUploadedFile = MarkdownUploadedFile
)

func (s *ProjectMarkdownUploadsService) UploadProjectMarkdown(pid any, content io.Reader, filename string, options ...RequestOptionFunc) (*ProjectMarkdownUploadedFile, *Response, error) {
	return do[*ProjectMarkdownUploadedFile](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/uploads", ProjectID{pid}),
		withUpload(content, filename, UploadFile),
		withRequestOpts(options...),
	)
}

func (s *ProjectMarkdownUploadsService) ListProjectMarkdownUploads(pid any, options ...RequestOptionFunc) ([]*ProjectMarkdownUpload, *Response, error) {
	return listMarkdownUploads[ProjectMarkdownUpload](s.client, ProjectResource, ProjectID{pid}, nil, options)
}

func (s *ProjectMarkdownUploadsService) DownloadProjectMarkdownUploadByID(pid any, uploadID int64, options ...RequestOptionFunc) ([]byte, *Response, error) {
	buffer, resp, err := downloadMarkdownUploadByID(s.client, ProjectResource, ProjectID{pid}, uploadID, options)
	if err != nil {
		return nil, resp, err
	}
	return buffer.Bytes(), resp, nil
}

func (s *ProjectMarkdownUploadsService) DownloadProjectMarkdownUploadBySecretAndFilename(pid any, secret string, filename string, options ...RequestOptionFunc) ([]byte, *Response, error) {
	buffer, resp, err := downloadMarkdownUploadBySecretAndFilename(s.client, ProjectResource, ProjectID{pid}, secret, filename, options)
	if err != nil {
		return nil, resp, err
	}
	return buffer.Bytes(), resp, nil
}

func (s *ProjectMarkdownUploadsService) DeleteProjectMarkdownUploadByID(pid any, uploadID int64, options ...RequestOptionFunc) (*Response, error) {
	return deleteMarkdownUploadByID(s.client, ProjectResource, ProjectID{pid}, uploadID, options)
}

func (s *ProjectMarkdownUploadsService) DeleteProjectMarkdownUploadBySecretAndFilename(pid any, secret string, filename string, options ...RequestOptionFunc) (*Response, error) {
	return deleteMarkdownUploadBySecretAndFilename(s.client, ProjectResource, ProjectID{pid}, secret, filename, options)
}
