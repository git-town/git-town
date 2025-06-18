//go:build go1.23
// +build go1.23

package gitlab

import (
	"fmt"
	"iter"
)

type PaginationOptionFunc = RequestOptionFunc

// Scan scans all pages for the given request function f and returns individual items in an iterator.
// If an error happens during pagination, the iterator stops immediately.
// The caller must consume the returned error function to retrieve potential errors.
//
//	opts := &ListProjectsOptions{}
//	it, hasErr := Scan(func(p PaginationOptionFunc) ([]*Project, *Response, error) {
//		return c.Projects.ListProjects(opts, p)
//	})
//	projects := slices.Collect(it)
//	if err := hasErr(); err != nil {
//		return err
//	}
//
// or with keyset-based pagination:
//
//	opts := &ListProjectsOptions{
//		ListOptions: ListOptions{
//			OrderBy:    "id",
//			Pagination: "keyset",
//		},
//	}
//	it, hasErr := Scan(func(p PaginationOptionFunc) ([]*Project, *Response, error) {
//		return c.Projects.ListProjects(opts, p)
//	})
//	projects := slices.Collect(it)
//	if err := hasErr(); err != nil {
//		return err
//	}
//
// Attention: This API is experimental and may be subject to breaking changes to improve the API in the future.
func Scan[T any](f func(p PaginationOptionFunc) ([]T, *Response, error)) (iter.Seq[T], func() error) {
	exhausted := false
	var e error
	it := func(yield func(T) bool) {
		defer func() {
			exhausted = true
		}()
		for t, err := range Scan2(f) {
			if err != nil {
				e = err
				return
			}

			if !yield(t) {
				return
			}
		}
	}
	hasErr := func() error {
		if !exhausted {
			panic("called error function of Scan iterator before iterator was exhausted")
		}
		return e
	}
	return it, hasErr
}

// Scan2 scans all pages for the given request function f and returns individual items and potential errors in an iterator.
// The caller must consume the error element of the iterator during each iteration
// to ensure that no errors happened.
//
//	opts := &ListProjectsOptions{}
//	for p, err := range Scan2(func(p PaginationOptionFunc) ([]*Project, *Response, error) {
//		return c.Projects.ListProjects(opts, p)
//	}) {
//		if err != nil {
//			return err
//		}
//		// do something with p
//	}
//
// or with keyset-based pagination:
//
//	opts := &ListProjectsOptions{
//		ListOptions: ListOptions{
//			OrderBy:    "id",
//			Pagination: "keyset",
//		},
//	}
//	for p, err := range Scan2(func(p PaginationOptionFunc) ([]*Project, *Response, error) {
//		return c.Projects.ListProjects(opts, p)
//	}) {
//		if err != nil {
//			return err
//		}
//		// do something with p
//	}
//
// Attention: This API is experimental and may be subject to breaking changes to improve the API in the future.
func Scan2[T any](f func(p PaginationOptionFunc) ([]T, *Response, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var nextOpt PaginationOptionFunc

	Pagination:
		for {
			ts, resp, err := f(nextOpt)
			if err != nil {
				var t T
				yield(t, err)
				return
			}

			for _, t := range ts {
				if !yield(t, nil) {
					return
				}
			}

			// the f request function was either configured for offset- or keyset-based
			// pagination. We support both here, by checking if the next link is provided (keyset)
			// or not. If both are provided, keyset-based pagination takes precedence.
			switch {
			case resp.NextLink != "":
				nextOpt = WithKeysetPaginationParameters(resp.NextLink)
			case resp.NextPage != 0:
				nextOpt = WithOffsetPaginationParameters(resp.NextPage)
			default:
				// no more pages
				break Pagination
			}
		}
	}
}

// Must provides a single item iterator for the provided two item iterator and panics if an error happens.
//
//	opts := &ListProjectsOptions{}
//	for p := range Must(Scan2(func(p PaginationOptionFunc) ([]*Project, *Response, error) {
//		return c.Projects.ListProjects(opts, p)
//	})) {
//		// do something with p
//	}
//
// Attention: This API is experimental and may be subject to breaking changes to improve the API in the future.
func Must[T any](it iter.Seq2[T, error]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x, err := range it {
			if err != nil {
				panic(fmt.Errorf("iterator produced an error: %w", err))
			}

			if !yield(x) {
				return
			}
		}
	}
}
