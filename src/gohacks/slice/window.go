package slice

import "github.com/git-town/git-town/v12/src/gohacks/math"

// Window provides size elements surrounding the given cursor.
func Window[T any](args WindowArgs[T]) (window []T, cursorRow int) {
	if len(args.Elements) == 0 {
		return args.Elements, 0
	}
	length := math.Min(args.Size, len(args.Elements))
	start := args.Cursor - (args.Size / 2)
	var end int
loop:
	for {
		end = start + length
		startsBeforeZero := start < 0
		startsAtZero := start == 0
		endsAtEnd := end == len(args.Elements)
		endsAfterEnd := end > len(args.Elements)
		switch {
		case startsBeforeZero && endsAfterEnd:
			start += 1
			length -= 2
		case startsBeforeZero && endsAtEnd:
			start += 1
			length -= 1
		case startsBeforeZero:
			start += 1
		case endsAfterEnd && startsAtZero:
			length -= 1
		case endsAfterEnd:
			start -= 1
		default:
			break loop
		}
	}
	return args.Elements[start:end], args.Cursor
}

type WindowArgs[T any] struct {
	Elements []T
	Cursor   int
	Size     int
}
