package validate_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v9/src/validate"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tt.want, have, fmt.Sprintf("%d.%d --> %t", tt.major, tt.minor, tt.want))
	}
}
