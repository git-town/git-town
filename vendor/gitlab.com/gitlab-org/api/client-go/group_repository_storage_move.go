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
	GroupRepositoryStorageMoveServiceInterface interface {
		RetrieveAllStorageMoves(opts RetrieveAllGroupStorageMovesOptions, options ...RequestOptionFunc) ([]*GroupRepositoryStorageMove, *Response, error)
		RetrieveAllStorageMovesForGroup(group int64, opts RetrieveAllGroupStorageMovesOptions, options ...RequestOptionFunc) ([]*GroupRepositoryStorageMove, *Response, error)
		GetStorageMove(repositoryStorage int64, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error)
		GetStorageMoveForGroup(group int64, repositoryStorage int64, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error)
		ScheduleStorageMoveForGroup(group int64, opts ScheduleStorageMoveForGroupOptions, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error)
		ScheduleAllStorageMoves(opts ScheduleAllGroupStorageMovesOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupRepositoryStorageMoveService handles communication with the
	// group repositories related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_repository_storage_moves/
	GroupRepositoryStorageMoveService struct {
		client *Client
	}
)

// GroupRepositoryStorageMove represents the status of a repository move.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/
type GroupRepositoryStorageMove struct {
	ID                     int64            `json:"id"`
	CreatedAt              *time.Time       `json:"created_at"`
	State                  string           `json:"state"`
	SourceStorageName      string           `json:"source_storage_name"`
	DestinationStorageName string           `json:"destination_storage_name"`
	Group                  *RepositoryGroup `json:"group"`
}

type RepositoryGroup struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

// RetrieveAllGroupStorageMovesOptions represents the available
// RetrieveAllStorageMoves() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#retrieve-all-group-repository-storage-moves
type RetrieveAllGroupStorageMovesOptions struct {
	ListOptions
}

// RetrieveAllStorageMoves retrieves all group repository storage moves
// accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#retrieve-all-group-repository-storage-moves
func (g GroupRepositoryStorageMoveService) RetrieveAllStorageMoves(opts RetrieveAllGroupStorageMovesOptions, options ...RequestOptionFunc) ([]*GroupRepositoryStorageMove, *Response, error) {
	return do[[]*GroupRepositoryStorageMove](g.client,
		withPath("group_repository_storage_moves"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// RetrieveAllStorageMovesForGroup retrieves all repository storage moves for
// a single group accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#retrieve-all-repository-storage-moves-for-a-single-group
func (g GroupRepositoryStorageMoveService) RetrieveAllStorageMovesForGroup(group int64, opts RetrieveAllGroupStorageMovesOptions, options ...RequestOptionFunc) ([]*GroupRepositoryStorageMove, *Response, error) {
	return do[[]*GroupRepositoryStorageMove](g.client,
		withPath("groups/%d/repository_storage_moves", group),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// GetStorageMove gets a single group repository storage move.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#get-a-single-group-repository-storage-move
func (g GroupRepositoryStorageMoveService) GetStorageMove(repositoryStorage int64, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error) {
	return do[*GroupRepositoryStorageMove](g.client,
		withPath("group_repository_storage_moves/%d", repositoryStorage),
		withRequestOpts(options...),
	)
}

// GetStorageMoveForGroup gets a single repository storage move for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#get-a-single-repository-storage-move-for-a-group
func (g GroupRepositoryStorageMoveService) GetStorageMoveForGroup(group int64, repositoryStorage int64, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error) {
	return do[*GroupRepositoryStorageMove](g.client,
		withPath("groups/%d/repository_storage_moves/%d", group, repositoryStorage),
		withRequestOpts(options...),
	)
}

// ScheduleStorageMoveForGroupOptions represents the available
// ScheduleStorageMoveForGroup() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#schedule-a-repository-storage-move-for-a-group
type ScheduleStorageMoveForGroupOptions struct {
	DestinationStorageName *string `url:"destination_storage_name,omitempty" json:"destination_storage_name,omitempty"`
}

// ScheduleStorageMoveForGroup schedule a repository to be moved for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#schedule-a-repository-storage-move-for-a-group
func (g GroupRepositoryStorageMoveService) ScheduleStorageMoveForGroup(group int64, opts ScheduleStorageMoveForGroupOptions, options ...RequestOptionFunc) (*GroupRepositoryStorageMove, *Response, error) {
	return do[*GroupRepositoryStorageMove](g.client,
		withMethod(http.MethodPost),
		withPath("groups/%d/repository_storage_moves", group),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// ScheduleAllGroupStorageMovesOptions represents the available
// ScheduleAllStorageMoves() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#schedule-repository-storage-moves-for-all-groups-on-a-storage-shard
type ScheduleAllGroupStorageMovesOptions struct {
	SourceStorageName      *string `url:"source_storage_name,omitempty" json:"source_storage_name,omitempty"`
	DestinationStorageName *string `url:"destination_storage_name,omitempty" json:"destination_storage_name,omitempty"`
}

// ScheduleAllStorageMoves schedules all group repositories to be moved.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_repository_storage_moves/#schedule-repository-storage-moves-for-all-groups-on-a-storage-shard
func (g GroupRepositoryStorageMoveService) ScheduleAllStorageMoves(opts ScheduleAllGroupStorageMovesOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](g.client,
		withMethod(http.MethodPost),
		withPath("group_repository_storage_moves"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
	return resp, err
}
