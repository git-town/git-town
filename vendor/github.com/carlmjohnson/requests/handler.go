package requests

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// ResponseHandler is used to validate or handle the response to a request.
type ResponseHandler = func(*http.Response) error

// ChainHandlers allows for the composing of validators or response handlers.
func ChainHandlers(handlers ...ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		for _, h := range handlers {
			if h == nil {
				continue
			}
			if err := h(r); err != nil {
				return err
			}
		}
		return nil
	}
}

func consumeBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}

// ToDeserializer decodes a response into v using a [Deserializer].
func ToDeserializer(d Deserializer, v any) ResponseHandler {
	return func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = d(data, v); err != nil {
			return err
		}
		return nil
	}
}

// ToJSON decodes a response as a JSON object.
//
// It uses [JSONDeserializer] to unmarshal the object.
func ToJSON(v any) ResponseHandler {
	return ToDeserializer(JSONDeserializer, v)
}

// ToString writes the response body to the provided string pointer.
func ToString(sp *string) ResponseHandler {
	return func(res *http.Response) error {
		var buf strings.Builder
		_, err := io.Copy(&buf, res.Body)
		if err == nil {
			*sp = buf.String()
		}
		return err
	}
}

// ToBytesBuffer writes the response body to the provided bytes.Buffer.
func ToBytesBuffer(buf *bytes.Buffer) ResponseHandler {
	return func(res *http.Response) error {
		_, err := io.Copy(buf, res.Body)
		return err
	}
}

// ToBufioReader takes a callback which wraps the response body in a bufio.Reader.
func ToBufioReader(f func(r *bufio.Reader) error) ResponseHandler {
	return func(res *http.Response) error {
		return f(bufio.NewReader(res.Body))
	}
}

// ToBufioScanner takes a callback which wraps the response body in a bufio.Scanner.
func ToBufioScanner(f func(r *bufio.Scanner) error) ResponseHandler {
	return func(res *http.Response) error {
		return f(bufio.NewScanner(res.Body))
	}
}

// ToHTML parses the page with x/net/html.Parse.
//
// Deprecated: Use reqhtml.To.
func ToHTML(n *html.Node) ResponseHandler {
	return ToBufioReader(func(r *bufio.Reader) error {
		n2, err := html.Parse(r)
		if err != nil {
			return err
		}
		*n = *n2
		return nil
	})
}

// ToWriter copies the response body to w.
func ToWriter(w io.Writer) ResponseHandler {
	return ToBufioReader(func(r *bufio.Reader) error {
		_, err := io.Copy(w, r)

		return err
	})
}

// ToFile writes the response body at the provided file path.
// The file and its parent directories are created automatically.
func ToFile(name string) ResponseHandler {
	return func(res *http.Response) error {
		_ = os.MkdirAll(filepath.Dir(name), 0777)

		f, err := os.Create(name)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, res.Body)
		return err
	}
}

// ToHeaders is an alias for backwards compatibility.
//
// Deprecated: Use CopyHeaders
var ToHeaders = CopyHeaders
