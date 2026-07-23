//
// Copyright 2023, Nick Westbury
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
	SnippetRepositoryStorageMoveServiceInterface interface {
		RetrieveAllStorageMoves(opts RetrieveAllSnippetStorageMovesOptions, options ...RequestOptionFunc) ([]*SnippetRepositoryStorageMove, *Response, error)
		RetrieveAllStorageMovesForSnippet(snippet int64, opts RetrieveAllSnippetStorageMovesOptions, options ...RequestOptionFunc) ([]*SnippetRepositoryStorageMove, *Response, error)
		GetStorageMove(repositoryStorage int64, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error)
		GetStorageMoveForSnippet(snippet int64, repositoryStorage int64, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error)
		ScheduleStorageMoveForSnippet(snippet int64, opts ScheduleStorageMoveForSnippetOptions, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error)
		ScheduleAllStorageMoves(opts ScheduleAllSnippetStorageMovesOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// SnippetRepositoryStorageMoveService handles communication with the
	// snippets related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/snippet_repository_storage_moves/
	SnippetRepositoryStorageMoveService struct {
		client *Client
	}
)

var _ SnippetRepositoryStorageMoveServiceInterface = (*SnippetRepositoryStorageMoveService)(nil)

// SnippetRepositoryStorageMove represents the status of a repository move.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/
type SnippetRepositoryStorageMove struct {
	ID                     int64              `json:"id"`
	CreatedAt              *time.Time         `json:"created_at"`
	State                  string             `json:"state"`
	SourceStorageName      string             `json:"source_storage_name"`
	DestinationStorageName string             `json:"destination_storage_name"`
	Snippet                *RepositorySnippet `json:"snippet"`
}

type RepositorySnippet struct {
	ID            int64           `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Visibility    VisibilityValue `json:"visibility"`
	UpdatedAt     *time.Time      `json:"updated_at"`
	CreatedAt     *time.Time      `json:"created_at"`
	ProjectID     int64           `json:"project_id"`
	WebURL        string          `json:"web_url"`
	RawURL        string          `json:"raw_url"`
	SSHURLToRepo  string          `json:"ssh_url_to_repo"`
	HTTPURLToRepo string          `json:"http_url_to_repo"`
}

// RetrieveAllSnippetStorageMovesOptions represents the available
// RetrieveAllStorageMoves() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#retrieve-all-snippet-repository-storage-moves
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#retrieve-all-repository-storage-moves-for-a-snippet
type RetrieveAllSnippetStorageMovesOptions struct {
	ListOptions
}

// RetrieveAllStorageMoves retrieves all snippet repository storage moves
// accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#retrieve-all-snippet-repository-storage-moves
func (s SnippetRepositoryStorageMoveService) RetrieveAllStorageMoves(opts RetrieveAllSnippetStorageMovesOptions, options ...RequestOptionFunc) ([]*SnippetRepositoryStorageMove, *Response, error) {
	return do[[]*SnippetRepositoryStorageMove](s.client,
		withPath("snippet_repository_storage_moves"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// RetrieveAllStorageMovesForSnippet retrieves all repository storage moves for
// a single snippet accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#retrieve-all-repository-storage-moves-for-a-snippet
func (s SnippetRepositoryStorageMoveService) RetrieveAllStorageMovesForSnippet(snippet int64, opts RetrieveAllSnippetStorageMovesOptions, options ...RequestOptionFunc) ([]*SnippetRepositoryStorageMove, *Response, error) {
	return do[[]*SnippetRepositoryStorageMove](s.client,
		withPath("snippets/%d/repository_storage_moves", snippet),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// GetStorageMove gets a single snippet repository storage move.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#get-a-single-snippet-repository-storage-move
func (s SnippetRepositoryStorageMoveService) GetStorageMove(repositoryStorage int64, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error) {
	return do[*SnippetRepositoryStorageMove](s.client,
		withPath("snippet_repository_storage_moves/%d", repositoryStorage),
		withRequestOpts(options...),
	)
}

// GetStorageMoveForSnippet gets a single repository storage move for a snippet.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#get-a-single-repository-storage-move-for-a-snippet
func (s SnippetRepositoryStorageMoveService) GetStorageMoveForSnippet(snippet int64, repositoryStorage int64, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error) {
	return do[*SnippetRepositoryStorageMove](s.client,
		withPath("snippets/%d/repository_storage_moves/%d", snippet, repositoryStorage),
		withRequestOpts(options...),
	)
}

// ScheduleStorageMoveForSnippetOptions represents the available
// ScheduleStorageMoveForSnippet() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#schedule-a-repository-storage-move-for-a-snippet
type ScheduleStorageMoveForSnippetOptions struct {
	DestinationStorageName *string `url:"destination_storage_name,omitempty" json:"destination_storage_name,omitempty"`
}

// ScheduleStorageMoveForSnippet schedule a repository to be moved for a snippet.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#schedule-a-repository-storage-move-for-a-snippet
func (s SnippetRepositoryStorageMoveService) ScheduleStorageMoveForSnippet(snippet int64, opts ScheduleStorageMoveForSnippetOptions, options ...RequestOptionFunc) (*SnippetRepositoryStorageMove, *Response, error) {
	return do[*SnippetRepositoryStorageMove](s.client,
		withMethod(http.MethodPost),
		withPath("snippets/%d/repository_storage_moves", snippet),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// ScheduleAllSnippetStorageMovesOptions represents the available
// ScheduleAllStorageMoves() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#schedule-repository-storage-moves-for-all-snippets-on-a-storage-shard
type ScheduleAllSnippetStorageMovesOptions struct {
	SourceStorageName      *string `url:"source_storage_name,omitempty" json:"source_storage_name,omitempty"`
	DestinationStorageName *string `url:"destination_storage_name,omitempty" json:"destination_storage_name,omitempty"`
}

// ScheduleAllStorageMoves schedules all snippet repositories to be moved.
//
// GitLab API docs:
// https://docs.gitlab.com/api/snippet_repository_storage_moves/#schedule-repository-storage-moves-for-all-snippets-on-a-storage-shard
func (s SnippetRepositoryStorageMoveService) ScheduleAllStorageMoves(opts ScheduleAllSnippetStorageMovesOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("snippet_repository_storage_moves"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
	return resp, err
}
