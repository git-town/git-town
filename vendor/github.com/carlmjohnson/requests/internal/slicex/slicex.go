package slicex

func Clip[T any](sp *[]T) {
	s := *sp
	*sp = s[:len(s):len(s)]
}
