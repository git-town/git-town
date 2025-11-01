package gitlab

import (
	"fmt"
	"net/http"
	"reflect"
)

type doConfig struct {
	method      string
	path        string
	apiOpts     any
	requestOpts []RequestOptionFunc
}

type doOption func(c *doConfig)

func withMethod(method string) doOption {
	return func(c *doConfig) {
		c.method = method
	}
}

func withPath(path string, args ...any) doOption {
	return func(c *doConfig) {
		as := make([]any, len(args))
		for i, a := range args {
			switch v := a.(type) {
			case string:
				as[i] = PathEscape(v)
			default:
				as[i] = v
			}
		}
		c.path = fmt.Sprintf(path, as...)
	}
}

func withAPIOpts(o any) doOption {
	return func(c *doConfig) {
		c.apiOpts = o
	}
}

func withRequestOpts(o ...RequestOptionFunc) doOption {
	return func(c *doConfig) {
		c.requestOpts = o
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
func do[T any](client *Client, opts ...doOption) (T, *Response, error) {
	// default config
	config := &doConfig{
		method:  http.MethodGet,
		apiOpts: nil,
	}

	// apply options to config
	for _, f := range opts {
		f(config)
	}

	req, err := client.NewRequest(config.method, config.path, config.apiOpts, config.requestOpts)
	if err != nil {
		var z T
		return z, nil, err
	}

	var (
		as   T
		resp *Response
	)
	if reflect.TypeOf(as) == reflect.TypeOf(none{}) {
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
