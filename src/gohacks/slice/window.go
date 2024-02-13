package slice

import "github.com/git-town/git-town/v12/src/gohacks/math"

// Window tells you which of the given elements to display given where the cursor is.
func Window(args WindowArgs) WindowResult {
	if args.ElementCount == 0 {
		return WindowResult{0, 0, 0}
	}
	length := math.Min(args.WindowSize, args.ElementCount)
	start := args.CursorPos - (args.WindowSize / 2)
	var end int
loop:
	for {
		end = start + length
		startsBeforeZero := start < 0
		startsAtZero := start == 0
		endsAtEnd := end == args.ElementCount
		endsAfterEnd := end > args.ElementCount
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
	return WindowResult{
		StartRow:  start,
		EndRow:    end,
		CursorRow: args.CursorPos,
	}
}

type WindowArgs struct {
	ElementCount int
	CursorPos    int
	WindowSize   int
}

type WindowResult struct {
	StartRow  int
	EndRow    int
	CursorRow int
}
