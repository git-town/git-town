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
	return do[[]*Project](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "projects", Search: query}),
		withRequestOpts(options...),
	)
}

// ProjectsByGroup searches the expression within projects for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#group-search-api
func (s *SearchService) ProjectsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	return do[[]*Project](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "projects", Search: query}),
		withRequestOpts(options...),
	)
}

// Issues searches the expression within issues
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues
func (s *SearchService) Issues(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	return do[[]*Issue](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "issues", Search: query}),
		withRequestOpts(options...),
	)
}

// IssuesByGroup searches the expression within issues for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues-1
func (s *SearchService) IssuesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	return do[[]*Issue](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "issues", Search: query}),
		withRequestOpts(options...),
	)
}

// IssuesByProject searches the expression within issues for
// the specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-issues-2
func (s *SearchService) IssuesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Issue, *Response, error) {
	return do[[]*Issue](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "issues", Search: query}),
		withRequestOpts(options...),
	)
}

// MergeRequests searches the expression within merge requests
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests
func (s *SearchService) MergeRequests(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	return do[[]*MergeRequest](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "merge_requests", Search: query}),
		withRequestOpts(options...),
	)
}

// MergeRequestsByGroup searches the expression within merge requests for
// the specified group
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests-1
func (s *SearchService) MergeRequestsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	return do[[]*MergeRequest](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "merge_requests", Search: query}),
		withRequestOpts(options...),
	)
}

// MergeRequestsByProject searches the expression within merge requests for
// the specified project
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-merge_requests-2
func (s *SearchService) MergeRequestsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	return do[[]*MergeRequest](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "merge_requests", Search: query}),
		withRequestOpts(options...),
	)
}

// Milestones searches the expression within milestones
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones
func (s *SearchService) Milestones(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	return do[[]*Milestone](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "milestones", Search: query}),
		withRequestOpts(options...),
	)
}

// MilestonesByGroup searches the expression within milestones for
// the specified group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones-1
func (s *SearchService) MilestonesByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	return do[[]*Milestone](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "milestones", Search: query}),
		withRequestOpts(options...),
	)
}

// MilestonesByProject searches the expression within milestones for
// the specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-milestones-2
func (s *SearchService) MilestonesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Milestone, *Response, error) {
	return do[[]*Milestone](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "milestones", Search: query}),
		withRequestOpts(options...),
	)
}

// SnippetTitles searches the expression within snippet titles
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-snippet_titles
func (s *SearchService) SnippetTitles(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Snippet, *Response, error) {
	return do[[]*Snippet](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "snippet_titles", Search: query}),
		withRequestOpts(options...),
	)
}

// NotesByProject searches the expression within notes for the specified
// project
//
// GitLab API docs: // https://docs.gitlab.com/api/search/#scope-notes
func (s *SearchService) NotesByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Note, *Response, error) {
	return do[[]*Note](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "notes", Search: query}),
		withRequestOpts(options...),
	)
}

// WikiBlobs searches the expression within all wiki blobs
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs
func (s *SearchService) WikiBlobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	return do[[]*Wiki](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "wiki_blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// WikiBlobsByGroup searches the expression within wiki blobs for
// specified group
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs-1
func (s *SearchService) WikiBlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	return do[[]*Wiki](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "wiki_blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// WikiBlobsByProject searches the expression within wiki blobs for
// the specified project
//
// GitLab API docs:
// https://docs.gitlab.com/api/search/#scope-wiki_blobs-2
func (s *SearchService) WikiBlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Wiki, *Response, error) {
	return do[[]*Wiki](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "wiki_blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// Commits searches the expression within all commits
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits
func (s *SearchService) Commits(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "commits", Search: query}),
		withRequestOpts(options...),
	)
}

// CommitsByGroup searches the expression within commits for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits-1
func (s *SearchService) CommitsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "commits", Search: query}),
		withRequestOpts(options...),
	)
}

// CommitsByProject searches the expression within commits for the
// specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-commits-2
func (s *SearchService) CommitsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "commits", Search: query}),
		withRequestOpts(options...),
	)
}

// Blob represents a single blob.
type Blob struct {
	Basename  string `json:"basename"`
	Data      string `json:"data"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	ID        string `json:"id"`
	Ref       string `json:"ref"`
	Startline int64  `json:"startline"`
	ProjectID int64  `json:"project_id"`
}

// Blobs searches the expression within all blobs
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs
func (s *SearchService) Blobs(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	return do[[]*Blob](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// BlobsByGroup searches the expression within blobs for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs-1
func (s *SearchService) BlobsByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	return do[[]*Blob](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// BlobsByProject searches the expression within blobs for the specified
// project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-blobs-2
func (s *SearchService) BlobsByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*Blob, *Response, error) {
	return do[[]*Blob](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "blobs", Search: query}),
		withRequestOpts(options...),
	)
}

// Users searches the expression within all users
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users
func (s *SearchService) Users(query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	return do[[]*User](s.client,
		withPath("search"),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "users", Search: query}),
		withRequestOpts(options...),
	)
}

// UsersByGroup searches the expression within users for the specified
// group
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users-1
func (s *SearchService) UsersByGroup(gid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	if opt == nil {
		opt = &SearchOptions{}
	}
	return do[[]*User](s.client,
		withPath("groups/%s/-/search", GroupID{gid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "users", Search: query}),
		withRequestOpts(options...),
	)
}

// UsersByProject searches the expression within users for the
// specified project
//
// GitLab API docs: https://docs.gitlab.com/api/search/#scope-users-2
func (s *SearchService) UsersByProject(pid any, query string, opt *SearchOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	return do[[]*User](s.client,
		withPath("projects/%s/-/search", ProjectID{pid}),
		withAPIOpts(&searchOptions{SearchOptions: *opt, Scope: "users", Search: query}),
		withRequestOpts(options...),
	)
}
