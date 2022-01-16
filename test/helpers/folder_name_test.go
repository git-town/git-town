package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v7/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestFolderName(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"foo":                                 "foo",
		`globally set to "true", local unset`: "globally_set_to_true_local_unset",
	}
	for give, want := range tests { //nolint:paralleltest
		t.Run(want, func(t *testing.T) {
			t.Parallel()
			have := helpers.FolderName(give)
			assert.Equal(t, want, have)
		})
	}
}
