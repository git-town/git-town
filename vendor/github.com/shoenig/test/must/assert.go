package must

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the must package.
type T interface {
	Helper()
	Fatalf(string, ...any)
}

func errorf(t T, msg string, args ...any) {
	t.Fatalf(msg, args...)
}
