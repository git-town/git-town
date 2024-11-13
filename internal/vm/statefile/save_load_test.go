package statefile_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	"github.com/git-town/git-town/v16/internal/vm/statefile"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLoadSave(t *testing.T) {
	t.Parallel()

	t.Run("SanitizePath", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"/home/user/development/git-town":        "home-user-development-git-town",
			"c:\\Users\\user\\development\\git-town": "c-users-user-development-git-town",
		}
		for give, want := range tests {
			rootDir := gitdomain.NewRepoRootDir(give)
			have := statefile.SanitizePath(rootDir)
			must.EqOp(t, want, have)
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			AbortProgram:          program.Program{},
			BeginBranchesSnapshot: gitdomain.EmptyBranchesSnapshot(),
			BeginConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
			BeginStashSize:        0,
			Command:               "command",
			DryRun:                true,
			EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
			EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
			EndStashSize:          Some(gitdomain.StashSize(1)),
			RunProgram: program.Program{
				&opcodes.BranchCreate{Branch: "branch", StartingPoint: "123456"},
				&opcodes.BranchCreateAndCheckoutExistingParent{Ancestors: gitdomain.NewLocalBranchNames("one", "two", "three"), Branch: "branch"},
				&opcodes.BranchCurrentReset{Base: "branch"},
				&opcodes.BranchCurrentResetToParent{CurrentBranch: "branch"},
				&opcodes.BranchCurrentResetToSHA{Hard: true, SetToSHA: "111111"},
				&opcodes.BranchCurrentResetToSHAIfNeeded{Hard: true, MustHaveSHA: "222222", SetToSHA: "111111"},
				&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: "branch"},
				&opcodes.BranchEnsureShippableChanges{Branch: "branch", Parent: "parent"},
				&opcodes.BranchLocalDelete{Branch: "branch"},
				&opcodes.BranchLocalRename{NewName: "new", OldName: "old"},
				&opcodes.BranchParentDelete{Branch: "branch"},
				&opcodes.BranchRemoteCreate{Branch: "branch", SHA: "123456"},
				&opcodes.BranchRemoteSetToSHA{Branch: "branch", SetToSHA: "222222"},
				&opcodes.BranchRemoteSetToSHAIfNeeded{Branch: "branch", MustHaveSHA: "111111", SetToSHA: "222222"},
				&opcodes.BranchReset{Target: "branch"},
				&opcodes.BranchTrackingCreate{Branch: "branch"},
				&opcodes.BranchTrackingDelete{Branch: "origin/branch"},
				&opcodes.BranchesContributionAdd{Branch: "branch"},
				&opcodes.BranchesContributionRemove{Branch: "branch"},
				&opcodes.BranchesObservedAdd{Branch: "branch"},
				&opcodes.BranchesObservedRemove{Branch: "branch"},
				&opcodes.BranchesParkedAdd{Branch: "branch"},
				&opcodes.BranchesParkedRemove{Branch: "branch"},
				&opcodes.BranchesPerennialAdd{Branch: "branch"},
				&opcodes.BranchesPerennialRemove{Branch: "branch"},
				&opcodes.BranchesPrototypeAdd{Branch: "branch"},
				&opcodes.BranchesPrototypeRemove{Branch: "branch"},
				&opcodes.ChangesDiscard{},
				&opcodes.ChangesStage{},
				&opcodes.Checkout{Branch: "branch"},
				&opcodes.CheckoutHistoryPreserve{PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{Some(gitdomain.NewLocalBranchName("previous"))}},
				&opcodes.CheckoutIfNeeded{Branch: "branch"},
				&opcodes.CheckoutUncached{Branch: "branch"},
				&opcodes.Commit{AuthorOverride: Some(gitdomain.Author("user@acme.com")), FallbackToDefaultCommitMessage: true, Message: Some(gitdomain.CommitMessage("my message"))},
				&opcodes.CommitAutoUndo{AuthorOverride: Some(gitdomain.Author("user@acme.com")), FallbackToDefaultCommitMessage: true, Message: Some(gitdomain.CommitMessage("my message"))},
				&opcodes.CommitMessageCommentOut{},
				&opcodes.CommitRevert{SHA: "123456"},
				&opcodes.CommitRevertIfNeeded{SHA: "123456"},
				&opcodes.CommitWithMessage{AuthorOverride: Some(gitdomain.Author("user@acme.com")), Message: "my message"},
				&opcodes.ConfigRemove{Key: configdomain.KeyOffline, Scope: configdomain.ConfigScopeLocal},
				&opcodes.ConfigSet{Key: configdomain.KeyOffline, Scope: configdomain.ConfigScopeLocal, Value: "1"},
				&opcodes.ConflictPhantomDetect{ParentBranch: Some(gitdomain.NewLocalBranchName("parent")), ParentSHA: Some(gitdomain.NewSHA("123456"))},
				&opcodes.ConflictPhantomFinalize{},
				&opcodes.ConflictPhantomResolve{FilePath: "file"},
				&opcodes.ConnectorProposalMerge{Branch: "branch", CommitMessage: Some(gitdomain.CommitMessage("commit message")), ProposalMessage: "proposal message", ProposalNumber: 123},
				&opcodes.FetchUpstream{Branch: "branch"},
				&opcodes.LineageBranchRemove{Branch: "branch"},
				&opcodes.LineageParentRemove{Branch: "branch"},
				&opcodes.LineageParentSet{Branch: "branch", Parent: "parent"},
				&opcodes.LineageParentSetFirstExisting{Ancestors: gitdomain.NewLocalBranchNames("one", "two"), Branch: "branch"},
				&opcodes.LineageParentSetIfExists{Branch: "branch", Parent: "parent"},
				&opcodes.LineageParentSetToGrandParent{Branch: "branch"},
				&opcodes.Merge{Branch: "branch"},
				&opcodes.MergeAbort{},
				&opcodes.MergeContinue{},
				&opcodes.MergeParent{CurrentParent: "parent", OriginalParentName: Some(gitdomain.NewLocalBranchName("original-parent")), OriginalParentSHA: Some(gitdomain.NewSHA("123456"))},
				&opcodes.MergeParentIfNeeded{Branch: "branch", OriginalParentName: Some(gitdomain.NewLocalBranchName("original-parent")), OriginalParentSHA: Some(gitdomain.NewSHA("123456"))},
				&opcodes.MergeSquashProgram{Authors: []gitdomain.Author{"author 1 <one@acme.com>", "author 2 <two@acme.com>"}, Branch: "branch", CommitMessage: Some(gitdomain.CommitMessage("commit message")), Parent: "parent"},
				&opcodes.MessageQueue{Message: "message"},
				&opcodes.ProgramEndOfBranch{},
				&opcodes.ProposalCreate{Branch: "branch", MainBranch: "main"},
				&opcodes.ProposalUpdateTarget{ProposalNumber: 123, NewBranch: "new-target", OldBranch: "old-target"},
				&opcodes.ProposalUpdateTargetToGrandParent{Branch: "branch", ProposalNumber: 123, OldTarget: "old-target"},
				&opcodes.ProposalUpdateSource{ProposalNumber: 123, NewBranch: "new-target", OldBranch: "old-target"},
				&opcodes.PullCurrentBranch{},
				&opcodes.PushCurrentBranch{CurrentBranch: "branch"},
				&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true},
				&opcodes.PushCurrentBranchIfLocal{CurrentBranch: "branch"},
				&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: "branch"},
				&opcodes.PushTags{},
				&opcodes.RebaseAbort{},
				&opcodes.RebaseBranch{Branch: "branch"},
				&opcodes.RebaseContinue{},
				&opcodes.RebaseContinueIfNeeded{},
				&opcodes.RebaseParentIfNeeded{Branch: "branch"},
				&opcodes.RebaseTrackingBranch{RemoteBranch: "origin/branch"},
				&opcodes.RegisterUndoablePerennialCommit{Parent: "parent"},
				&opcodes.SnapshotInitialUpdateLocalSHA{Branch: "branch", SHA: "111111"},
				&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: "branch"},
				&opcodes.StashDrop{},
				&opcodes.StashPop{},
				&opcodes.StashPopIfNeeded{},
				&opcodes.StashOpenChanges{},
			},
			TouchedBranches: []gitdomain.BranchName{"branch-1", "branch-2"},
			UnfinishedDetails: SomeP(&runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: gitdomain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			}),
			UndoablePerennialCommits: []gitdomain.SHA{},
		}

		wantJSON := `
{
  "AbortProgram": [],
  "BeginBranchesSnapshot": {
    "Active": null,
    "Branches": []
  },
  "BeginConfigSnapshot": {
    "Global": {},
    "Local": {}
  },
  "BeginStashSize": 0,
  "Command": "command",
  "DryRun": true,
  "EndBranchesSnapshot": null,
  "EndConfigSnapshot": null,
  "EndStashSize": 1,
  "FinalUndoProgram": [],
  "RunProgram": [
    {
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "BranchCreate"
    },
    {
      "data": {
        "Ancestors": [
          "one",
          "two",
          "three"
        ],
        "Branch": "branch"
      },
      "type": "BranchCreateAndCheckoutExistingParent"
    },
    {
      "data": {
        "Base": "branch"
      },
      "type": "BranchCurrentReset"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "BranchCurrentResetToParent"
    },
    {
      "data": {
        "Hard": true,
        "SetToSHA": "111111"
      },
      "type": "BranchCurrentResetToSHA"
    },
    {
      "data": {
        "Hard": true,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "BranchCurrentResetToSHAIfNeeded"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchDeleteIfEmptyAtRuntime"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "BranchEnsureShippableChanges"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchLocalDelete"
    },
    {
      "data": {
        "NewName": "new",
        "OldName": "old"
      },
      "type": "BranchLocalRename"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchParentDelete"
    },
    {
      "data": {
        "Branch": "branch",
        "SHA": "123456"
      },
      "type": "BranchRemoteCreate"
    },
    {
      "data": {
        "Branch": "branch",
        "SetToSHA": "222222"
      },
      "type": "BranchRemoteSetToSHA"
    },
    {
      "data": {
        "Branch": "branch",
        "MustHaveSHA": "111111",
        "SetToSHA": "222222"
      },
      "type": "BranchRemoteSetToSHAIfNeeded"
    },
    {
      "data": {
        "Target": "branch"
      },
      "type": "BranchReset"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchTrackingCreate"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "BranchTrackingDelete"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesContributionAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesContributionRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesObservedAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesObservedRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesParkedAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesParkedRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesPerennialAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesPerennialRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesPrototypeAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesPrototypeRemove"
    },
    {
      "data": {},
      "type": "ChangesDiscard"
    },
    {
      "data": {},
      "type": "ChangesStage"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Checkout"
    },
    {
      "data": {
        "PreviousBranchCandidates": [
          "previous"
        ]
      },
      "type": "CheckoutHistoryPreserve"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CheckoutIfNeeded"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CheckoutUncached"
    },
    {
      "data": {
        "AuthorOverride": "user@acme.com",
        "FallbackToDefaultCommitMessage": true,
        "Message": "my message"
      },
      "type": "Commit"
    },
    {
      "data": {
        "AuthorOverride": "user@acme.com",
        "FallbackToDefaultCommitMessage": true,
        "Message": "my message"
      },
      "type": "CommitAutoUndo"
    },
    {
      "data": {},
      "type": "CommitMessageCommentOut"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "CommitRevert"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "CommitRevertIfNeeded"
    },
    {
      "data": {
        "AuthorOverride": "user@acme.com",
        "Message": "my message"
      },
      "type": "CommitWithMessage"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Scope": "local"
      },
      "type": "ConfigRemove"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Scope": "local",
        "Value": "1"
      },
      "type": "ConfigSet"
    },
    {
      "data": {
        "ParentBranch": "parent",
        "ParentSHA": "123456"
      },
      "type": "ConflictPhantomDetect"
    },
    {
      "data": {},
      "type": "ConflictPhantomFinalize"
    },
    {
      "data": {
        "FilePath": "file"
      },
      "type": "ConflictPhantomResolve"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "ProposalMessage": "proposal message",
        "ProposalNumber": 123
      },
      "type": "ConnectorProposalMerge"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "FetchUpstream"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "LineageBranchRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "LineageParentRemove"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "LineageParentSet"
    },
    {
      "data": {
        "Ancestors": [
          "one",
          "two"
        ],
        "Branch": "branch"
      },
      "type": "LineageParentSetFirstExisting"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "LineageParentSetIfExists"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "LineageParentSetToGrandParent"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Merge"
    },
    {
      "data": {},
      "type": "MergeAbort"
    },
    {
      "data": {},
      "type": "MergeContinue"
    },
    {
      "data": {
        "CurrentParent": "parent",
        "OriginalParentName": "original-parent",
        "OriginalParentSHA": "123456"
      },
      "type": "MergeParent"
    },
    {
      "data": {
        "Branch": "branch",
        "OriginalParentName": "original-parent",
        "OriginalParentSHA": "123456"
      },
      "type": "MergeParentIfNeeded"
    },
    {
      "data": {
        "Authors": [
          "author 1 \u003cone@acme.com\u003e",
          "author 2 \u003ctwo@acme.com\u003e"
        ],
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "MergeSquashProgram"
    },
    {
      "data": {
        "Message": "message"
      },
      "type": "MessageQueue"
    },
    {
      "data": {},
      "type": "ProgramEndOfBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "MainBranch": "main",
        "ProposalBody": "",
        "ProposalTitle": ""
      },
      "type": "ProposalCreate"
    },
    {
      "data": {
        "NewBranch": "new-target",
        "OldBranch": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateTarget"
    },
    {
      "data": {
        "Branch": "branch",
        "OldTarget": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateTargetToParent"
    },
    {
      "data": {
        "NewBranch": "new-target",
        "OldBranch": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateSource"
    },
    {
      "data": {},
      "type": "PullCurrentBranch"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "PushCurrentBranch"
    },
    {
      "data": {
        "ForceIfIncludes": true
      },
      "type": "PushCurrentBranchForceIfNeeded"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "PushCurrentBranchIfLocal"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "PushCurrentBranchIfNeeded"
    },
    {
      "data": {},
      "type": "PushTags"
    },
    {
      "data": {},
      "type": "RebaseAbort"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RebaseBranch"
    },
    {
      "data": {},
      "type": "RebaseContinue"
    },
    {
      "data": {},
      "type": "RebaseContinueIfNeeded"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RebaseParentIfNeeded"
    },
    {
      "data": {
        "PushBranches": false,
        "RemoteBranch": "origin/branch"
      },
      "type": "RebaseTrackingBranch"
    },
    {
      "data": {
        "Parent": "parent"
      },
      "type": "RegisterUndoablePerennialCommit"
    },
    {
      "data": {
        "Branch": "branch",
        "SHA": "111111"
      },
      "type": "SnapshotInitialUpdateLocalSHA"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "SnapshotInitialUpdateLocalSHAIfNeeded"
    },
    {
      "data": {},
      "type": "StashDrop"
    },
    {
      "data": {},
      "type": "StashPop"
    },
    {
      "data": {},
      "type": "StashPopIfNeeded"
    },
    {
      "data": {},
      "type": "StashOpenChanges"
    }
  ],
  "TouchedBranches": [
    "branch-1",
    "branch-2"
  ],
  "UndoAPIProgram": [],
  "UndoablePerennialCommits": [],
  "UnfinishedDetails": {
    "CanSkip": true,
    "EndBranch": "end-branch",
    "EndTime": "0001-01-01T00:00:00Z"
  }
}`[1:]

		repoRoot := gitdomain.NewRepoRootDir("/path/to/git-town-unit-tests")
		err := statefile.Save(runState, repoRoot)
		must.NoError(t, err)
		filepath, err := statefile.FilePath(repoRoot)
		must.NoError(t, err)
		content, err := os.ReadFile(filepath)
		must.NoError(t, err)
		must.EqOp(t, wantJSON, string(content))
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		must.NoError(t, err)
		// NOTE: comparing runState and newState directly leads to incorrect test failures
		// solely due to different pointer addresses, even when using reflect.DeepEqual.
		// Comparing the serialization seems to work better here.
		runStateText, err := json.MarshalIndent(runState, "", "  ")
		must.NoError(t, err)
		newStateText, err := json.MarshalIndent(newState, "", "  ")
		must.NoError(t, err)
		must.EqOp(t, string(runStateText), string(newStateText))
	})
}
