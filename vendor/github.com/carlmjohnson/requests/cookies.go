package requests

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// NewCookieJar returns a cookie jar using the standard public suffix list.
func NewCookieJar() http.CookieJar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	// As of Go 1.16, cookiejar.New err is hardcoded nil
	if err != nil {
		panic(err)
	}
	return jar
}
