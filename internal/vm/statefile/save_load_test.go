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
				&opcodes.BranchCreate{Branch: gitdomain.NewLocalBranchName("branch"), StartingPoint: gitdomain.NewSHA("123456").Location()},
				&opcodes.BranchCreateAndCheckoutExistingParent{Ancestors: gitdomain.NewLocalBranchNames("one", "two", "three"), Branch: "branch"},
				&opcodes.BranchCurrentReset{Base: "branch"},
				&opcodes.BranchCurrentResetToParent{CurrentBranch: "branch"},
				&opcodes.BranchDeleteIfEmptyAtRuntime{Branch: "branch"},
				&opcodes.BranchEnsureShippableChanges{Branch: gitdomain.NewLocalBranchName("branch"), Parent: gitdomain.NewLocalBranchName("parent")},
				&opcodes.BranchLocalDelete{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchParentDelete{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchRemoteCreate{Branch: gitdomain.NewLocalBranchName("branch"), SHA: gitdomain.NewSHA("123456")},
				&opcodes.BranchRemoteSetToSHA{Branch: "branch", SetToSHA: "222222"},
				&opcodes.BranchRemoteSetToSHAIfNeeded{Branch: "branch", MustHaveSHA: "111111", SetToSHA: "222222"},
				&opcodes.BranchReset{Target: "branch"},
				&opcodes.BranchTrackingCreate{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchTrackingDelete{Branch: gitdomain.NewRemoteBranchName("origin/branch")},
				&opcodes.BranchesContributionAdd{Branch: gitdomain.NewLocalBranchName("branch")}, // TODO: use string constants here, they get converted to the right data type
				&opcodes.BranchesObservedAdd{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesParkedAdd{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesPerennialAdd{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesPrototypeAdd{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.ChangesDiscard{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.CheckoutUncached{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.Commit{AuthorOverride: Some(gitdomain.Author("user@acme.com")), FallbackToDefaultCommitMessage: true, Message: Some(gitdomain.CommitMessage("my message"))},
				&opcodes.CommitWithMessage{AuthorOverride: Some(gitdomain.Author("user@acme.com")), Message: gitdomain.CommitMessage("my message")},
				&opcodes.ConnectorProposalMerge{Branch: gitdomain.NewLocalBranchName("branch"), CommitMessage: Some(gitdomain.CommitMessage("commit message")), ProposalMessage: "proposal message", ProposalNumber: 123},
				&opcodes.FetchUpstream{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.LineageBranchRemove{Branch: "branch"},
				&opcodes.LineageParentRemove{Branch: "branch"},
				&opcodes.LineageParentSet{Branch: gitdomain.NewLocalBranchName("branch"), Parent: gitdomain.NewLocalBranchName("parent")},
				&opcodes.LineageParentSetFirstExisting{Ancestors: gitdomain.NewLocalBranchNames("one", "two"), Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.LineageParentSetIfExists{Branch: gitdomain.NewLocalBranchName("branch"), Parent: gitdomain.NewLocalBranchName("parent")},
				&opcodes.LineageParentSetToGrandParent{Branch: "branch"},
				&opcodes.Merge{Branch: gitdomain.NewBranchName("branch")},
				&opcodes.MergeAbort{},
				&opcodes.MergeContinue{},
				&opcodes.MergeParent{Parent: gitdomain.NewBranchName("parent")},
				&opcodes.MergeParentIfNeeded{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.ProgramEndOfBranch{},
				&opcodes.ProposalCreate{Branch: gitdomain.NewLocalBranchName("branch"), MainBranch: gitdomain.NewLocalBranchName("main")},
				&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true},
				&opcodes.RebaseAbort{},
				&opcodes.RebaseContinue{},
				&opcodes.RebaseContinueIfNeeded{},
				&opcodes.StashDrop{},
				&opcodes.StashPop{},
				&opcodes.CheckoutHistoryPreserve{PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{Some(gitdomain.NewLocalBranchName("previous"))}},
				&opcodes.PullCurrentBranch{},
				&opcodes.PushCurrentBranch{CurrentBranch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.PushCurrentBranchIfLocal{CurrentBranch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.PushTags{},
				&opcodes.MessageQueue{Message: "message"},
				&opcodes.RebaseBranch{Branch: gitdomain.NewBranchName("branch")},
				&opcodes.RebaseParentIfNeeded{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.RebaseTrackingBranch{RemoteBranch: gitdomain.NewRemoteBranchName("origin/branch")},
				&opcodes.BranchesContributionRemove{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesObservedRemove{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesParkedRemove{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesPerennialRemove{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.BranchesPrototypeRemove{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.ConfigGlobalRemove{Key: configdomain.KeyOffline},
				&opcodes.ConfigLocalRemove{Key: configdomain.KeyOffline},
				&opcodes.BranchLocalRename{NewName: "new", OldName: "old"},
				&opcodes.BranchCurrentResetToSHA{Hard: true, SetToSHA: gitdomain.NewSHA("111111")},
				&opcodes.BranchCurrentResetToSHAIfNeeded{Hard: true, MustHaveSHA: gitdomain.NewSHA("222222"), SetToSHA: gitdomain.NewSHA("111111")},
				&opcodes.StashPopIfNeeded{},
				&opcodes.CommitRevert{SHA: gitdomain.NewSHA("123456")},
				&opcodes.CommitRevertIfNeeded{SHA: gitdomain.NewSHA("123456")},
				&opcodes.ConfigGlobalSet{Key: configdomain.KeyOffline, Value: "1"},
				&opcodes.ConfigLocalSet{Key: configdomain.KeyOffline, Value: "1"},
				&opcodes.MergeSquash{Branch: gitdomain.NewLocalBranchName("branch"), CommitMessage: Some(gitdomain.CommitMessage("commit message")), Parent: gitdomain.NewLocalBranchName("parent")},
				&opcodes.ChangesStage{},
				&opcodes.SnapshotInitialUpdateLocalSHA{Branch: "branch", SHA: "111111"},
				&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: "branch"},
				&opcodes.StashOpenChanges{},
				&opcodes.ProposalUpdateBase{ProposalNumber: 123, NewTarget: gitdomain.NewLocalBranchName("new-target"), OldTarget: gitdomain.NewLocalBranchName("old-target")},
				&opcodes.ProposalUpdateBaseToParent{Branch: "branch", ProposalNumber: 123, OldTarget: gitdomain.NewLocalBranchName("old-target")},
				&opcodes.ProposalUpdateHead{ProposalNumber: 123, NewTarget: gitdomain.NewLocalBranchName("new-target"), OldTarget: gitdomain.NewLocalBranchName("old-target")},
			},
			TouchedBranches: []gitdomain.BranchName{"branch-1", "branch-2"},
			UnfinishedDetails: SomeP(&runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: gitdomain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			}),
			UndoablePerennialCommits: []gitdomain.SHA(nil),
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
      "type": "BranchesObservedAdd"
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
      "type": "BranchesPerennialAdd"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchesPrototypeAdd"
    },
    {
      "data": {},
      "type": "ChangesDiscard"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Checkout"
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
        "Message": "my message"
      },
      "type": "CommitWithMessage"
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
        "Parent": "parent"
      },
      "type": "MergeParent"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "MergeParentIfNeeded"
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
        "ForceIfIncludes": true
      },
      "type": "PushCurrentBranchForceIfNeeded"
    },
    {
      "data": {},
      "type": "RebaseAbort"
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
      "data": {},
      "type": "StashDrop"
    },
    {
      "data": {},
      "type": "StashPop"
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
      "data": {
        "Message": "message"
      },
      "type": "MessageQueue"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RebaseBranch"
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
        "Branch": "branch"
      },
      "type": "BranchesContributionRemove"
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
      "type": "BranchesParkedRemove"
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
      "type": "BranchesPrototypeRemove"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "ConfigGlobalRemove"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "ConfigLocalRemove"
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
      "data": {},
      "type": "StashPopIfNeeded"
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
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "ConfigGlobalSet"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "ConfigLocalSet"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "MergeSquash"
    },
    {
      "data": {},
      "type": "ChangesStage"
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
      "type": "StashOpenChanges"
    },
    {
      "data": {
        "NewTarget": "new-target",
        "OldTarget": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateBase"
    },
    {
      "data": {
        "Branch": "branch",
        "OldTarget": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateBaseToParent"
    },
    {
      "data": {
        "NewTarget": "new-target",
        "OldTarget": "old-target",
        "ProposalNumber": 123
      },
      "type": "ProposalUpdateHead"
    }
  ],
  "TouchedBranches": [
    "branch-1",
    "branch-2"
  ],
  "UndoAPIProgram": [],
  "UndoablePerennialCommits": null,
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
