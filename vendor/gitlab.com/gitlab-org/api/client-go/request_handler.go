package gitlab

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

type Pather interface {
	forPath() (string, error)
}

type ProjectID struct {
	Value any
}

func (i ProjectID) forPath() (string, error) {
	id, err := parseID(i.Value)
	if err != nil {
		return "", err
	}

	return PathEscape(id), nil
}

type GroupID struct {
	Value any
}

func (i GroupID) forPath() (string, error) {
	id, err := parseID(i.Value)
	if err != nil {
		return "", err
	}

	return PathEscape(id), nil
}

type RunnerID struct {
	Value any
}

func (i RunnerID) forPath() (string, error) {
	id, err := parseID(i.Value)
	if err != nil {
		return "", err
	}

	return PathEscape(id), nil
}

// UserID represents a user identifier for API paths. It accepts either a
// numeric user ID or a username string. If a username is provided with a
// leading "@" character (e.g., "@johndoe"), the "@" will be trimmed.
type UserID struct {
	Value any
}

func (i UserID) forPath() (string, error) {
	id, err := parseID(i.Value)
	if err != nil {
		return "", err
	}

	return PathEscape(strings.TrimPrefix(id, "@")), nil
}

type LabelID struct {
	Value any
}

func (i LabelID) forPath() (string, error) {
	id, err := parseID(i.Value)
	if err != nil {
		return "", err
	}

	return PathEscape(id), nil
}

type NoEscape struct {
	Value string
}

func (n NoEscape) forPath() (string, error) {
	return n.Value, nil
}

type doConfig struct {
	method      string
	path        string
	apiOpts     any
	requestOpts []RequestOptionFunc
	upload      *uploadConfig
}

type uploadConfig struct {
	content    io.Reader
	filename   string
	uploadType UploadType
}

type doOption func(c *doConfig) error

func withMethod(method string) doOption {
	return func(c *doConfig) error {
		c.method = method
		return nil
	}
}

func withPath(path string, args ...any) doOption {
	return func(c *doConfig) error {
		as := make([]any, len(args))
		for i, a := range args {
			switch v := a.(type) {
			case Pather:
				path, err := v.forPath()
				if err != nil {
					return err
				}
				as[i] = path
			case string:
				as[i] = PathEscape(v)
			default:
				as[i] = v
			}
		}
		c.path = fmt.Sprintf(path, as...)

		return nil
	}
}

func withAPIOpts(o any) doOption {
	return func(c *doConfig) error {
		c.apiOpts = o
		return nil
	}
}

func withRequestOpts(o ...RequestOptionFunc) doOption {
	return func(c *doConfig) error {
		c.requestOpts = o
		return nil
	}
}

func withUpload(content io.Reader, filename string, uploadType UploadType) doOption {
	return func(c *doConfig) error {
		c.upload = &uploadConfig{
			content:    content,
			filename:   filename,
			uploadType: uploadType,
		}
		return nil
	}
}

// none is a sentinel type to signal that a request performed with do does not return a value.
type none struct{}

// do constructs an API requests, performs it and processes the response.
//
// Use the opts to configure the request.
// If the response body shouldn't be handled, use the none sentinel type
// and ignore the first return argument.
//
// Example:
//
// // Get Request to return single *Agent:
// return do[*Agent](s.client,
//
//	withPath("projects/%s/cluster_agents/%d", project, id),
//	withRequestOpts(options...),
//
// )
//
// // Get Request to return multiple []*Agents
// return do[[]*Agent](s.client,
//
//	withPath("projects/%s/cluster_agents", project),
//	withRequestOpts(options...),
//
// )
//
// // Post Request to create Agent and return *Agents
// return do[*Agent](s.client,
//
//	withMethod(http.MethodPost),
//	withPath("projects/%s/cluster_agents", project),
//	withAPIOpts(opt),
//	withRequestOpts(options...),
//
// )
//
// // Delete Request that returns nothing:
// _, resp, err := do[none](s.client,
//
//	withMethod(http.MethodDelete),
//	withPath("projects/%s/cluster_agents/%d", project, id),
//	withRequestOpts(options...),
//
// )
//
// // Upload file Request:
// return do[*WikiAttachment](s.client,
//
//	withMethod(http.MethodPost),
//	withPath("projects/%s/wikis/attachments", project),
//	withUpload(content, filename, UploadFile),
//	withAPIOpts(opt),
//	withRequestOpts(options...),
//
// )
func do[T any](client *Client, opts ...doOption) (T, *Response, error) {
	// default config
	config := &doConfig{
		method:  http.MethodGet,
		apiOpts: nil,
	}

	// apply options to config
	for _, f := range opts {
		err := f(config)
		if err != nil {
			var z T
			return z, nil, err
		}
	}

	var (
		req *retryablehttp.Request
		err error
	)
	switch {
	case config.upload != nil:
		req, err = client.UploadRequest(
			config.method,
			config.path,
			config.upload.content,
			config.upload.filename,
			config.upload.uploadType,
			config.apiOpts,
			config.requestOpts,
		)
	default:
		req, err = client.NewRequest(config.method, config.path, config.apiOpts, config.requestOpts)
	}

	if err != nil {
		var z T
		return z, nil, err
	}

	var (
		as   T
		resp *Response
	)
	if reflect.TypeOf(as) == reflect.TypeFor[none]() {
		resp, err = client.Do(req, nil)
	} else {
		resp, err = client.Do(req, &as)
	}

	if err != nil {
		var z T
		return z, resp, err
	}

	return as, resp, nil
}
