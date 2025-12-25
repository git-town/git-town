package requests

import (
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
)

// Config allows Builder to be extended by setting several options at once.
// For example, a Config might set a Body and its ContentType.
type Config = func(rb *Builder)

// GzipConfig writes a gzip stream to its request body using a callback.
// It also sets the appropriate Content-Encoding header and automatically
// closes and the stream when the callback returns.
func GzipConfig(level int, h func(gw *gzip.Writer) error) Config {
	return func(rb *Builder) {
		rb.
			Header("Content-Encoding", "gzip").
			BodyWriter(func(w io.Writer) error {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return err
				}
				if err = h(gw); err != nil {
					gw.Close()
					return err
				}
				return gw.Close()
			})
	}
}

// TestServerConfig returns a Config
// which sets the Builder's BaseURL to s.URL
// and the Builder's Client to s.Client().
//
// Deprecated: Use reqtest.Server.
func TestServerConfig(s *httptest.Server) Config {
	return func(rb *Builder) {
		rb.
			BaseURL(s.URL).
			Client(s.Client())
	}
}

// BodyMultipart returns a Config
// that uses a multipart.Writer for the request body.
// If boundary is "", a multipart boundary is chosen at random.
// The content type of the request is set to multipart/form-data
// with the correct boundary.
// The multipart.Writer is automatically closed if the callback succeeds.
func BodyMultipart(boundary string, h func(multi *multipart.Writer) error) Config {
	return func(rb *Builder) {
		if boundary == "" {
			multi := multipart.NewWriter(nil)
			boundary = multi.Boundary()
		}
		rb.
			ContentType("multipart/form-data; boundary=" + boundary).
			BodyWriter(func(w io.Writer) error {
				multi := multipart.NewWriter(w)
				if err := multi.SetBoundary(boundary); err != nil {
					return fmt.Errorf("setting boundary: %w", err)
				}
				if err := h(multi); err != nil {
					return err
				}
				if err := multi.Close(); err != nil {
					return fmt.Errorf("closing multipart writer: %w", err)
				}
				return nil
			})
	}
}
