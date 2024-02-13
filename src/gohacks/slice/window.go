package slice

// Window provides size elements surrounding the given cursor.
func Window[T any](args WindowArgs[T]) (window []T, cursorRow int) {
	return args.Elements, args.Cursor
}

type WindowArgs[T any] struct {
	Elements []T
	Cursor   int
	Size     int
}
