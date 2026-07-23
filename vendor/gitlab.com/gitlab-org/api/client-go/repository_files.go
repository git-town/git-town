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
	"strconv"
	"time"
)

type (
	RepositoryFilesServiceInterface interface {
		GetFile(pid any, fileName string, opt *GetFileOptions, options ...RequestOptionFunc) (*File, *Response, error)
		GetFileMetaData(pid any, fileName string, opt *GetFileMetaDataOptions, options ...RequestOptionFunc) (*File, *Response, error)
		GetFileBlame(pid any, file string, opt *GetFileBlameOptions, options ...RequestOptionFunc) ([]*FileBlameRange, *Response, error)
		GetRawFile(pid any, fileName string, opt *GetRawFileOptions, options ...RequestOptionFunc) ([]byte, *Response, error)
		GetRawFileMetaData(pid any, fileName string, opt *GetRawFileOptions, options ...RequestOptionFunc) (*File, *Response, error)
		CreateFile(pid any, fileName string, opt *CreateFileOptions, options ...RequestOptionFunc) (*FileInfo, *Response, error)
		UpdateFile(pid any, fileName string, opt *UpdateFileOptions, options ...RequestOptionFunc) (*FileInfo, *Response, error)
		DeleteFile(pid any, fileName string, opt *DeleteFileOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// RepositoryFilesService handles communication with the repository files
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/repository_files/
	RepositoryFilesService struct {
		client *Client
	}
)

var _ RepositoryFilesServiceInterface = (*RepositoryFilesService)(nil)

// File represents a GitLab repository file.
//
// GitLab API docs: https://docs.gitlab.com/api/repository_files/
type File struct {
	FileName        string `json:"file_name"`
	FilePath        string `json:"file_path"`
	Size            int64  `json:"size"`
	Encoding        string `json:"encoding"`
	Content         string `json:"content"`
	ExecuteFilemode bool   `json:"execute_filemode"`
	Ref             string `json:"ref"`
	BlobID          string `json:"blob_id"`
	CommitID        string `json:"commit_id"`
	SHA256          string `json:"content_sha256"`
	LastCommitID    string `json:"last_commit_id"`
}

func (r File) String() string {
	return Stringify(r)
}

// GetFileOptions represents the available GetFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-from-repository
type GetFileOptions struct {
	Ref *string `url:"ref,omitempty" json:"ref,omitempty"`
}

// GetFile allows you to receive information about a file in repository like
// name, size, content. Note that file content is Base64 encoded.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-from-repository
func (s *RepositoryFilesService) GetFile(pid any, fileName string, opt *GetFileOptions, options ...RequestOptionFunc) (*File, *Response, error) {
	return do[*File](s.client,
		withPath("projects/%s/repository/files/%s", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetFileMetaDataOptions represents the available GetFileMetaData() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-from-repository
type GetFileMetaDataOptions struct {
	Ref *string `url:"ref,omitempty" json:"ref,omitempty"`
}

// GetFileMetaData allows you to receive meta information about a file in
// repository like name, size.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-from-repository
func (s *RepositoryFilesService) GetFileMetaData(pid any, fileName string, opt *GetFileMetaDataOptions, options ...RequestOptionFunc) (*File, *Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodHead),
		withPath("projects/%s/repository/files/%s", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}

	file, err := getMetaDataFileFromHeaders(resp)
	if err != nil {
		return nil, resp, err
	}

	return file, resp, nil
}

// getMetaDataFileFromHeaders extracts file metadata from response
// headers and converts it into a File object.
func getMetaDataFileFromHeaders(resp *Response) (*File, error) {
	file := &File{
		BlobID:          resp.Header.Get("X-Gitlab-Blob-Id"),
		CommitID:        resp.Header.Get("X-Gitlab-Commit-Id"),
		Encoding:        resp.Header.Get("X-Gitlab-Encoding"),
		FileName:        resp.Header.Get("X-Gitlab-File-Name"),
		FilePath:        resp.Header.Get("X-Gitlab-File-Path"),
		ExecuteFilemode: resp.Header.Get("X-Gitlab-Execute-Filemode") == "true",
		Ref:             resp.Header.Get("X-Gitlab-Ref"),
		SHA256:          resp.Header.Get("X-Gitlab-Content-Sha256"),
		LastCommitID:    resp.Header.Get("X-Gitlab-Last-Commit-Id"),
	}

	if sizeString := resp.Header.Get("X-Gitlab-Size"); sizeString != "" {
		size, err := strconv.ParseInt(sizeString, 10, 64)
		if err != nil {
			return nil, err
		}
		file.Size = size
	}

	return file, nil
}

// FileBlameRange represents one item of blame information.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-blame-from-repository
type FileBlameRange struct {
	Commit FileBlameRangeCommit `json:"commit"`
	Lines  []string             `json:"lines"`
}

func (b FileBlameRange) String() string {
	return Stringify(b)
}

// FileBlameRangeCommit represents one item of blame information's commit.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-blame-from-repository
type FileBlameRangeCommit struct {
	ID             string     `json:"id"`
	ParentIDs      []string   `json:"parent_ids"`
	Message        string     `json:"message"`
	AuthoredDate   *time.Time `json:"authored_date"`
	AuthorName     string     `json:"author_name"`
	AuthorEmail    string     `json:"author_email"`
	CommittedDate  *time.Time `json:"committed_date"`
	CommitterName  string     `json:"committer_name"`
	CommitterEmail string     `json:"committer_email"`
}

func (c FileBlameRangeCommit) String() string {
	return Stringify(c)
}

// GetFileBlameOptions represents the available GetFileBlame() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-blame-from-repository
type GetFileBlameOptions struct {
	Ref        *string `url:"ref,omitempty" json:"ref,omitempty"`
	RangeStart *int64  `url:"range[start],omitempty" json:"range[start],omitempty"`
	RangeEnd   *int64  `url:"range[end],omitempty" json:"range[end],omitempty"`
}

// GetFileBlame allows you to receive blame information. Each blame range
// contains lines and corresponding commit info.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-file-blame-from-repository
func (s *RepositoryFilesService) GetFileBlame(pid any, file string, opt *GetFileBlameOptions, options ...RequestOptionFunc) ([]*FileBlameRange, *Response, error) {
	return do[[]*FileBlameRange](s.client,
		withPath("projects/%s/repository/files/%s/blame", ProjectID{pid}, file),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetRawFileOptions represents the available GetRawFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-raw-file-from-repository
type GetRawFileOptions struct {
	Ref *string `url:"ref,omitempty" json:"ref,omitempty"`
	LFS *bool   `url:"lfs,omitempty" json:"lfs,omitempty"`
}

// GetRawFile gets the contents of a raw file from a repository.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-raw-file-from-repository
func (s *RepositoryFilesService) GetRawFile(pid any, fileName string, opt *GetRawFileOptions, options ...RequestOptionFunc) ([]byte, *Response, error) {
	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/repository/files/%s/raw", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return buf.Bytes(), resp, nil
}

// GetRawFileMetaData gets the metadata of a raw file from a repository.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#get-raw-file-from-repository
func (s *RepositoryFilesService) GetRawFileMetaData(pid any, fileName string, opt *GetRawFileOptions, options ...RequestOptionFunc) (*File, *Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodHead),
		withPath("projects/%s/repository/files/%s/raw", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}

	file, err := getMetaDataFileFromHeaders(resp)
	if err != nil {
		return nil, resp, err
	}

	return file, resp, nil
}

// FileInfo represents file details of a GitLab repository file.
//
// GitLab API docs: https://docs.gitlab.com/api/repository_files/
type FileInfo struct {
	FilePath string `json:"file_path"`
	Branch   string `json:"branch"`
}

func (r FileInfo) String() string {
	return Stringify(r)
}

// CreateFileOptions represents the available CreateFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#create-new-file-in-repository
type CreateFileOptions struct {
	Branch          *string `url:"branch,omitempty" json:"branch,omitempty"`
	StartBranch     *string `url:"start_branch,omitempty" json:"start_branch,omitempty"`
	Encoding        *string `url:"encoding,omitempty" json:"encoding,omitempty"`
	AuthorEmail     *string `url:"author_email,omitempty" json:"author_email,omitempty"`
	AuthorName      *string `url:"author_name,omitempty" json:"author_name,omitempty"`
	Content         *string `url:"content,omitempty" json:"content,omitempty"`
	CommitMessage   *string `url:"commit_message,omitempty" json:"commit_message,omitempty"`
	ExecuteFilemode *bool   `url:"execute_filemode,omitempty" json:"execute_filemode,omitempty"`
}

// CreateFile creates a new file in a repository.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#create-new-file-in-repository
func (s *RepositoryFilesService) CreateFile(pid any, fileName string, opt *CreateFileOptions, options ...RequestOptionFunc) (*FileInfo, *Response, error) {
	return do[*FileInfo](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/repository/files/%s", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateFileOptions represents the available UpdateFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#update-existing-file-in-repository
type UpdateFileOptions struct {
	Branch          *string `url:"branch,omitempty" json:"branch,omitempty"`
	StartBranch     *string `url:"start_branch,omitempty" json:"start_branch,omitempty"`
	Encoding        *string `url:"encoding,omitempty" json:"encoding,omitempty"`
	AuthorEmail     *string `url:"author_email,omitempty" json:"author_email,omitempty"`
	AuthorName      *string `url:"author_name,omitempty" json:"author_name,omitempty"`
	Content         *string `url:"content,omitempty" json:"content,omitempty"`
	CommitMessage   *string `url:"commit_message,omitempty" json:"commit_message,omitempty"`
	LastCommitID    *string `url:"last_commit_id,omitempty" json:"last_commit_id,omitempty"`
	ExecuteFilemode *bool   `url:"execute_filemode,omitempty" json:"execute_filemode,omitempty"`
}

// UpdateFile updates an existing file in a repository
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#update-existing-file-in-repository
func (s *RepositoryFilesService) UpdateFile(pid any, fileName string, opt *UpdateFileOptions, options ...RequestOptionFunc) (*FileInfo, *Response, error) {
	return do[*FileInfo](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/repository/files/%s", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteFileOptions represents the available DeleteFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#delete-existing-file-in-repository
type DeleteFileOptions struct {
	Branch        *string `url:"branch,omitempty" json:"branch,omitempty"`
	StartBranch   *string `url:"start_branch,omitempty" json:"start_branch,omitempty"`
	AuthorEmail   *string `url:"author_email,omitempty" json:"author_email,omitempty"`
	AuthorName    *string `url:"author_name,omitempty" json:"author_name,omitempty"`
	CommitMessage *string `url:"commit_message,omitempty" json:"commit_message,omitempty"`
	LastCommitID  *string `url:"last_commit_id,omitempty" json:"last_commit_id,omitempty"`
}

// DeleteFile deletes an existing file in a repository
//
// GitLab API docs:
// https://docs.gitlab.com/api/repository_files/#delete-existing-file-in-repository
func (s *RepositoryFilesService) DeleteFile(pid any, fileName string, opt *DeleteFileOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/repository/files/%s", ProjectID{pid}, fileName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
