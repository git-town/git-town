package requests

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"strings"
)

// BodyGetter provides a Builder with a source for a request body.
type BodyGetter = func() (io.ReadCloser, error)

// BodyReader is a BodyGetter that returns an io.Reader.
func BodyReader(r io.Reader) BodyGetter {
	return func() (io.ReadCloser, error) {
		if rc, ok := r.(io.ReadCloser); ok {
			return rc, nil
		}
		return rc(r), nil
	}
}

// BodyWriter is a BodyGetter that pipes writes into a request body.
func BodyWriter(f func(w io.Writer) error) BodyGetter {
	return func() (io.ReadCloser, error) {
		r, w := io.Pipe()
		go func() {
			var err error
			defer func() {
				w.CloseWithError(err)
			}()
			err = f(w)
		}()
		return r, nil
	}
}

// BodyBytes is a BodyGetter that returns the provided raw bytes.
func BodyBytes(b []byte) BodyGetter {
	return func() (io.ReadCloser, error) {
		return rc(bytes.NewReader(b)), nil
	}
}

// BodySerializer is a BodyGetter
// that uses the provided [Serializer]
// to build the body of a request from v.
func BodySerializer(s Serializer, v any) BodyGetter {
	return func() (io.ReadCloser, error) {
		b, err := s(v)
		if err != nil {
			return nil, err
		}
		return rc(bytes.NewReader(b)), nil
	}
}

// BodyJSON is a [BodySerializer]
// that uses [JSONSerializer] to marshal the object.
func BodyJSON(v any) BodyGetter {
	return BodySerializer(JSONSerializer, v)
}

// BodyForm is a BodyGetter that builds an encoded form body.
func BodyForm(data url.Values) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return rc(strings.NewReader(data.Encode())), nil
	}
}

// BodyFile is a BodyGetter that reads the provided file path.
func BodyFile(name string) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return os.Open(name)
	}
}
