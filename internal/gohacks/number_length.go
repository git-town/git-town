package gohacks

func NumberLength(number int) int {
	switch {
	case number < 0:
		return NumberLength(0-number) + 1
	case number < 10:
		return 1
	case number < 100:
		return 2
	default:
		return 3
	}
}
