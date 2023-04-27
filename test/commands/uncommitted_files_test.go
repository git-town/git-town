package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestUncommittedFiles(t *testing.T) {
	t.Parallel()
	runtime := runtime.Create(t)
	err := commands.CreateFile(runtime.Dir(), "f1.txt", "one")
	assert.NoError(t, err)
	err = commands.CreateFile(runtime.Dir(), "f2.txt", "two")
	assert.NoError(t, err)
	files, err := commands.UncommittedFiles(&runtime)
	assert.NoError(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
}
