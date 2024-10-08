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
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
				&opcodes.AddToContributionBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AddToObservedBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AddToParkedBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AddToPerennialBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AddToPrototypeBranches{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.ChangeParent{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.CheckoutUncached{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.Commit{
					AuthorOverride:                 Some(gitdomain.Author("user@acme.com")),
					FallbackToDefaultCommitMessage: true,
					Message:                        Some(gitdomain.CommitMessage("my message")),
				},
				&opcodes.CommitWithMessage{
					AuthorOverride: Some(gitdomain.Author("user@acme.com")),
					Message:        gitdomain.CommitMessage("my message"),
				},
				&opcodes.ConnectorMergeProposal{
					Branch:          gitdomain.NewLocalBranchName("branch"),
					CommitMessage:   Some(gitdomain.CommitMessage("commit message")),
					ProposalMessage: "proposal message",
					ProposalNumber:  123,
				},
				&opcodes.ContinueMerge{},
				&opcodes.ContinueRebase{},
				&opcodes.ContinueRebaseIfNeeded{},
				&opcodes.CreateBranch{
					Branch:        gitdomain.NewLocalBranchName("branch"),
					StartingPoint: gitdomain.NewSHA("123456").Location(),
				},
				&opcodes.CreateProposal{
					Branch:     gitdomain.NewLocalBranchName("branch"),
					MainBranch: gitdomain.NewLocalBranchName("main"),
				},
				&opcodes.CreateRemoteBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
					SHA:    gitdomain.NewSHA("123456"),
				},
				&opcodes.CreateTrackingBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.DeleteParentBranch{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.DeleteTrackingBranch{
					Branch: gitdomain.NewRemoteBranchName("origin/branch"),
				},
				&opcodes.DiscardOpenChanges{},
				&opcodes.DropStash{},
				&opcodes.EndOfBranchProgram{},
				&opcodes.EnsureHasShippableChanges{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcodes.FetchUpstream{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.ForcePushCurrentBranch{ForceIfIncludes: true},
				&opcodes.Merge{Branch: gitdomain.NewBranchName("branch")},
				&opcodes.MergeParent{
					Parent: gitdomain.NewBranchName("parent"),
				},
				&opcodes.MergeParentIfNeeded{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.PreserveCheckoutHistory{
					PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{Some(gitdomain.NewLocalBranchName("previous"))},
				},
				&opcodes.PullCurrentBranch{},
				&opcodes.PushCurrentBranch{
					CurrentBranch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.PushCurrentBranchIfNeeded{
					CurrentBranch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.PushTags{},
				&opcodes.RebaseBranch{
					Branch: gitdomain.NewBranchName("branch"),
				},
				&opcodes.RebaseParentIfNeeded{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RebaseFeatureTrackingBranch{
					RemoteBranch: gitdomain.NewRemoteBranchName("origin/branch"),
				},
				&opcodes.RemoveFromContributionBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RemoveFromObservedBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RemoveFromParkedBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RemoveFromPerennialBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RemoveFromPrototypeBranches{
					Branch: gitdomain.NewLocalBranchName("branch"),
				},
				&opcodes.RemoveGlobalConfig{
					Key: configdomain.KeyOffline,
				},
				&opcodes.RemoveLocalConfig{
					Key: configdomain.KeyOffline,
				},
				&opcodes.RenameBranch{
					NewName: "new",
					OldName: "old",
				},
				&opcodes.ResetCurrentBranchToSHA{
					Hard:     true,
					SetToSHA: gitdomain.NewSHA("111111"),
				},
				&opcodes.ResetCurrentBranchToSHAIfNeeded{
					Hard:        true,
					MustHaveSHA: gitdomain.NewSHA("222222"),
					SetToSHA:    gitdomain.NewSHA("111111"),
				},
				&opcodes.RestoreOpenChanges{},
				&opcodes.RevertCommit{
					SHA: gitdomain.NewSHA("123456"),
				},
				&opcodes.RevertCommitIfNeeded{
					SHA: gitdomain.NewSHA("123456"),
				},
				&opcodes.SetGlobalConfig{
					Key:   configdomain.KeyOffline,
					Value: "1",
				},
				&opcodes.SetLocalConfig{
					Key:   configdomain.KeyOffline,
					Value: "1",
				},
				&opcodes.SetParent{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcodes.SetParentIfBranchExists{
					Branch: gitdomain.NewLocalBranchName("branch"),
					Parent: gitdomain.NewLocalBranchName("parent"),
				},
				&opcodes.SkipCurrentBranch{},
				&opcodes.SquashMerge{
					Branch:        gitdomain.NewLocalBranchName("branch"),
					CommitMessage: Some(gitdomain.CommitMessage("commit message")),
					Parent:        gitdomain.NewLocalBranchName("parent"),
				},
				&opcodes.StageOpenChanges{},
				&opcodes.StashOpenChanges{},
				&opcodes.UpdateProposalBase{
					ProposalNumber: 123,
					NewTarget:      gitdomain.NewLocalBranchName("new-target"),
					OldTarget:      gitdomain.NewLocalBranchName("old-target"),
				},
				&opcodes.UpdateProposalHead{
					ProposalNumber: 123,
					NewTarget:      gitdomain.NewLocalBranchName("new-target"),
					OldTarget:      gitdomain.NewLocalBranchName("old-target"),
				},
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
      "data": {},
      "type": "AbortMerge"
    },
    {
      "data": {},
      "type": "AbortRebase"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToContributionBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToObservedBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToParkedBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToPerennialBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "AddToPrototypeBranches"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "ChangeParent"
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
      "type": "ConnectorMergeProposal"
    },
    {
      "data": {},
      "type": "ContinueMerge"
    },
    {
      "data": {},
      "type": "ContinueRebase"
    },
    {
      "data": {},
      "type": "ContinueRebaseIfNeeded"
    },
    {
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "CreateBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "MainBranch": "main",
        "ProposalBody": "",
        "ProposalTitle": ""
      },
      "type": "CreateProposal"
    },
    {
      "data": {
        "Branch": "branch",
        "SHA": "123456"
      },
      "type": "CreateRemoteBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CreateTrackingBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "DeleteLocalBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "DeleteParentBranch"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "DeleteTrackingBranch"
    },
    {
      "data": {},
      "type": "DiscardOpenChanges"
    },
    {
      "data": {},
      "type": "DropStash"
    },
    {
      "data": {},
      "type": "EndOfBranchProgram"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "EnsureHasShippableChanges"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "FetchUpstream"
    },
    {
      "data": {
        "ForceIfIncludes": true
      },
      "type": "ForcePushCurrentBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "Merge"
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
      "data": {
        "PreviousBranchCandidates": [
          "previous"
        ]
      },
      "type": "PreserveCheckoutHistory"
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
      "type": "PushCurrentBranchIfNeeded"
    },
    {
      "data": {},
      "type": "PushTags"
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
      "type": "RebaseFeatureTrackingBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromContributionBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromObservedBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromParkedBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromPerennialBranches"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromPrototypeBranches"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveGlobalConfig"
    },
    {
      "data": {
        "Key": "git-town.offline"
      },
      "type": "RemoveLocalConfig"
    },
    {
      "data": {
        "NewName": "new",
        "OldName": "old"
      },
      "type": "RenameBranch"
    },
    {
      "data": {
        "Hard": true,
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
    },
    {
      "data": {
        "Hard": true,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHAIfNeeded"
    },
    {
      "data": {},
      "type": "RestoreOpenChanges"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "RevertCommit"
    },
    {
      "data": {
        "SHA": "123456"
      },
      "type": "RevertCommitIfNeeded"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetGlobalConfig"
    },
    {
      "data": {
        "Key": "git-town.offline",
        "Value": "1"
      },
      "type": "SetLocalConfig"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "SetParent"
    },
    {
      "data": {
        "Branch": "branch",
        "Parent": "parent"
      },
      "type": "SetParentIfBranchExists"
    },
    {
      "data": {},
      "type": "SkipCurrentBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Parent": "parent"
      },
      "type": "SquashMerge"
    },
    {
      "data": {},
      "type": "StageOpenChanges"
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
      "type": "UpdateProposalBase"
    },
    {
      "data": {
        "NewTarget": "new-target",
        "OldTarget": "old-target",
        "ProposalNumber": 123
      },
      "type": "UpdateProposalHead"
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
