package bitbucket

import (
	"errors"
	"fmt"
	"net/url"
)

//"github.com/k0kubun/pp"

type Repositories struct {
	c                  *Client
	PullRequests       *PullRequests
	Issues             *Issues
	Pipelines          *Pipelines
	Repository         *Repository
	Commits            *Commits
	Diff               *Diff
	BranchRestrictions *BranchRestrictions
	Webhooks           *Webhooks
	Downloads          *Downloads
	DeployKeys         *DeployKeys
	repositories
}

type RepositoriesRes struct {
	Page    int32
	Pagelen int32
	Size    int32
	Items   []Repository
}

func (r *Repositories) ListForAccount(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlPath := "/repositories"
	if ro.Owner != "" {
		urlPath += fmt.Sprintf("/%s", ro.Owner)
	}
	urlStr := r.c.requestUrl(urlPath)
	urlAsUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	q := urlAsUrl.Query()
	if ro.Role != "" {
		q.Set("role", ro.Role)
	}
	if ro.Keyword != nil && *ro.Keyword != "" {
		// https://developer.atlassian.com/cloud/bitbucket/rest/intro/#operators
		query := fmt.Sprintf("full_name ~ \"%s\"", *ro.Keyword)
		q.Set("q", query)
	}
	urlAsUrl.RawQuery = q.Encode()
	urlStr = urlAsUrl.String()
	repos, err := r.c.executePaginated("GET", urlStr, "", ro.Page)
	if err != nil {
		return nil, err
	}
	return decodeRepositories(repos)
}

// Deprecated: Use ListForAccount instead
func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	return r.ListForAccount(ro)
}

// Return all repositories that belong to a project
func (r *Repositories) ListProject(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlPath := r.c.requestUrl("/repositories")
	urlPath += fmt.Sprintf("/%s/?q=project.key=\"%s\"", ro.Owner, ro.Project)
	repos, err := r.c.executePaginated("GET", urlPath, "", nil)
	if err != nil {
		return nil, err
	}
	return decodeRepositories(repos)
}

func (r *Repositories) ListPublic() (*RepositoriesRes, error) {
	urlStr := r.c.requestUrl("/repositories/")
	repos, err := r.c.executePaginated("GET", urlStr, "", nil)
	if err != nil {
		return nil, err
	}
	return decodeRepositories(repos)
}

func decodeRepositories(reposResponse interface{}) (*RepositoriesRes, error) {
	reposResponseMap, ok := reposResponse.(map[string]interface{})
	if !ok {
		return nil, errors.New("Not a valid format")
	}

	repoArray := reposResponseMap["values"].([]interface{})
	var repos []Repository
	for _, repoEntry := range repoArray {
		repo, err := decodeRepository(repoEntry)
		if err == nil {
			repos = append(repos, *repo)
		}
	}

	page, ok := reposResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := reposResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	size, ok := reposResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	repositories := RepositoriesRes{
		Page:    int32(page),
		Pagelen: int32(pagelen),
		Size:    int32(size),
		Items:   repos,
	}
	return &repositories, nil
}
