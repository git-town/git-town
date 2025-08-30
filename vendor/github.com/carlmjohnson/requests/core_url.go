package requests

import (
	"cmp"
	"net/url"

	"github.com/carlmjohnson/requests/internal/slicex"
)

type multimap struct {
	key      string
	values   []string
	optional bool
}

type kvpair struct {
	key, value string
}

type urlBuilder struct {
	baseurl      string
	scheme, host string
	paths        []string
	params       []multimap
}

func (ub *urlBuilder) BaseURL(baseurl string) {
	ub.baseurl = baseurl
}

func (ub *urlBuilder) Scheme(scheme string) {
	ub.scheme = scheme
}

func (ub *urlBuilder) Host(host string) {
	ub.host = host

}

func (ub *urlBuilder) Path(path string) {
	ub.paths = append(ub.paths, path)
}

func (ub *urlBuilder) Param(key string, values ...string) {
	ub.params = append(ub.params, multimap{key, values, false})
}

func (ub *urlBuilder) ParamOptional(key string, values ...string) {
	ub.params = append(ub.params, multimap{key, values, true})
}

func (ub *urlBuilder) Clone() *urlBuilder {
	ub2 := *ub
	slicex.Clip(&ub2.paths)
	slicex.Clip(&ub2.params)
	return &ub2
}

func (ub *urlBuilder) URL() (u *url.URL, err error) {
	u, err = url.Parse(ub.baseurl)
	if err != nil {
		return new(url.URL), err
	}
	u.Scheme = cmp.Or(
		ub.scheme,
		u.Scheme,
		"https",
	)
	u.Host = cmp.Or(ub.host, u.Host)
	for _, p := range ub.paths {
		u.Path = u.ResolveReference(&url.URL{Path: p}).Path
	}
	if len(ub.params) > 0 {
		q := u.Query()
		for _, kv := range ub.params {
			if !kv.optional {
				q[kv.key] = kv.values
			}
		}
		for _, kv := range ub.params {
			if kv.optional &&
				q.Get(kv.key) == "" &&
				cmp.Or(kv.values...) != "" {
				q[kv.key] = kv.values
			}
		}
		u.RawQuery = q.Encode()
	}
	// Reparsing, in case the path rewriting broke the URL
	u, err = url.Parse(u.String())
	if err != nil {
		return new(url.URL), err
	}
	return u, nil
}
