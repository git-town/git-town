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
)

type (
	GroupMarkdownUploadsServiceInterface interface {
		ListGroupMarkdownUploads(gid any, opt *ListMarkdownUploadsOptions, options ...RequestOptionFunc) ([]*GroupMarkdownUpload, *Response, error)
		DownloadGroupMarkdownUploadByID(gid any, uploadID int, options ...RequestOptionFunc) (io.Reader, *Response, error)
		DownloadGroupMarkdownUploadBySecretAndFilename(gid any, secret string, filename string, options ...RequestOptionFunc) (io.Reader, *Response, error)
		DeleteGroupMarkdownUploadByID(gid any, uploadID int, options ...RequestOptionFunc) (*Response, error)
		DeleteGroupMarkdownUploadBySecretAndFilename(gid any, secret string, filename string, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupMarkdownUploadsService handles communication with the group
	// markdown uploads related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_access_tokens/
	GroupMarkdownUploadsService struct {
		client *Client
	}
)

var _ GroupMarkdownUploadsServiceInterface = (*GroupMarkdownUploadsService)(nil)

// Type aliases for backward compatibility
type (
	GroupMarkdownUpload = MarkdownUpload
)

// ListGroupMarkdownUploads gets all markdown uploads for a group.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_markdown_uploads/#list-uploads
func (s *GroupMarkdownUploadsService) ListGroupMarkdownUploads(gid any, opt *ListMarkdownUploadsOptions, options ...RequestOptionFunc) ([]*GroupMarkdownUpload, *Response, error) {
	return listMarkdownUploads[GroupMarkdownUpload](s.client, GroupResource, gid, opt, options)
}

// DownloadGroupMarkdownUploadByID downloads a specific upload by ID.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_markdown_uploads/#download-an-uploaded-file-by-id
func (s *GroupMarkdownUploadsService) DownloadGroupMarkdownUploadByID(gid any, uploadID int, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	buffer, resp, err := downloadMarkdownUploadByID(s.client, GroupResource, gid, uploadID, options)
	if err != nil {
		return nil, resp, err
	}
	return buffer, resp, nil
}

// DownloadGroupMarkdownUploadBySecretAndFilename downloads a specific upload
// by secret and filename.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_markdown_uploads/#download-an-uploaded-file-by-secret-and-filename
func (s *GroupMarkdownUploadsService) DownloadGroupMarkdownUploadBySecretAndFilename(gid any, secret string, filename string, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	buffer, resp, err := downloadMarkdownUploadBySecretAndFilename(s.client, GroupResource, gid, secret, filename, options)
	if err != nil {
		return nil, resp, err
	}
	return buffer, resp, nil
}

// DeleteGroupMarkdownUploadByID deletes an upload by ID.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_markdown_uploads/#delete-an-uploaded-file-by-id
func (s *GroupMarkdownUploadsService) DeleteGroupMarkdownUploadByID(gid any, uploadID int, options ...RequestOptionFunc) (*Response, error) {
	return deleteMarkdownUploadByID(s.client, GroupResource, gid, uploadID, options)
}

// DeleteGroupMarkdownUploadBySecretAndFilename deletes an upload
// by secret and filename.
//
// GitLab API Docs:
// https://docs.gitlab.com/api/group_markdown_uploads/#delete-an-uploaded-file-by-secret-and-filename
func (s *GroupMarkdownUploadsService) DeleteGroupMarkdownUploadBySecretAndFilename(gid any, secret string, filename string, options ...RequestOptionFunc) (*Response, error) {
	return deleteMarkdownUploadBySecretAndFilename(s.client, GroupResource, gid, secret, filename, options)
}
