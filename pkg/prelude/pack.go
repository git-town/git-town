package prelude

type Result[A any] struct {
	value A
	err   error
}

// packs the result of a fallible operation with one return value into a single variable
// so that both values can be provided as a single parameter to another function.
func P[A any](value A, err error) Result[A] {
	return Result[A]{
		value: value,
		err:   err,
	}
}

// checks the given result instance (created using P) using the given ErrorCollector instance
func C[A any](c ErrorCollector, value Result[A]) A {
	c.Check(value.err)
	return value.value
}
