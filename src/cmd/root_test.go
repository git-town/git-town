package cmd_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/src/cmd"
	"github.com/stretchr/testify/assert"
)

func TestIsAcceptableVersion(t *testing.T) {
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
	for _, test := range tests {
		have := cmd.IsAcceptableVersion(test.major, test.minor)
		assert.Equal(t, test.want, have, fmt.Sprintf("%d.%d --> %t", test.major, test.minor, test.want))
	}
}
