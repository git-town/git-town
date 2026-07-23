package gitlab

import (
	"net/http"
)

type (
	// MergeRequestContextCommitsServiceInterface handles communication with the
	// merge request context commits related methods of the GitLab API.
	MergeRequestContextCommitsServiceInterface interface {
		// ListMergeRequestContextCommits gets a list of merge request context commits.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/merge_request_context_commits/#list-mr-context-commits
		ListMergeRequestContextCommits(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		// CreateMergeRequestContextCommits creates a list of merge request context
		// commits.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/merge_request_context_commits/#create-mr-context-commits
		CreateMergeRequestContextCommits(pid any, mergeRequest int64, opt *CreateMergeRequestContextCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error)
		// DeleteMergeRequestContextCommits deletes a list of merge request context
		// commits.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/merge_request_context_commits/#delete-mr-context-commits
		DeleteMergeRequestContextCommits(pid any, mergeRequest int64, opt *DeleteMergeRequestContextCommitsOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// MergeRequestContextCommitsService handles communication with the merge
	// request context commits related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/merge_request_context_commits/
	MergeRequestContextCommitsService struct {
		client *Client
	}
)

var _ MergeRequestContextCommitsServiceInterface = (*MergeRequestContextCommitsService)(nil)

func (s *MergeRequestContextCommitsService) ListMergeRequestContextCommits(pid any, mergeRequest int64, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withPath("projects/%s/merge_requests/%d/context_commits", ProjectID{pid}, mergeRequest),
		withRequestOpts(options...),
	)
}

// CreateMergeRequestContextCommitsOptions represents the available
// CreateMergeRequestContextCommits() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_context_commits/#create-mr-context-commits
type CreateMergeRequestContextCommitsOptions struct {
	Commits *[]string `url:"commits,omitempty" json:"commits,omitempty"`
}

func (s *MergeRequestContextCommitsService) CreateMergeRequestContextCommits(pid any, mergeRequest int64, opt *CreateMergeRequestContextCommitsOptions, options ...RequestOptionFunc) ([]*Commit, *Response, error) {
	return do[[]*Commit](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/context_commits", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteMergeRequestContextCommitsOptions represents the available
// DeleteMergeRequestContextCommits() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_context_commits/#delete-mr-context-commits
type DeleteMergeRequestContextCommitsOptions struct {
	Commits *[]string `url:"commits,omitempty" json:"commits,omitempty"`
}

func (s *MergeRequestContextCommitsService) DeleteMergeRequestContextCommits(pid any, mergeRequest int64, opt *DeleteMergeRequestContextCommitsOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/merge_requests/%d/context_commits", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
