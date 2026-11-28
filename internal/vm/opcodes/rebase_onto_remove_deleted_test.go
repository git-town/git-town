package opcodes_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v24/internal/config"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v24/internal/test/commands"
	"github.com/git-town/git-town/v24/internal/test/testruntime"
	"github.com/git-town/git-town/v24/internal/vm/opcodes"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestRebaseOntoRemoveDeleted(t *testing.T) {
	t.Parallel()

	t.Run("does not resolve conflicts from an already active rebase", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file", "target content")
		runtime.MustRun("git", "add", "file")
		runtime.MustRun("git", "commit", "-m", "target commit")
		runtime.MustRun("git", "branch", "target")
		runtime.MustRun("git", "checkout", "-b", "removed")
		runtime.CreateFile("file", "removed content")
		runtime.MustRun("git", "add", "file")
		runtime.MustRun("git", "commit", "-m", "removed commit")
		runtime.MustRun("git", "branch", "oldbase")
		runtime.MustRun("git", "checkout", "-b", "oldtopic", "oldbase")
		runtime.CreateFile("file", "incoming content")
		runtime.MustRun("git", "add", "file")
		runtime.MustRun("git", "commit", "-m", "incoming commit")
		runtime.MustRun("git", "checkout", "-b", "oldonto", "oldbase")
		runtime.CreateFile("file", "current content")
		runtime.MustRun("git", "add", "file")
		runtime.MustRun("git", "commit", "-m", "current commit")
		runtime.MustRun("git", "checkout", "oldtopic")
		rebaseStartErr := runtime.TestRunner.Run("git", "rebase", "--onto", "oldonto", "oldbase")
		must.True(t, rebaseStartErr != nil)
		expectedErr := runtime.TestRunner.Run("git", "-c", "rebase.updateRefs=false", "rebase", "--onto", "target", "removed")
		must.True(t, expectedErr != nil)
		opcode := opcodes.RebaseOntoRemoveDeleted{BranchToRebaseOnto: "target", CommitsToRemove: "removed"}
		prepended, err := runRebaseOntoRemoveDeleted(runtime, opcode)
		must.True(t, err != nil)
		must.EqOp(t, expectedErr.Error(), err.Error())
		must.False(t, prepended)
		must.EqOp(t, "removed content", runtime.MustQuery("git", "show", ":1:file"))
		must.EqOp(t, "current content", runtime.MustQuery("git", "show", ":2:file"))
		must.EqOp(t, "incoming content", runtime.MustQuery("git", "show", ":3:file"))
	})

	t.Run("preserves the original error when the rebase cannot start", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.MustRun("git", "checkout", "-b", "feature-1")
		expectedErr := runtime.TestRunner.Run("git", "-c", "rebase.updateRefs=false", "rebase", "--onto", "missing", "feature-1")
		must.True(t, expectedErr != nil)
		opcode := opcodes.RebaseOntoRemoveDeleted{BranchToRebaseOnto: "missing", CommitsToRemove: "feature-1"}
		prepended, err := runRebaseOntoRemoveDeleted(runtime, opcode)
		must.True(t, err != nil)
		must.EqOp(t, expectedErr.Error(), err.Error())
		must.False(t, prepended)
	})
}

func runRebaseOntoRemoveDeleted(runtime commands.TestCommands, opcode opcodes.RebaseOntoRemoveDeleted) (bool, error) {
	prepended := false
	validatedConfig := config.ValidatedConfig{
		NormalConfig: config.DefaultNormalConfig(),
		ValidatedConfigData: configdomain.ValidatedConfigData{
			MainBranch: "main",
		},
	}
	err := opcode.Run(shared.RunArgs{
		Backend:           runtime.TestRunner,
		BranchInfos:       gitdomain.BranchInfos{},
		Config:            NewMutable(&validatedConfig),
		Connector:         None[forgedomain.Connector](),
		DetectedForgeType: None[forgedomain.DetectedForgeType](),
		FinalMessages:     stringslice.NewCollector(),
		Frontend:          runtime.TestRunner,
		Git:               *runtime.Git,
		Inputs:            dialogcomponents.NewInputs(),
		PrependOpcodes: func(...shared.Opcode) {
			prepended = true
		},
		RegisterUndoablePerennialCommit: func(gitdomain.SHA) {},
		RemoveOpcodes:                   func(func(shared.Opcode) bool) {},
		UpdateInitialSnapshotLocalSHA: func(gitdomain.LocalBranchName, gitdomain.SHA) error {
			return nil
		},
	})
	return prepended, err
}
