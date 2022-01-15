package git_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

func TestMainFirst(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give []string
		want []string
	}{
		{give: []string{"main", "one", "two"}, want: []string{"main", "one", "two"}},
		{give: []string{"alpha", "main", "omega"}, want: []string{"main", "alpha", "omega"}},
		{give: []string{"main"}, want: []string{"main"}},
		{give: []string{}, want: []string{}},
	}
	for tt := range tests {
		have := git.MainFirst(tests[tt].give)
		assert.Equal(t, tests[tt].want, have)
	}
}
