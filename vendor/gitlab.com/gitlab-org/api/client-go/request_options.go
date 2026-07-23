//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// RequestOptionFunc can be passed to all API requests to customize the API request.
type RequestOptionFunc func(*retryablehttp.Request) error

// WithContext runs the request with the provided context
func WithContext(ctx context.Context) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		newCtx := copyContextValues(req.Context(), ctx)

		*req = *req.WithContext(newCtx)
		return nil
	}
}

// copyContextValues copy some context key and values in old context
func copyContextValues(oldCtx context.Context, newCtx context.Context) context.Context {
	checkRetry := checkRetryFromContext(oldCtx)

	if checkRetry != nil {
		newCtx = contextWithCheckRetry(newCtx, checkRetry)
	}

	return newCtx
}

// WithHeader takes a header name and value and appends it to the request headers.
func WithHeader(name, value string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		req.Header.Set(name, value)
		return nil
	}
}

// WithHeaders takes a map of header name/value pairs and appends them to the
// request headers.
func WithHeaders(headers map[string]string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	}
}

// WithKeysetPaginationParameters takes a "next" link from the Link header of a
// response to a keyset-based paginated request and modifies the values of each
// query parameter in the request with its corresponding response parameter.
func WithKeysetPaginationParameters(nextLink string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		nextURL, err := url.Parse(nextLink)
		if err != nil {
			return err
		}
		q := req.URL.Query()
		for k, values := range nextURL.Query() {
			q.Del(k)
			for _, v := range values {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

// WithOffsetPaginationParameters takes a page number and modifies the request
// to use that page for offset-based pagination, overriding any existing page value.
func WithOffsetPaginationParameters(page int64) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		q := req.URL.Query()
		q.Del("page")
		q.Add("page", strconv.FormatInt(page, 10))
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

// withGraphQLPaginationParamters takes a PageInfo from a GraphQL response and
// modifies the request to use that cursor for GraphQL pagination, overriding
// any existing "after" variable.
//
// GraphQL API docs:
// https://docs.gitlab.com/development/graphql_guide/pagination/
func withGraphQLPaginationParamters(pi PageInfo) RequestOptionFunc {
	if !pi.HasNextPage {
		return nil
	}

	return func(req *retryablehttp.Request) error {
		var q GraphQLQuery

		data, err := req.BodyBytes()
		if err != nil {
			return fmt.Errorf("reading request body failed: %w", err)
		}

		if err := json.Unmarshal(data, &q); err != nil {
			return fmt.Errorf("decoding request body failed: %w", err)
		}

		if q.Variables == nil {
			q.Variables = make(map[string]any)
		}

		q.Variables["after"] = pi.EndCursor

		data, err = json.Marshal(q)
		if err != nil {
			return fmt.Errorf("encoding request body failed: %w", err)
		}

		return req.SetBody(data)
	}
}

// WithNext returns a RequestOptionFunc that configures the next page of a paginated
// request based on pagination metadata from a previous response. It automatically
// detects and handles all three pagination styles used by GitLab's APIs:
//
//   - GraphQL cursor pagination: Uses PageInfo.EndCursor with the "after" variable
//   - REST keyset pagination: Extracts parameters from the "next" link header
//   - REST offset pagination: Uses the NextPage number with "page" parameter
//
// If multiple pagination styles are present in the response, keyset/cursor pagination
// is preferred over offset pagination for better performance and consistency.
//
// The boolean return value indicates whether more pages are available, similar to
// the comma-ok idiom used for map accesses. When false, the returned
// RequestOptionFunc is nil.
func WithNext(resp *Response) (RequestOptionFunc, bool) {
	switch {
	case resp.PageInfo != nil:
		return withGraphQLPaginationParamters(*resp.PageInfo), resp.PageInfo.HasNextPage

	case resp.NextLink != "":
		return WithKeysetPaginationParameters(resp.NextLink), true

	case resp.NextPage != 0:
		return WithOffsetPaginationParameters(resp.NextPage), true

	default:
		return nil, false
	}
}

// WithSudo takes either a username or user ID and sets the Sudo request header.
func WithSudo(uid any) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		user, err := parseID(uid)
		if err != nil {
			return err
		}
		req.Header.Set("Sudo", user)
		return nil
	}
}

// WithToken takes a token which is then used when making this one request.
func WithToken(authType AuthType, token string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		switch authType {
		case JobToken:
			req.Header.Set("Job-Token", token)
		case OAuthToken:
			req.Header.Set("Authorization", "Bearer "+token)
		case PrivateToken:
			req.Header.Set("Private-Token", token)
		}
		return nil
	}
}

// WithRequestRetry takes a `retryablehttp.CheckRetry` which is then used when making this one request.
func WithRequestRetry(checkRetry retryablehttp.CheckRetry) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		// Store checkRetry to context
		ctx := contextWithCheckRetry(req.Context(), checkRetry)
		*req = *req.WithContext(ctx)
		return nil
	}
}
