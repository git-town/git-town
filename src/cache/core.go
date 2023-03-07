// Package cache provides infrastructure to cache things in memory.
package cache

// Bool is a cache for bool variables.
type Bool = Cache[bool]

// String is a cache for string variables.
type String = Cache[string]

// Strings is a cache for string variables.
type Strings = Cache[[]string]
