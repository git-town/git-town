package gitlab

import (
	"context"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// contextKey is key of context used internal client-go
type contextKey struct{}

// checkRetryKey is context key of requestRetry.
// Value type of this key must be `retryablehttp.CheckRetry`
// This is used in [WithRequestRetry].
var checkRetryKey = &contextKey{}

// checkRetryFromContext returns checkRetry from Context.
// If checkRetry doesn't exist in context, return nil
func checkRetryFromContext(ctx context.Context) retryablehttp.CheckRetry {
	val := ctx.Value(checkRetryKey)

	// There is no checkRetry in context
	if val == nil {
		return nil
	}

	return val.(retryablehttp.CheckRetry)
}

// contextWithCheckRetry create and return new context with checkRetry
func contextWithCheckRetry(ctx context.Context, checkRetry retryablehttp.CheckRetry) context.Context {
	return context.WithValue(ctx, checkRetryKey, checkRetry)
}
