package commands_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestCreateChildFeatureBranch(t *testing.T) {
	t.Parallel()
	runtime := runtime.CreateGitTown(t)
	err := runtime.CreateFeatureBranch("f1")
	assert.NoError(t, err)
	err = commands.CreateChildFeatureBranch(&runtime, "f1a", "f1")
	assert.NoError(t, err)
	output, err := runtime.BackendRunner.Run("git-town", "config")
	assert.NoError(t, err)
	has := strings.Contains(output, "Branch Ancestry:\n  main\n    f1\n      f1a")
	if !has {
		fmt.Printf("unexpected output: %s", output)
	}
	assert.True(t, has)
}
