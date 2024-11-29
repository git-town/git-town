package requests

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
)

// DefaultValidator is the validator applied by Builder unless otherwise specified.
var DefaultValidator ResponseHandler = CheckStatus(
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNonAuthoritativeInfo,
	http.StatusNoContent,
)

// ResponseError is the error type produced by CheckStatus and CheckContentType.
type ResponseError http.Response

// Error fulfills the error interface.
func (se *ResponseError) Error() string {
	return fmt.Sprintf("response error for %s", se.Request.URL.Redacted())
}

// CheckStatus validates the response has an acceptable status code.
func CheckStatus(acceptStatuses ...int) ResponseHandler {
	return func(res *http.Response) error {
		for _, code := range acceptStatuses {
			if res.StatusCode == code {
				return nil
			}
		}

		return fmt.Errorf("%w: unexpected status: %d",
			(*ResponseError)(res), res.StatusCode)
	}
}

// HasStatusErr returns true if err is a ResponseError caused by any of the codes given.
func HasStatusErr(err error, codes ...int) bool {
	if err == nil {
		return false
	}
	if se := new(ResponseError); errors.As(err, &se) {
		for _, code := range codes {
			if se.StatusCode == code {
				return true
			}
		}
	}
	return false
}

// CheckContentType validates that a response has one of the given content type headers.
func CheckContentType(cts ...string) ResponseHandler {
	return func(res *http.Response) error {
		mt, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("%w: problem matching Content-Type",
				(*ResponseError)(res))
		}
		for _, ct := range cts {
			if mt == ct {
				return nil
			}
		}
		return fmt.Errorf("%w: unexpected Content-Type: %s",
			(*ResponseError)(res), mt)
	}
}

type bufioCloser struct {
	*bufio.Reader
	io.Closer
}

// CheckPeek wraps the body of a response in a bufio.Reader and
// gives f a peek at the first n bytes for validation.
func CheckPeek(n int, f func([]byte) error) ResponseHandler {
	return func(res *http.Response) error {
		// ensure buffer is at least minimum size
		buf := bufio.NewReader(res.Body)
		// ensure large peeks will fit in the buffer
		buf = bufio.NewReaderSize(buf, n)
		res.Body = &bufioCloser{
			buf,
			res.Body,
		}
		b, err := buf.Peek(n)
		if err != nil && err != io.EOF {
			return err
		}
		return f(b)
	}
}

// CopyHeaders copies the response headers to h.
func CopyHeaders(h map[string][]string) ResponseHandler {
	return func(res *http.Response) error {
		for k, v := range res.Header {
			h[k] = v
		}

		return nil
	}
}
