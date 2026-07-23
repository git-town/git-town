package bitbucket

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type BranchRestrictions struct {
	c *Client

	ID      int
	Pattern string
	Kind    string
	Value   *int
}

func (b *BranchRestrictions) Gets(bo *BranchRestrictionsOptions) (interface{}, error) {
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.RepoSlug)
	return b.c.executePaginated("GET", urlStr, "", nil)
}

func (b *BranchRestrictions) Create(bo *BranchRestrictionsOptions) (*BranchRestrictions, error) {
	data, err := b.buildBranchRestrictionsBody(bo)
	if err != nil {
		return nil, err
	}
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.RepoSlug)
	response, err := b.c.executeWithContext("POST", urlStr, data, bo.ctx)
	if err != nil {
		return nil, err
	}

	return decodeBranchRestriction(response)
}

func (b *BranchRestrictions) Get(bo *BranchRestrictionsOptions) (*BranchRestrictions, error) {
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	response, err := b.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeBranchRestriction(response)
}

func (b *BranchRestrictions) Update(bo *BranchRestrictionsOptions) (interface{}, error) {
	data, err := b.buildBranchRestrictionsBody(bo)
	if err != nil {
		return nil, err
	}
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	response, err := b.c.execute("PUT", urlStr, data)
	if err != nil {
		return nil, err
	}

	return decodeBranchRestriction(response)
}

func (b *BranchRestrictions) Delete(bo *BranchRestrictionsOptions) (interface{}, error) {
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	return b.c.execute("DELETE", urlStr, "")
}

type branchRestrictionsBody struct {
	Type            string                        `json:"type"`
	Kind            string                        `json:"kind"`
	BranchMatchKind string                        `json:"branch_match_kind"`
	BranchType      string                        `json:"branch_type,omitempty"`
	Pattern         string                        `json:"pattern"`
	Value           interface{}                   `json:"value,omitempty"`
	Users           []branchRestrictionsBodyUser  `json:"users"`
	Groups          []branchRestrictionsBodyGroup `json:"groups"`
}

type branchRestrictionsBodyGroup struct {
	Name  string `json:"name"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	FullSlug string `json:"full_slug"`
	Slug     string `json:"slug"`
}

type branchRestrictionsBodyUser struct {
	Username     string `json:"username"`
	Website      string `json:"website"`
	Display_name string `json:"display_name"`
	UUID         string `json:"uuid"`
	Created_on   string `json:"created_on"`
	Links        struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Repositories struct {
			Href string `json:"href"`
		} `json:"repositories"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
		Followers struct {
			Href string `json:"href"`
		} `json:"followers"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		Following struct {
			Href string `json:"href"`
		} `json:"following"`
	} `json:"links"`
}

func (b *BranchRestrictions) buildBranchRestrictionsBody(bo *BranchRestrictionsOptions) (string, error) {
	users := make([]branchRestrictionsBodyUser, 0, len(bo.Users))
	groups := make([]branchRestrictionsBodyGroup, 0, len(bo.Groups))
	for _, u := range bo.Users {
		user := branchRestrictionsBodyUser{
			Username: u,
		}
		users = append(users, user)
	}
	for _, g := range bo.Groups {
		group := branchRestrictionsBodyGroup{
			Slug: g,
		}
		groups = append(groups, group)
	}

	branchMatchKind := bo.BranchMatchKind
	if branchMatchKind == "" {
		branchMatchKind = "glob"
	}

	body := branchRestrictionsBody{
		Type:            "branchrestriction",
		Kind:            bo.Kind,
		BranchMatchKind: branchMatchKind,
		BranchType:      bo.BranchType,
		Pattern:         bo.Pattern,
		Users:           users,
		Groups:          groups,
		Value:           bo.Value,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func decodeBranchRestriction(branchResponse interface{}) (*BranchRestrictions, error) {
	branchMap := branchResponse.(map[string]interface{})

	if branchMap["type"] == "error" {
		return nil, DecodeError(branchMap)
	}

	var branchRestriction = new(BranchRestrictions)
	err := mapstructure.Decode(branchMap, branchRestriction)
	if err != nil {
		return nil, err
	}
	return branchRestriction, nil
}
