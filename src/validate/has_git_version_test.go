package validate_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/validate"
	"github.com/shoenig/test/must"
)

func TestIsAcceptableGitVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		major int
		minor int
		want  bool
	}{
		{2, 7, true},
		{3, 0, true},
		{2, 6, false},
		{1, 8, false},
	}
	for _, tt := range tests {
		have := validate.IsAcceptableGitVersion(tt.major, tt.minor)
		must.EqOp(t, tt.want, have)
	}
}
