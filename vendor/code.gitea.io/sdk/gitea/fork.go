// Copyright 2016 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ListForksOptions options for listing repository's forks
type ListForksOptions struct {
	ListOptions
}

// ListForks list a repository's forks
func (c *Client) ListForks(user string, repo string, opt ListForksOptions) ([]*Repository, error) {
	opt.setDefaults()
	forks := make([]*Repository, opt.PageSize)
	return forks, c.getParsedResponse("GET",
		fmt.Sprintf("/repos/%s/%s/forks?%s", user, repo, opt.getURLQuery().Encode()),
		nil, nil, &forks)
}

// CreateForkOption options for creating a fork
type CreateForkOption struct {
	// organization name, if forking into an organization
	Organization *string `json:"organization"`
}

// CreateFork create a fork of a repository
func (c *Client) CreateFork(user, repo string, form CreateForkOption) (*Repository, error) {
	body, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	fork := new(Repository)
	return fork, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/forks", user, repo), jsonHeader, bytes.NewReader(body), &fork)
}
