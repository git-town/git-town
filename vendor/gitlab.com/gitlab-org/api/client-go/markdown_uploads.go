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
	"net/http"
	"time"
)

// MarkdownUpload represents a single markdown upload.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_markdown_uploads/
// https://docs.gitlab.com/api/group_markdown_uploads/
type MarkdownUpload struct {
	ID         int64      `json:"id"`
	Size       int64      `json:"size"`
	Filename   string     `json:"filename"`
	CreatedAt  *time.Time `json:"created_at"`
	UploadedBy *User      `json:"uploaded_by"`
}

// String gets a string representation of a MarkdownUpload.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_markdown_uploads/
// https://docs.gitlab.com/api/group_markdown_uploads/
func (m MarkdownUpload) String() string {
	return Stringify(m)
}

// MarkdownUploadedFile represents a single markdown uploaded file.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_markdown_uploads/
type MarkdownUploadedFile struct {
	ID       int64  `json:"id"`
	Alt      string `json:"alt"`
	URL      string `json:"url"`
	FullPath string `json:"full_path"`
	Markdown string `json:"markdown"`
}

// ResourceType represents the type of resource (project or group)
type ResourceType string

const (
	ProjectResource ResourceType = "projects"
	GroupResource   ResourceType = "groups"
)

type ListMarkdownUploadsOptions struct {
	ListOptions
}

// listMarkdownUploads gets all markdown uploads for a resource
func listMarkdownUploads[T any](client *Client, resourceType ResourceType, id Pather, opt *ListMarkdownUploadsOptions, options []RequestOptionFunc) ([]*T, *Response, error) {
	return do[[]*T](client,
		withMethod(http.MethodGet),
		withPath("%s/%s/uploads", resourceType, id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// downloadMarkdownUploadByID downloads a specific upload by ID
func downloadMarkdownUploadByID(client *Client, resourceType ResourceType, id Pather, uploadID int64, options []RequestOptionFunc) (*bytes.Buffer, *Response, error) {
	file, resp, err := do[bytes.Buffer](client,
		withMethod(http.MethodGet),
		withPath("%s/%s/uploads/%d", resourceType, id, uploadID),
		withRequestOpts(options...),
	)

	return &file, resp, err
}

// downloadMarkdownUploadBySecretAndFilename downloads a specific upload by secret and filename
func downloadMarkdownUploadBySecretAndFilename(client *Client, resourceType ResourceType, id Pather, secret string, filename string, options []RequestOptionFunc) (*bytes.Buffer, *Response, error) {
	file, resp, err := do[bytes.Buffer](client,
		withMethod(http.MethodGet),
		withPath("%s/%s/uploads/%s/%s", resourceType, id, secret, filename),
		withRequestOpts(options...),
	)

	return &file, resp, err
}

// deleteMarkdownUploadByID deletes an upload by ID
func deleteMarkdownUploadByID(client *Client, resourceType ResourceType, id Pather, uploadID int64, options []RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](client,
		withMethod(http.MethodDelete),
		withPath("%s/%s/uploads/%d", resourceType, id, uploadID),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// deleteMarkdownUploadBySecretAndFilename deletes an upload by secret and filename
func deleteMarkdownUploadBySecretAndFilename(client *Client, resourceType ResourceType, id Pather, secret string, filename string, options []RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](client,
		withMethod(http.MethodDelete),
		withPath("%s/%s/uploads/%s/%s", resourceType, id, secret, filename),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}
