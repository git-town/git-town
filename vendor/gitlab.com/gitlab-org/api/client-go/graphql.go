package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	// GraphQLAPIEndpoint defines the endpoint URI for the GraphQL backend
	GraphQLAPIEndpoint = "/api/graphql"
)

type (
	GraphQLInterface interface {
		Do(query GraphQLQuery, response any, options ...RequestOptionFunc) (*Response, error)
	}

	GraphQL struct {
		client *Client
	}

	GraphQLQuery struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables,omitempty"`
	}

	GenericGraphQLErrors struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	GraphQLResponseError struct {
		Err    error
		Errors GenericGraphQLErrors
	}
)

var _ GraphQLInterface = (*GraphQL)(nil)

func (e *GraphQLResponseError) Error() string {
	if len(e.Errors.Errors) == 0 {
		return fmt.Sprintf("%s (no additional error messages)", e.Err)
	}

	var sb strings.Builder
	sb.WriteString(e.Err.Error())
	sb.WriteString(" (GraphQL errors: ")

	for i, err := range e.Errors.Errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Message)
	}
	sb.WriteString(")")

	return sb.String()
}

// Do sends a GraphQL query and returns the response in the given response argument
// The response must be JSON serializable. The *Response return value is the HTTP response
// and must be used to retrieve additional HTTP information, like status codes and also
// error messages from failed queries.
//
// Example:
//
//	var response struct {
//		Data struct {
//			Project struct {
//				ID string `json:"id"`
//			} `json:"project"`
//		} `json:"data"`
//	}
//	_, err := client.GraphQL.Do(GraphQLQuery{Query: `query { project(fullPath: "gitlab-org/gitlab") { id } }`}, &response, gitlab.WithContext(ctx))
//
// Attention: This API is experimental and may be subject to breaking changes to improve the API in the future.
func (g *GraphQL) Do(query GraphQLQuery, response any, options ...RequestOptionFunc) (*Response, error) {
	request, err := g.client.NewRequest(http.MethodPost, "", query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL request: %w", err)
	}
	// Overwrite the path of the existing request, as otherwise client-go appends /api/v4 instead.
	request.URL.Path = GraphQLAPIEndpoint
	resp, err := g.client.Do(request, response)
	if err != nil {
		// return error, details can be read from Response
		if errResp, ok := err.(*ErrorResponse); ok { //nolint:errorlint
			var v GenericGraphQLErrors
			if json.Unmarshal(errResp.Body, &v) == nil {
				return resp, &GraphQLResponseError{
					Err:    err,
					Errors: v,
				}
			}
		}
		return resp, fmt.Errorf("failed to execute GraphQL query: %w", err)
	}
	return resp, nil
}

// gidGQL is a global ID. It is used by GraphQL to uniquely identify resources.
type gidGQL struct {
	Type  string
	Int64 int64
}

var gidGQLRegex = regexp.MustCompile(`^gid://gitlab/([^/]+)/(\d+)$`)

func (id *gidGQL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	m := gidGQLRegex.FindStringSubmatch(s)
	if len(m) != 3 {
		return fmt.Errorf("invalid global ID format: %q", s)
	}

	i, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		return fmt.Errorf("failed parsing %q as numeric ID: %w", s, err)
	}

	id.Type = m[1]
	id.Int64 = i

	return nil
}

func (id gidGQL) String() string {
	return fmt.Sprintf("gid://gitlab/%s/%d", id.Type, id.Int64)
}

// iidGQL represents an int64 ID that is encoded by GraphQL as a string.
// This type is used unmarshal the string response into an int64 type.
type iidGQL int64

func (id *iidGQL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("failed parsing %q as numeric ID: %w", s, err)
	}

	*id = iidGQL(i)
	return nil
}

// PageInfo contains cursor-based pagination metadata for GraphQL connections following the Relay
// cursor pagination specification. Use EndCursor and HasNextPage for forward pagination
// (most common), or StartCursor and HasPreviousPage for backward pagination.
//
// Cursors are opaque strings that should not be parsed or constructed manually - always
// use the cursors returned by the API.
//
// Note: GraphQL cursor pagination differs from GitLab's REST API keyset pagination.
// In REST, the pagination link points to the first item of the next page. In GraphQL,
// EndCursor points to the last item of the current page - you pass this to the "after"
// parameter to fetch items after it (essentially an off-by-one difference in semantics).
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#pageinfo
type PageInfo struct {
	EndCursor       string `json:"endCursor"`       // Cursor of the last item in this page (pass to "after" for next page)
	HasNextPage     bool   `json:"hasNextPage"`     // True if more items exist after this page
	StartCursor     string `json:"startCursor"`     // Cursor of the first item in this page (pass to "before" for previous page)
	HasPreviousPage bool   `json:"hasPreviousPage"` // True if items exist before this page
}

// connectionGQL represents a paginated GraphQL connection response following the Relay
// cursor pagination specification. It wraps a list of nodes of any type T along with
// pagination metadata. This type is used internally to unmarshal GraphQL responses from
// GitLab's API, which consistently uses this connection pattern for all paginated fields.
//
// The PageInfo field provides cursors and flags for iterating through pages, while Nodes
// contains the actual data items for the current page.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#connection-fields
type connectionGQL[T any] struct {
	PageInfo PageInfo `json:"pageInfo"`
	Nodes    []T      `json:"nodes"`
}
