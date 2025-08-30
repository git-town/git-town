package requests

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests/internal/minitrue"
	"github.com/carlmjohnson/requests/internal/slicex"
)

// Builder is a convenient way to build, send, and handle HTTP requests.
// Builder has a fluent API with methods returning a pointer to the same
// struct, which allows for declaratively describing a request by method chaining.
//
// Builder can build a url.URL,
// build an http.Request,
// or handle a full http.Client request and response with validation.
//
// # Build a url.URL with Builder.URL
//
// Set the base URL by creating a new Builder with [requests.URL]
// or by calling [Builder.BaseURL]
// then customize it with
// [Builder.Scheme], [Builder.Host], [Builder.Hostf], [Builder.Path],
// [Builder.Pathf], [Builder.Param], and [Builder.ParamInt].
// [Builder.ParamOptional] can be used to add a query parameter
// only if it has not been otherwise set.
//
// # Build an http.Request with Builder.Request
//
// Set the method for a request with [Builder.Method]
// or use the [Builder.Delete], [Builder.Head], [Builder.Patch], [Builder.Post], and [Builder.Put] methods.
// By default, requests without a body are GET,
// and those with a body are POST.
//
// Set headers with [Builder.Header]
// or set conventional header keys with
// [Builder.Accept], [Builder.BasicAuth], [Builder.Bearer], [Builder.CacheControl],
// [Builder.ContentType], [Builder.Cookie], and [Builder.UserAgent].
// [Builder.HeaderOptional] can be used to add a header
// only if it has not been otherwise set.
//
// Set the body of the request, if any, with [Builder.Body]
// or use built in [Builder.BodyBytes], [Builder.BodyFile], [Builder.BodyForm],
// [Builder.BodyJSON], [Builder.BodyReader], or [Builder.BodyWriter].
//
// # Handle a request and response with Builder.Do or Builder.Fetch
//
// Set the http.Client to use for a request with [Builder.Client]
// and/or set an http.RoundTripper with [Builder.Transport].
//
// Add a response validator to the Builder with [Builder.AddValidator]
// or use the built in [Builder.CheckStatus], [Builder.CheckContentType],
// [Builder.CheckPeek], [Builder.CopyHeaders], and [Builder.ErrorJSON].
// If no validator has been added, Builder will use [DefaultValidator].
//
// Set a handler for a response with [Builder.Handle]
// or use the built in [Builder.ToHeaders], [Builder.ToJSON], [Builder.ToString],
// [Builder.ToBytesBuffer], or [Builder.ToWriter].
//
// [Builder.Fetch] creates an http.Request with [Builder.Request]
// and validates and handles it with [Builder.Do].
//
// # Other methods
//
// [Builder.Config] can be used to set several options on a Builder at once.
// [New] creates a new Builder and applies [Config] options to it.
//
// In many cases, it will be possible to set most options for an API endpoint
// in a Builder at the package or struct level
// and then call [Builder.Clone] in a function
// to add request specific details for the URL, parameters, headers, body, or handler.
//
// Errors returned by Builder methods will have an [ErrorKind] indicating their origin.
//
// The zero value of Builder is usable.
type Builder struct {
	ub         urlBuilder
	rb         requestBuilder
	cl         *http.Client
	rt         http.RoundTripper
	validators []ResponseHandler
	handler    ResponseHandler
}

// BaseURL sets the base URL that other URL methods modify.
// It is usually more convenient to use [URL] instead.
func (rb *Builder) BaseURL(baseurl string) *Builder {
	rb.ub.BaseURL(baseurl)
	return rb
}

// Scheme sets the scheme for a Builder's URL.
// It overrides the scheme set by BaseURL.
func (rb *Builder) Scheme(scheme string) *Builder {
	rb.ub.Scheme(scheme)
	return rb
}

// Host sets the host for a Builder's URL.
// It overrides the host set by BaseURL.
func (rb *Builder) Host(host string) *Builder {
	rb.ub.Host(host)
	return rb
}

// Path joins a path to a Builder's URL per the path joining rules of RFC 3986.
// If the path begins with /, it overrides any existing path.
// If the path begins with ./ or ../, the final path will be rewritten in its absolute form when creating a request.
func (rb *Builder) Path(path string) *Builder {
	rb.ub.Path(path)
	return rb
}

// Param sets a query parameter on a Builder's URL.
// It overwrites the existing values of a key.
func (rb *Builder) Param(key string, values ...string) *Builder {
	rb.ub.Param(key, values...)
	return rb
}

// ParamOptional sets a query parameter on a Builder's URL
// only if it is not set by some other call to Param or ParamOptional
// and one of the values is a non-blank string.
func (rb *Builder) ParamOptional(key string, values ...string) *Builder {
	rb.ub.ParamOptional(key, values...)
	return rb
}

// Header sets a header on a request. It overwrites the existing values of a key.
func (rb *Builder) Header(key string, values ...string) *Builder {
	rb.rb.Header(key, values...)
	return rb
}

// HeaderOptional sets a header on a request
// only if it has not already been set by another call to Header or HeaderOptional
// and one of the values is a non-blank string.
func (rb *Builder) HeaderOptional(key string, values ...string) *Builder {
	rb.rb.HeaderOptional(key, values...)
	return rb
}

// Cookie adds a cookie to a request.
// Unlike other headers, adding a cookie does not overwrite existing values.
func (rb *Builder) Cookie(name, value string) *Builder {
	rb.rb.Cookie(name, value)
	return rb
}

// Method sets the HTTP method for a request.
// By default, requests without a body are GET,
// and those with a body are POST.
func (rb *Builder) Method(method string) *Builder {
	rb.rb.Method(method)
	return rb
}

// Body sets the BodyGetter to use to build the body of a request.
// The provided BodyGetter is used as an http.Request.GetBody func.
// It implicitly sets method to POST.
func (rb *Builder) Body(src BodyGetter) *Builder {
	rb.rb.Body(src)
	return rb
}

// Client sets the http.Client to use for requests. If nil, it uses http.DefaultClient.
func (rb *Builder) Client(cl *http.Client) *Builder {
	rb.cl = cl
	return rb
}

// Transport sets the http.RoundTripper to use for requests.
// If set, it makes a shallow copy of the http.Client before modifying it.
func (rb *Builder) Transport(rt http.RoundTripper) *Builder {
	rb.rt = rt
	return rb
}

// AddValidator adds a response validator to the Builder.
// Adding a validator disables DefaultValidator.
// To disable all validation, just add nil.
func (rb *Builder) AddValidator(h ResponseHandler) *Builder {
	rb.validators = append(rb.validators, h)
	return rb
}

// Handle sets the response handler for a Builder.
// To use multiple handlers, use ChainHandlers.
func (rb *Builder) Handle(h ResponseHandler) *Builder {
	rb.handler = h
	return rb
}

// Config allows Builder to be extended by functions that set several options at once.
func (rb *Builder) Config(cfgs ...Config) *Builder {
	for _, cfg := range cfgs {
		cfg(rb)
	}
	return rb
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *Builder) Clone() *Builder {
	rb2 := *rb
	rb2.ub = *rb.ub.Clone()
	rb2.rb = *rb.rb.Clone()
	slicex.Clip(&rb2.validators)
	return &rb2
}

func joinerrs(a, b error) error {
	return fmt.Errorf("%w: %w", a, b)
}

// URL builds a *url.URL from the base URL and options set on the Builder.
// If a valid url.URL cannot be built,
// URL() nevertheless returns a new url.URL,
// so it is always safe to call u.String().
func (rb *Builder) URL() (u *url.URL, err error) {
	u, err = rb.ub.URL()
	if err != nil {
		return u, joinerrs(ErrURL, err)
	}
	return u, nil
}

// Request builds a new http.Request with its context set.
func (rb *Builder) Request(ctx context.Context) (req *http.Request, err error) {
	u, err := rb.URL()
	if err != nil {
		return nil, err
	}
	req, err = rb.rb.Request(ctx, u)
	if err != nil {
		return nil, joinerrs(ErrRequest, err)
	}
	return req, nil
}

// Do calls the underlying http.Client and validates and handles any resulting response. The response body is closed after all validators and the handler run.
func (rb *Builder) Do(req *http.Request) (err error) {
	cl := cmp.Or(rb.cl, http.DefaultClient)
	if rb.rt != nil {
		cl2 := *cl
		cl2.Transport = rb.rt
		cl = &cl2
	}
	validators := rb.validators
	if len(validators) == 0 {
		validators = []ResponseHandler{DefaultValidator}
	}
	h := minitrue.Cond(rb.handler != nil,
		rb.handler,
		consumeBody)

	code, err := do(cl, req, validators, h)
	switch code {
	case doOK:
		return nil
	case doConnect:
		err = joinerrs(ErrTransport, err)
	case doValidate:
		err = joinerrs(ErrValidator, err)
	case doHandle:
		err = joinerrs(ErrHandler, err)
	}
	return err
}

// Fetch builds a request, sends it, and handles the response.
func (rb *Builder) Fetch(ctx context.Context) (err error) {
	req, err := rb.Request(ctx)
	if err != nil {
		return err
	}
	return rb.Do(req)
}
