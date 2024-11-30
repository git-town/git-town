package requests

import "net/http"

// CheckRedirectPolicy is a function suitable for use as CheckRedirect on an http.Client.
type CheckRedirectPolicy = func(req *http.Request, via []*http.Request) error

// MaxFollow returns a CheckRedirectPolicy that follows a maximum of n redirects.
func MaxFollow(n int) CheckRedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) > n {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

// NoFollow is a CheckRedirectPolicy that does not follow redirects.
var NoFollow CheckRedirectPolicy = MaxFollow(0)
