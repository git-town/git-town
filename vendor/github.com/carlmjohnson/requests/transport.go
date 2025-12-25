package requests

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Transport is an alias of http.RoundTripper for documentation purposes.
type Transport = http.RoundTripper

// RoundTripFunc is an adaptor to use a function as an http.RoundTripper.
type RoundTripFunc func(req *http.Request) (res *http.Response, err error)

// RoundTrip implements http.RoundTripper.
func (rtf RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtf(r)
}

var _ Transport = RoundTripFunc(nil)

// ReplayString returns an http.RoundTripper that always responds with a
// request built from rawResponse. It is intended for use in one-off tests.
//
// Deprecated: Use reqtest.ReplayString.
func ReplayString(rawResponse string) Transport {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r := bufio.NewReader(strings.NewReader(rawResponse))
		res, err = http.ReadResponse(r, req)
		return
	})
}

// UserAgentTransport returns a wrapped http.RoundTripper that sets the User-Agent header on requests to s.
func UserAgentTransport(rt http.RoundTripper, s string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r2 := *req
		r2.Header = r2.Header.Clone()
		r2.Header.Set("User-Agent", s)
		return rt.RoundTrip(&r2)
	})
}

// PermitURLTransport returns a wrapped http.RoundTripper that rejects any requests whose URL doesn't match the provided regular expression string.
//
// PermitURLTransport will panic if the regexp does not compile.
func PermitURLTransport(rt http.RoundTripper, regex string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	re := regexp.MustCompile(regex)
	reErr := fmt.Errorf("requested URL not permitted by regexp: %s", regex)
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		if u := req.URL.String(); !re.MatchString(u) {
			return nil, reErr
		}
		return rt.RoundTrip(req)
	})
}

// LogTransport returns a wrapped http.RoundTripper
// that calls fn with details when a response has finished.
// A response is considered finished
// when the wrapper http.RoundTripper returns an error
// or the Response.Body is closed,
// whichever comes first.
// To simplify logging code,
// a nil *http.Response is replaced with a new http.Response.
func LogTransport(rt http.RoundTripper, fn func(req *http.Request, res *http.Response, err error, duration time.Duration)) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		start := time.Now()
		res, err = rt.RoundTrip(req)
		if err != nil {
			res2 := res
			if res == nil {
				res2 = new(http.Response)
			}
			fn(req, res2, err, time.Since(start))
			return
		}

		res.Body = closeLogger{res.Body, func() {
			fn(req, res, err, time.Since(start))
		}}
		return
	})
}

type closeLogger struct {
	io.ReadCloser
	fn func()
}

func (cl closeLogger) Close() error {
	cl.fn()
	return cl.ReadCloser.Close()
}

// DoerTransport converts a Doer into a Transport.
// It exists for compatibility with other libraries.
// A Doer is an interface with a Do method.
// Users should prefer Transport,
// because Do is the interface of http.Client
// which has higher level concerns.
func DoerTransport(cl interface {
	Do(req *http.Request) (*http.Response, error)
}) Transport {
	return RoundTripFunc(cl.Do)
}

// ErrorTransport always returns the specified error instead of connecting.
// It is intended for use in testing
// or to prevent accidental use of http.DefaultClient.
func ErrorTransport(err error) Transport {
	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, err
	})
}
