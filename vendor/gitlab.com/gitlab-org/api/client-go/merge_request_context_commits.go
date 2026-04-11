package gitlab

import (
	"fmt"
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
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/context_commits", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var c []*Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
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
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/context_commits", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var c []*Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
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
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/context_commits", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodDelete, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
