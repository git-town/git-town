package requests

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

// Record returns an http.RoundTripper that writes out its
// requests and their responses to text files in basepath.
// Requests are named according to a hash of their contents.
// Responses are named according to the request that made them.
//
// Deprecated: Use reqtest.Record.
func Record(rt http.RoundTripper, basepath string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while recording transport: %w", err)
			}
		}()
		_ = os.MkdirAll(basepath, 0755)
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		reqname, resname := buildName(b)
		name := filepath.Join(basepath, reqname)
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		if res, err = rt.RoundTrip(req); err != nil {
			return
		}
		b, err = httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		name = filepath.Join(basepath, resname)
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		return
	})
}

// Replay returns an http.RoundTripper that reads its
// responses from text files in basepath.
// Responses are looked up according to a hash of the request.
//
// Deprecated: Use reqtest.Replay.
func Replay(basepath string) Transport {
	return ReplayFS(os.DirFS(basepath))
}

var errNotFound = errors.New("response not found")

// ReplayFS returns an http.RoundTripper that reads its
// responses from text files in the fs.FS.
// Responses are looked up according to a hash of the request.
// Response file names may optionally be prefixed with comments for better human organization.
//
// Deprecated: Use reqtest.ReplayFS.
func ReplayFS(fsys fs.FS) Transport {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while replaying transport: %w", err)
			}
		}()
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		_, name := buildName(b)
		glob := "*" + name
		matches, err := fs.Glob(fsys, glob)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("%w: no replay file matches %q", errNotFound, glob)
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("ambiguous response: multiple replay files match %q", glob)
		}
		b, err = fs.ReadFile(fsys, matches[0])
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		return http.ReadResponse(r, req)
	})
}

func buildName(b []byte) (reqname, resname string) {
	h := md5.New()
	h.Write(b)
	s := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return s[:8] + ".req.txt", s[:8] + ".res.txt"
}

// Caching returns an http.RoundTripper that attempts to read its
// responses from text files in basepath. If the response is absent,
// it caches the result of issuing the request with rt in basepath.
// Requests are named according to a hash of their contents.
// Responses are named according to the request that made them.
//
// Deprecated: Use reqtest.Caching.
func Caching(rt http.RoundTripper, basepath string) Transport {
	replay := Replay(basepath).RoundTrip
	record := Record(rt, basepath).RoundTrip
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		res, err = replay(req)
		if errors.Is(err, errNotFound) {
			res, err = record(req)
		}
		return
	})
}
