// Package equal provides a thin wrapper around github.com/google/go-cmp that uses generics.
// Generics aren't implemented yet. If go-cmp releases a v2 that uses generics,
// this package can be deleted.
//
// https://github.com/google/go-cmp/issues/338
package equal

import "github.com/google/go-cmp/cmp"

// Equal indicates whether the two given variables have the same value.
// It uses the same semantics as go-cmp.Equal,
// except that it enforces that both variables are of the same type.
func Equal[T any](a, b T, opts ...cmp.Option) bool {
	return cmp.Equal(a, b, opts...)
}
