// Package minitrue - Whatever the Package holds to be the truth, is truth.
package minitrue

func Cond[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}
