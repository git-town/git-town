// Package option helps use pointers similar to Rust option types.
//
// One of the uses of nil in Go is to express optionality:
// A function that might or might not return a type typically returns a pointer to the type.
// If the pointer is nil, it means there is no value.
// This package helps use such poor-man option types.
package option
