// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TransferRepoOption options when transfer a repository's ownership
type TransferRepoOption struct {
	// required: true
	NewOwner string `json:"new_owner"`
	// ID of the team or teams to add to the repository. Teams can only be added to organization-owned repositories.
	TeamIDs *[]int64 `json:"team_ids"`
}

// TransferRepo transfers the ownership of a repository
func (c *Client) TransferRepo(owner, reponame string, opt TransferRepoOption) (*Repository, error) {
	if err := c.CheckServerVersionConstraint(">=1.12.0"); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/transfer", owner, reponame), jsonHeader, bytes.NewReader(body), repo)
}
