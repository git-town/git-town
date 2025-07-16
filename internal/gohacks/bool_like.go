package gohacks

// BoolLike describes newtypes that wrap a bool value
type BoolLike interface {
	IsTrue() bool
}
