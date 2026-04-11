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
	"fmt"
	"net/http"
)

type (
	SearchServiceInterface interface {
		Projects(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Project, *Response, error)
		ProjectsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Project, *Response, error)
		Issues(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		IssuesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		IssuesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error)
		MergeRequests(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error)
		MergeRequestsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error)
		MergeRequestsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error)
		Milestones(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error)
		MilestonesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error)
		MilestonesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error)
		SnippetTitles(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Snippet, *Response, error)
		NotesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Note, *Response, error)
		WikiBlobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error)
		WikiBlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error)
		WikiBlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error)
		Commits(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		CommitsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		CommitsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		Blobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error)
		BlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error)
		BlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error)
		Users(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
		UsersByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
		UsersByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
	}

	// SearchService handles communication with the search related methods of the
	// GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/search/
	SearchService struct {
		client *Client
	}
)

var _ SearchServiceInterface = (*SearchService)(nil)

// SearchOptions represents the available options for all search methods.
//
// GitLab API docs: https://docs.gitlab.com/api/search/
type SearchOptions struct {
	ListOptions
	Ref *string `url:"ref,omitempty" json:"ref,omitempty"`
}

type searchOptions struct {
	SearchOptions
	Scope  string `url:"scope" json:"scope"`
	Search string `url:"search" json:"search"`
}

// Projects searches the expression within projects
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-projects
func (s *SearchService) Projects(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	var ps []*Project
	resp, err := s.search("projects", query, &ps, opt, options...)
	return ps, resp, err
}

// ProjectsByGroup searches the expression within projects for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#group-search-api
func (s *SearchService) ProjectsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	var ps []*Project
	resp, err := s.searchByGroup(gid, "projects", query, &ps, opt, options...)
	return ps, resp, err
}

// Issues searches the expression within issues
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues
func (s *SearchService) Issues(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	var is []*Issue
	resp, err := s.search("issues", query, &is, opt, options...)
	return is, resp, err
}

// IssuesByGroup searches the expression within issues for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues-1
func (s *SearchService) IssuesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	var is []*Issue
	resp, err := s.searchByGroup(gid, "issues", query, &is, opt, options...)
	return is, resp, err
}

// IssuesByProject searches the expression within issues for
// the specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues-2
func (s *SearchService) IssuesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	var is []*Issue
	resp, err := s.searchByProject(pid, "issues", query, &is, opt, options...)
	return is, resp, err
}

// MergeRequests searches the expression within merge requests
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests
func (s *SearchService) MergeRequests(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	var ms []*MergeRequest
	resp, err := s.search("merge_requests", query, &ms, opt, options...)
	return ms, resp, err
}

// MergeRequestsByGroup searches the expression within merge requests for
// the specified group
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests-1
func (s *SearchService) MergeRequestsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	var ms []*MergeRequest
	resp, err := s.searchByGroup(gid, "merge_requests", query, &ms, opt, options...)
	return ms, resp, err
}

// MergeRequestsByProject searches the expression within merge requests for
// the specified project
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests-2
func (s *SearchService) MergeRequestsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	var ms []*MergeRequest
	resp, err := s.searchByProject(pid, "merge_requests", query, &ms, opt, options...)
	return ms, resp, err
}

// Milestones searches the expression within milestones
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones
func (s *SearchService) Milestones(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	var ms []*Milestone
	resp, err := s.search("milestones", query, &ms, opt, options...)
	return ms, resp, err
}

// MilestonesByGroup searches the expression within milestones for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones-1
func (s *SearchService) MilestonesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	var ms []*Milestone
	resp, err := s.searchByGroup(gid, "milestones", query, &ms, opt, options...)
	return ms, resp, err
}

// MilestonesByProject searches the expression within milestones for
// the specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones-2
func (s *SearchService) MilestonesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	var ms []*Milestone
	resp, err := s.searchByProject(pid, "milestones", query, &ms, opt, options...)
	return ms, resp, err
}

// SnippetTitles searches the expression within snippet titles
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-snippet_titles
func (s *SearchService) SnippetTitles(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Snippet, *Response, error) {
	var ss []*Snippet
	resp, err := s.search("snippet_titles", query, &ss, opt, options...)
	return ss, resp, err
}

// NotesByProject searches the expression within notes for the specified
// project
//
// GitLab API docs: // https://docs.gitlab.com/api/search/#scope-notes
func (s *SearchService) NotesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Note, *Response, error) {
	var ns []*Note
	resp, err := s.searchByProject(pid, "notes", query, &ns, opt, options...)
	return ns, resp, err
}

// WikiBlobs searches the expression within all wiki blobs
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs
func (s *SearchService) WikiBlobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	var ws []*Wiki
	resp, err := s.search("wiki_blobs", query, &ws, opt, options...)
	return ws, resp, err
}

// WikiBlobsByGroup searches the expression within wiki blobs for
// specified group
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs-1
func (s *SearchService) WikiBlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	var ws []*Wiki
	resp, err := s.searchByGroup(gid, "wiki_blobs", query, &ws, opt, options...)
	return ws, resp, err
}

// WikiBlobsByProject searches the expression within wiki blobs for
// the specified project
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs-2
func (s *SearchService) WikiBlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	var ws []*Wiki
	resp, err := s.searchByProject(pid, "wiki_blobs", query, &ws, opt, options...)
	return ws, resp, err
}

// Commits searches the expression within all commits
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits
func (s *SearchService) Commits(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	var cs []*Commit
	resp, err := s.search("commits", query, &cs, opt, options...)
	return cs, resp, err
}

// CommitsByGroup searches the expression within commits for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits-1
func (s *SearchService) CommitsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	var cs []*Commit
	resp, err := s.searchByGroup(gid, "commits", query, &cs, opt, options...)
	return cs, resp, err
}

// CommitsByProject searches the expression within commits for the
// specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits-2
func (s *SearchService) CommitsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	var cs []*Commit
	resp, err := s.searchByProject(pid, "commits", query, &cs, opt, options...)
	return cs, resp, err
}

// Blob represents a single blob.
type Blob struct {
	Basename  string `json:"basename"`
	Data      string `json:"data"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	ID        string `json:"id"`
	Ref       string `json:"ref"`
	Startline int    `json:"startline"`
	ProjectID int    `json:"project_id"`
}

// Blobs searches the expression within all blobs
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs
func (s *SearchService) Blobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	var bs []*Blob
	resp, err := s.search("blobs", query, &bs, opt, options...)
	return bs, resp, err
}

// BlobsByGroup searches the expression within blobs for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs-1
func (s *SearchService) BlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	var bs []*Blob
	resp, err := s.searchByGroup(gid, "blobs", query, &bs, opt, options...)
	return bs, resp, err
}

// BlobsByProject searches the expression within blobs for the specified
// project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs-2
func (s *SearchService) BlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	var bs []*Blob
	resp, err := s.searchByProject(pid, "blobs", query, &bs, opt, options...)
	return bs, resp, err
}

// Users searches the expression within all users
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users
func (s *SearchService) Users(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	var ret []*User
	resp, err := s.search("users", query, &ret, opt, options...)
	return ret, resp, err
}

// UsersByGroup searches the expression within users for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users-1
func (s *SearchService) UsersByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	if opt == nil {
		opt = &SearchOptions{}
	}
	var ret []*User
	resp, err := s.searchByGroup(gid, "users", query, &ret, opt, options...)
	return ret, resp, err
}

// UsersByProject searches the expression within users for the
// specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users-2
func (s *SearchService) UsersByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	var ret []*User
	resp, err := s.searchByProject(pid, "users", query, &ret, opt, options...)
	return ret, resp, err
}

func (s *SearchService) search(scope, query string, result any, opt *SearchOptions, options ...RequestOptionFunc) (*Response, error) {
	opts := &searchOptions{SearchOptions: *opt, Scope: scope, Search: query}

	req, err := s.client.NewRequest(http.MethodGet, "search", opts, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, result)
}

func (s *SearchService) searchByGroup(gid any, scope, query string, result any, opt *SearchOptions, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/-/search", PathEscape(group))

	opts := &searchOptions{SearchOptions: *opt, Scope: scope, Search: query}

	req, err := s.client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, result)
}

func (s *SearchService) searchByProject(pid any, scope, query string, result any, opt *SearchOptions, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/-/search", PathEscape(project))

	opts := &searchOptions{SearchOptions: *opt, Scope: scope, Search: query}

	req, err := s.client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, result)
}
