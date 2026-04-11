package bitbucket

import "fmt"

type Downloads struct {
	c *Client
}

func (dl *Downloads) Create(do *DownloadsOptions) (interface{}, error) {
	urlStr := dl.c.requestUrl("/repositories/%s/%s/downloads", do.Owner, do.RepoSlug)

	if do.FileName != "" {
		if len(do.Files) > 0 {
			return nil, fmt.Errorf("can't specify both files and filename")
		}
		do.Files = []File{{
			Path: do.FileName,
			Name: do.FileName,
		}}
	}
	return dl.c.executeFileUpload("POST", urlStr, do.Files, []string{}, make(map[string]string), do.ctx)
}

func (dl *Downloads) List(do *DownloadsOptions) (interface{}, error) {
	urlStr := dl.c.requestUrl("/repositories/%s/%s/downloads", do.Owner, do.RepoSlug)
	return dl.c.executePaginated("GET", urlStr, "", nil)
}
