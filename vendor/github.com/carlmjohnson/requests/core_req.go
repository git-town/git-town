package requests

import (
	"cmp"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests/internal/minitrue"
	"github.com/carlmjohnson/requests/internal/slicex"
)

// nopCloser is like io.NopCloser(),
// but it is a concrete type so we can strip it out
// before setting a body on a request.
// See https://github.com/carlmjohnson/requests/discussions/49
type nopCloser struct {
	io.Reader
}

func rc(r io.Reader) nopCloser {
	return nopCloser{r}
}

func (nopCloser) Close() error { return nil }

var _ io.ReadCloser = nopCloser{}

type requestBuilder struct {
	headers []multimap
	cookies []kvpair
	getBody BodyGetter
	method  string
}

func (rb *requestBuilder) Header(key string, values ...string) {
	rb.headers = append(rb.headers, multimap{key, values, false})
}

func (rb *requestBuilder) HeaderOptional(key string, values ...string) {
	rb.headers = append(rb.headers, multimap{key, values, true})
}

func (rb *requestBuilder) Cookie(name, value string) {
	rb.cookies = append(rb.cookies, kvpair{name, value})
}

func (rb *requestBuilder) Method(method string) {
	rb.method = method
}

func (rb *requestBuilder) Body(src BodyGetter) {
	rb.getBody = src
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *requestBuilder) Clone() *requestBuilder {
	rb2 := *rb
	slicex.Clip(&rb2.headers)
	slicex.Clip(&rb2.cookies)
	return &rb2
}

// Request builds a new http.Request with its context set.
func (rb *requestBuilder) Request(ctx context.Context, u *url.URL) (req *http.Request, err error) {
	var body io.Reader
	if rb.getBody != nil {
		if body, err = rb.getBody(); err != nil {
			return nil, err
		}
		if nopper, ok := body.(nopCloser); ok {
			body = nopper.Reader
		}
	}
	method := cmp.Or(rb.method,
		minitrue.Cond(rb.getBody == nil,
			http.MethodGet,
			http.MethodPost))

	req, err = http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.GetBody = rb.getBody

	for _, kv := range rb.headers {
		if !kv.optional {
			req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
		}
	}
	for _, kv := range rb.headers {
		if kv.optional &&
			req.Header.Get(kv.key) == "" &&
			cmp.Or(kv.values...) != "" {
			req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
		}
	}
	for _, kv := range rb.cookies {
		req.AddCookie(&http.Cookie{
			Name:  kv.key,
			Value: kv.value,
		})
	}
	return req, nil
}
