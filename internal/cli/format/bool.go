package format

import (
	"github.com/git-town/git-town/v21/internal/gohacks"
)

// Bool converts the given bool into either "yes" or "no".
func Bool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func BoolOpt(value gohacks.BoolLike, exists bool) string {
	if exists {
		return Bool(value.IsTrue())
	}
	return "(not provided)"
}
