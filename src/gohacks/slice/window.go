package slice

// Window provides the largest window in a collection
// with the given number of elements around the given cursor position.
func Window(args WindowArgs) WindowResult {
	if args.ElementCount == 0 {
		return WindowResult{0, 0}
	}
	length := args.WindowSize
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
		EndRow:   end,
		StartRow: start,
	}
}

type WindowArgs struct {
	CursorPos    int
	ElementCount int
	WindowSize   int
}

type WindowResult struct {
	EndRow   int
	StartRow int
}
