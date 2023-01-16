// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// Milestone milestone is a collection of issues on one repository
type Milestone struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	State        StateType  `json:"state"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	Closed       *time.Time `json:"closed_at"`
	Deadline     *time.Time `json:"due_on"`
}

// ListMilestoneOption list milestone options
type ListMilestoneOption struct {
	ListOptions
	// open, closed, all
	State StateType
}

// QueryEncode turns options into querystring argument
func (opt *ListMilestoneOption) QueryEncode() string {
	query := opt.getURLQuery()
	if opt.State != "" {
		query.Add("state", string(opt.State))
	}
	return query.Encode()
}

// ListRepoMilestones list all the milestones of one repository
func (c *Client) ListRepoMilestones(owner, repo string, opt ListMilestoneOption) ([]*Milestone, error) {
	opt.setDefaults()
	milestones := make([]*Milestone, 0, opt.PageSize)

	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/milestones", owner, repo))
	link.RawQuery = opt.QueryEncode()
	return milestones, c.getParsedResponse("GET", link.String(), nil, nil, &milestones)
}

// GetMilestone get one milestone by repo name and milestone id
func (c *Client) GetMilestone(owner, repo string, id int64) (*Milestone, error) {
	milestone := new(Milestone)
	return milestone, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, id), nil, nil, milestone)
}

// CreateMilestoneOption options for creating a milestone
type CreateMilestoneOption struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"due_on"`
}

// CreateMilestone create one milestone with options
func (c *Client) CreateMilestone(owner, repo string, opt CreateMilestoneOption) (*Milestone, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	milestone := new(Milestone)
	return milestone, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/milestones", owner, repo), jsonHeader, bytes.NewReader(body), milestone)
}

// EditMilestoneOption options for editing a milestone
type EditMilestoneOption struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	State       *string    `json:"state"`
	Deadline    *time.Time `json:"due_on"`
}

// EditMilestone modify milestone with options
func (c *Client) EditMilestone(owner, repo string, id int64, opt EditMilestoneOption) (*Milestone, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	milestone := new(Milestone)
	return milestone, c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, id), jsonHeader, bytes.NewReader(body), milestone)
}

// DeleteMilestone delete one milestone by milestone id
func (c *Client) DeleteMilestone(owner, repo string, id int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/milestones/%d", owner, repo, id), nil, nil)
	return err
}
