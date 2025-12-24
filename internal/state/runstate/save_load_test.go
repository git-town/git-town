package runstate_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/state"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
			have := state.SanitizePath(rootDir)
			must.EqOp(t, want, have)
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			AbortProgram:          program.Program{},
			BeginBranchesSnapshot: gitdomain.EmptyBranchesSnapshot(),
			BeginConfigSnapshot: configdomain.BeginConfigSnapshot{
				Global:   configdomain.SingleSnapshot{},
				Local:    configdomain.SingleSnapshot{},
				Unscoped: configdomain.SingleSnapshot{},
			},
			BeginStashSize:      0,
			BranchInfosLastRun:  None[gitdomain.BranchInfos](),
			Command:             "command",
			DryRun:              true,
			EndBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
			EndConfigSnapshot:   None[configdomain.EndConfigSnapshot](),
			EndStashSize:        Some(gitdomain.StashSize(1)),
			RunProgram: program.Program{
				&opcodes.BranchCreate{Branch: "branch", StartingPoint: "123456"},
				&opcodes.BranchCreateAndCheckoutExistingParent{Ancestors: gitdomain.NewLocalBranchNames("one", "two", "three"), Branch: "branch"},
				&opcodes.BranchCurrentReset{Base: "branch"},
				&opcodes.BranchCurrentResetToParent{CurrentBranch: "branch"},
				&opcodes.BranchCurrentResetToSHA{SHA: "111111"},
				&opcodes.BranchCurrentResetToSHAIfNeeded{MustHaveSHA: "222222", SetToSHA: "111111"},
				&opcodes.BranchEnsureShippableChanges{Branch: "branch", Parent: "parent"},
				&opcodes.BranchLocalDelete{Branch: "branch"},
				&opcodes.BranchLocalDeleteContent{BranchToDelete: "branch", BranchToRebaseOnto: "main"},
				&opcodes.BranchLocalRename{NewName: "new", OldName: "old"},
				&opcodes.BranchRemoteCreate{Branch: "branch", SHA: "123456"},
				&opcodes.BranchRemoteSetToSHA{Branch: "branch", SetToSHA: "222222"},
				&opcodes.BranchRemoteSetToSHAIfNeeded{Branch: "branch", MustHaveSHA: "111111", SetToSHA: "222222"},
				&opcodes.BranchReset{Target: "branch"},
				&opcodes.BranchTrackingCreate{Branch: "branch"},
				&opcodes.BranchTrackingCreateIfLocalExists{Branch: "branch"},
				&opcodes.BranchTrackingCreateIfNeeded{CurrentBranch: "branch"},
				&opcodes.BranchTrackingDelete{Branch: "origin/branch"},
				&opcodes.BranchTypeOverrideSet{Branch: "branch", BranchType: configdomain.BranchTypeFeatureBranch},
				&opcodes.BranchTypeOverrideRemove{Branch: "branch"},
				&opcodes.BranchWithRemoteGoneDeleteIfEmptyAtRuntime{Branch: "branch"},
				&opcodes.ChangesDiscard{},
				&opcodes.ChangesStage{},
				&opcodes.ChangesUnstageAll{},
				&opcodes.Checkout{Branch: "branch"},
				&opcodes.CheckoutHistoryPreserve{PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{gitdomain.NewLocalBranchNameOption("previous")}},
				&opcodes.CheckoutIfNeeded{Branch: "branch"},
				&opcodes.CheckoutUncached{Branch: "branch"},
				&opcodes.CherryPick{SHA: "123456"},
				&opcodes.CherryPickContinue{},
				&opcodes.Commit{AuthorOverride: Some(gitdomain.Author("user@acme.com")), FallbackToDefaultCommitMessage: true, Message: Some(gitdomain.CommitMessage("my message"))},
				&opcodes.CommitAutoUndo{AuthorOverride: Some(gitdomain.Author("user@acme.com")), FallbackToDefaultCommitMessage: true, Message: Some(gitdomain.CommitMessage("my message"))},
				&opcodes.CommitMessageCommentOut{},
				&opcodes.CommitRemove{SHA: "123456"},
				&opcodes.CommitRevert{SHA: "123456"},
				&opcodes.CommitRevertIfNeeded{SHA: "123456"},
				&opcodes.CommitWithMessage{AuthorOverride: Some(gitdomain.Author("user@acme.com")), Message: "my message", CommitHook: configdomain.CommitHookEnabled},
				&opcodes.ConfigRemove{Key: configdomain.KeyOffline, Scope: configdomain.ConfigScopeLocal},
				&opcodes.ConfigSet{Key: configdomain.KeyOffline, Scope: configdomain.ConfigScopeLocal, Value: "1"},
				&opcodes.ConflictMergePhantomFinalize{},
				&opcodes.ConflictMergePhantomResolveAll{CurrentBranch: "current", ParentBranch: gitdomain.NewLocalBranchNameOption("parent"), ParentSHA: Some(gitdomain.NewSHA("123456"))},
				&opcodes.ConflictResolve{FilePath: "file", Resolution: gitdomain.ConflictResolutionOurs},
				&opcodes.ConnectorProposalMerge{Branch: "branch", CommitMessage: Some(gitdomain.CommitMessage("commit message")), Proposal: forgedomain.Proposal{Data: forgedomain.BitbucketCloudProposalData{ProposalData: forgedomain.ProposalData{Active: true, Body: gitdomain.NewProposalBodyOpt("body"), MergeWithAPI: true, Number: 123, Source: "source", Target: "target", Title: "title", URL: "url"}}, ForgeType: forgedomain.ForgeTypeBitbucket}},
				&opcodes.ExecuteShellCommand{Args: []string{"arg1", "arg2"}, Executable: "executable"},
				&opcodes.ExitToShell{},
				&opcodes.FetchUpstream{Branch: "branch"},
				&opcodes.FileRemove{FilePath: "file"},
				&opcodes.FileStage{FilePath: "file"},
				&opcodes.LineageBranchRemove{Branch: "branch"},
				&opcodes.LineageParentRemove{Branch: "branch"},
				&opcodes.LineageParentSet{Branch: "branch", Parent: "parent"},
				&opcodes.LineageParentSetFirstExisting{Ancestors: gitdomain.NewLocalBranchNames("one", "two"), Branch: "branch"},
				&opcodes.LineageParentSetIfExists{Branch: "branch", Parent: "parent"},
				&opcodes.LineageParentSetToGrandParent{Branch: "branch"},
				&opcodes.MergeIntoCurrentBranch{BranchToMerge: "branch"},
				&opcodes.MergeAbort{},
				&opcodes.MergeAlwaysProgram{Branch: "branch", CommitMessage: Some(gitdomain.CommitMessage("commit message"))},
				&opcodes.MergeContinue{},
				&opcodes.MergeParentResolvePhantomConflicts{CurrentBranch: "current", CurrentParent: "parent", InitialParentName: gitdomain.NewLocalBranchNameOption("original-parent"), InitialParentSHA: Some(gitdomain.NewSHA("123456"))},
				&opcodes.MergeSquashProgram{Authors: []gitdomain.Author{"author 1 <one@acme.com>", "author 2 <two@acme.com>"}, Branch: "branch", CommitMessage: Some(gitdomain.CommitMessage("commit message")), Parent: "parent"},
				&opcodes.MessageQueue{Message: "message"},
				&opcodes.ProgramEndOfBranch{},
				&opcodes.ProposalCreate{Branch: "branch", MainBranch: "main"},
				&opcodes.ProposalUpdateTarget{Proposal: forgedomain.Proposal{Data: forgedomain.ProposalData{Active: true, Body: gitdomain.NewProposalBodyOpt("body"), MergeWithAPI: true, Number: 123, Source: "source", Target: "target", Title: "title", URL: "url"}, ForgeType: forgedomain.ForgeTypeGitLab}, NewBranch: "new-target", OldBranch: "old-target"},
				&opcodes.ProposalUpdateTargetToGrandParent{Branch: "branch", Proposal: forgedomain.Proposal{Data: forgedomain.ProposalData{Active: true, Body: gitdomain.NewProposalBodyOpt("body"), MergeWithAPI: true, Number: 123, Source: "source", Target: "target", Title: "title", URL: "url"}, ForgeType: forgedomain.ForgeTypeGitea}, OldTarget: "old-target"},
				&opcodes.ProposalUpdateSource{Proposal: forgedomain.Proposal{Data: forgedomain.ProposalData{Active: true, Body: None[gitdomain.ProposalBody](), MergeWithAPI: false, Number: 123, Source: "source", Target: "target", Title: "title", URL: "url"}, ForgeType: forgedomain.ForgeTypeForgejo}, NewBranch: "new-target", OldBranch: "old-target"},
				&opcodes.PullCurrentBranch{},
				&opcodes.PushCurrentBranch{},
				&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
				&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: "branch", ForceIfIncludes: true, TrackingBranch: "origin/branch"},
				&opcodes.PushCurrentBranchForceIgnoreError{},
				&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: "branch"},
				&opcodes.PushTags{},
				&opcodes.RebaseAbort{},
				&opcodes.RebaseBranch{Branch: "branch"},
				&opcodes.RebaseContinue{},
				&opcodes.RebaseContinueIfNeeded{},
				&opcodes.RebaseOntoRemoveDeleted{BranchToRebaseOnto: "branch-2", CommitsToRemove: "branch-1"},
				&opcodes.RebaseAncestorsUntilLocal{Branch: "branch", CommitsToRemove: Some(gitdomain.SHA("123456"))},
				&opcodes.RebaseTrackingBranch{RemoteBranch: "origin/branch", PushBranches: true},
				&opcodes.RegisterUndoablePerennialCommit{Parent: "parent"},
				&opcodes.SnapshotInitialUpdateLocalSHA{Branch: "branch", SHA: "111111"},
				&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: "branch"},
				&opcodes.StashDrop{},
				&opcodes.StashPop{},
				&opcodes.StashPopIfExists{},
				&opcodes.StashPopIfNeeded{InitialStashSize: 2},
				&opcodes.StashOpenChanges{},
				&opcodes.SyncFeatureBranchCompress{CommitMessage: Some(gitdomain.CommitMessage("commit message")), CurrentBranch: "branch", Offline: true, InitialParentName: gitdomain.NewLocalBranchNameOption("parent"), InitialParentSHA: Some(gitdomain.NewSHA("111111")), TrackingBranch: Some(gitdomain.NewRemoteBranchName("origin/branch")), PushBranches: true},
				&opcodes.SyncFeatureBranchMerge{Branch: "branch", InitialParentName: gitdomain.NewLocalBranchNameOption("original-parent"), InitialParentSHA: Some(gitdomain.NewSHA("123456")), TrackingBranch: Some(gitdomain.NewRemoteBranchName("origin/branch"))},
				&opcodes.SyncFeatureBranchRebase{Branch: "branch", ParentSHAPreviousRun: Some(gitdomain.NewSHA("111111")), PushBranches: true, TrackingBranch: Some(gitdomain.NewRemoteBranchName("origin/branch"))},
			},
			TouchedBranches: []gitdomain.BranchName{"branch-1", "branch-2"},
			UnfinishedDetails: MutableSome(&runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: "end-branch",
				EndTime:   time.Time{},
			}),
			UndoablePerennialCommits: []gitdomain.SHA{},
			FinalUndoProgram:         program.Program{},
			UndoAPIProgram:           program.Program{},
		}

		wantJSON := `
{
  "AbortProgram": [],
  "BeginBranchesSnapshot": {
    "Active": null,
    "Branches": [],
    "DetachedHead": false
  },
  "BeginConfigSnapshot": {
    "Global": {},
    "Local": {},
    "Unscoped": {}
  },
  "BeginStashSize": 0,
  "BranchInfosLastRun": null,
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
        "SHA": "111111"
      },
      "type": "BranchCurrentResetToSHA"
    },
    {
      "data": {
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "BranchCurrentResetToSHAIfNeeded"
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
        "BranchToDelete": "branch",
        "BranchToRebaseOnto": "main"
      },
      "type": "BranchLocalDeleteContent"
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
        "Branch": "branch"
      },
      "type": "BranchTrackingCreateIfLocalExists"
    },
    {
      "data": {
        "CurrentBranch": "branch"
      },
      "type": "BranchTrackingCreateIfNeeded"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "BranchTrackingDelete"
    },
    {
      "data": {
        "Branch": "branch",
        "BranchType": "feature"
      },
      "type": "BranchTypeOverrideSet"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchTypeOverrideRemove"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "BranchWithRemoteGoneDeleteIfEmptyAtRuntime"
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
      "data": {},
      "type": "ChangesUnstageAll"
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
        "SHA": "123456"
      },
      "type": "CherryPick"
    },
    {
      "data": {},
      "type": "CherryPickContinue"
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
      "type": "CommitRemove"
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
        "CommitHook": true,
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
      "data": {},
      "type": "ConflictMergePhantomFinalize"
    },
    {
      "data": {
        "CurrentBranch": "current",
        "ParentBranch": "parent",
        "ParentSHA": "123456"
      },
      "type": "ConflictMergePhantomResolveAll"
    },
    {
      "data": {
        "FilePath": "file",
        "Resolution": "ours"
      },
      "type": "ConflictResolve"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message",
        "Proposal": {
          "data": {
            "Active": true,
            "Body": "body",
            "MergeWithAPI": true,
            "Number": 123,
            "Source": "source",
            "Target": "target",
            "Title": "title",
            "URL": "url",
            "CloseSourceBranch": false,
            "Draft": false
          },
          "forge-type": "bitbucket"
        }
      },
      "type": "ConnectorProposalMerge"
    },
    {
      "data": {
        "Args": [
          "arg1",
          "arg2"
        ],
        "Executable": "executable"
      },
      "type": "ExecuteShellCommand"
    },
    {
      "data": {},
      "type": "ExitToShell"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "FetchUpstream"
    },
    {
      "data": {
        "FilePath": "file"
      },
      "type": "FileRemove"
    },
    {
      "data": {
        "FilePath": "file"
      },
      "type": "FileStage"
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
        "BranchToMerge": "branch"
      },
      "type": "MergeIntoCurrentBranch"
    },
    {
      "data": {},
      "type": "MergeAbort"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitMessage": "commit message"
      },
      "type": "MergeAlwaysProgram"
    },
    {
      "data": {},
      "type": "MergeContinue"
    },
    {
      "data": {
        "CurrentBranch": "current",
        "CurrentParent": "parent",
        "InitialParentName": "original-parent",
        "InitialParentSHA": "123456"
      },
      "type": "MergeParentResolvePhantomConflicts"
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
        "ProposalBody": null,
        "ProposalTitle": null
      },
      "type": "ProposalCreate"
    },
    {
      "data": {
        "NewBranch": "new-target",
        "OldBranch": "old-target",
        "Proposal": {
          "data": {
            "Active": true,
            "Body": "body",
            "MergeWithAPI": true,
            "Number": 123,
            "Source": "source",
            "Target": "target",
            "Title": "title",
            "URL": "url"
          },
          "forge-type": "gitlab"
        }
      },
      "type": "ProposalUpdateTarget"
    },
    {
      "data": {
        "Branch": "branch",
        "OldTarget": "old-target",
        "Proposal": {
          "data": {
            "Active": true,
            "Body": "body",
            "MergeWithAPI": true,
            "Number": 123,
            "Source": "source",
            "Target": "target",
            "Title": "title",
            "URL": "url"
          },
          "forge-type": "gitea"
        }
      },
      "type": "ProposalUpdateTargetToGrandParent"
    },
    {
      "data": {
        "NewBranch": "new-target",
        "OldBranch": "old-target",
        "Proposal": {
          "data": {
            "Active": true,
            "Body": null,
            "MergeWithAPI": false,
            "Number": 123,
            "Source": "source",
            "Target": "target",
            "Title": "title",
            "URL": "url"
          },
          "forge-type": "forgejo"
        }
      },
      "type": "ProposalUpdateSource"
    },
    {
      "data": {},
      "type": "PullCurrentBranch"
    },
    {
      "data": {},
      "type": "PushCurrentBranch"
    },
    {
      "data": {
        "ForceIfIncludes": true
      },
      "type": "PushCurrentBranchForce"
    },
    {
      "data": {
        "CurrentBranch": "branch",
        "ForceIfIncludes": true,
        "TrackingBranch": "origin/branch"
      },
      "type": "PushCurrentBranchForceIfNeeded"
    },
    {
      "data": {},
      "type": "PushCurrentBranchForceIgnoreError"
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
        "BranchToRebaseOnto": "branch-2",
        "CommitsToRemove": "branch-1"
      },
      "type": "RebaseOntoRemoveDeleted"
    },
    {
      "data": {
        "Branch": "branch",
        "CommitsToRemove": "123456"
      },
      "type": "RebaseAncestorsUntilLocal"
    },
    {
      "data": {
        "PushBranches": true,
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
      "type": "StashPopIfExists"
    },
    {
      "data": {
        "InitialStashSize": 2
      },
      "type": "StashPopIfNeeded"
    },
    {
      "data": {},
      "type": "StashOpenChanges"
    },
    {
      "data": {
        "CommitMessage": "commit message",
        "CurrentBranch": "branch",
        "InitialParentName": "parent",
        "InitialParentSHA": "111111",
        "Offline": true,
        "PushBranches": true,
        "TrackingBranch": "origin/branch"
      },
      "type": "SyncFeatureBranchCompress"
    },
    {
      "data": {
        "Branch": "branch",
        "InitialParentName": "original-parent",
        "InitialParentSHA": "123456",
        "TrackingBranch": "origin/branch"
      },
      "type": "SyncFeatureBranchMerge"
    },
    {
      "data": {
        "Branch": "branch",
        "ParentSHAPreviousRun": "111111",
        "PushBranches": true,
        "TrackingBranch": "origin/branch"
      },
      "type": "SyncFeatureBranchRebase"
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
		err := runstate.Save(runState, repoRoot)
		must.NoError(t, err)
		filepath, err := state.FilePath(repoRoot, state.FileTypeRunstate)
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
