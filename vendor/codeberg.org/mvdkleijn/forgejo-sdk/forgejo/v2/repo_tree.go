// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
)

// GitEntry represents a git tree
type GitEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

// GitTreeResponse returns a git tree
type GitTreeResponse struct {
	SHA        string     `json:"sha"`
	URL        string     `json:"url"`
	Entries    []GitEntry `json:"tree"`
	Truncated  bool       `json:"truncated"`
	Page       int        `json:"page"`
	TotalCount int        `json:"total_count"`
}

// GetTreesOptions options for listing a repository's trees
type GetTreesOptions struct {
	Recursive bool
	ListOptions
}

// GetTrees downloads a file of repository, ref can be branch/tag/commit.
// e.g.: ref -> main, tree -> macaron.go(no leading slash)
func (c *Client) GetTrees(user, repo, ref string, opt GetTreesOptions) (*GitTreeResponse, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo, &ref); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()

	trees := new(GitTreeResponse)
	recInt := 0
	if opt.Recursive {
		recInt = 1
	}
	perPage := opt.PageSize // workaround for api endpoint using per_page instead of limit
	path := fmt.Sprintf("/repos/%s/%s/git/trees/%s?recursive=%d&per_page=%d&%s", user, repo, ref, recInt, perPage, opt.getURLQuery().Encode())

	resp, err := c.getParsedResponse("GET", path, nil, nil, &trees)

	return trees, resp, err
}
