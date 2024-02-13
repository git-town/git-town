package slice

import "github.com/git-town/git-town/v12/src/gohacks/math"

// Window provides size elements surrounding the given cursor.
func Window[T any](args WindowArgs[T]) (window []T, cursorRow int) {
	if len(args.Elements) == 0 {
		return args.Elements, 0
	}
	end := math.Min(args.Size, len(args.Elements))
	return args.Elements[0:end], args.Cursor
}

type WindowArgs[T any] struct {
	Elements []T
	Cursor   int
	Size     int
}
