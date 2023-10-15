package persistence_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/persistence"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/git-town/git-town/v9/src/vm/runstate"
	"github.com/git-town/git-town/v9/src/vm/shared"
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
			rootDir := domain.NewRepoRootDir(give)
			have := persistence.SanitizePath(rootDir)
			must.EqOp(t, want, have)
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		t.Parallel()
		runState := runstate.RunState{
			Command:      "command",
			IsAbort:      true,
			IsUndo:       true,
			AbortProgram: program.Program{},
			RunProgram: program.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
					&opcode.AddToPerennialBranches{Branch: domain.NewLocalBranchName("branch")},
					&opcode.ChangeParent{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch")},
					&opcode.CommitOpenChanges{},
					&opcode.ConnectorMergeProposal{
						Branch:          domain.NewLocalBranchName("branch"),
						CommitMessage:   "commit message",
						ProposalMessage: "proposal message",
						ProposalNumber:  123,
					},
					&opcode.ContinueMerge{},
					&opcode.ContinueRebase{},
					&opcode.CreateBranch{
						Branch:        domain.NewLocalBranchName("branch"),
						StartingPoint: domain.NewSHA("123456").Location(),
					},
					&opcode.CreateProposal{Branch: domain.NewLocalBranchName("branch")},
					&opcode.CreateRemoteBranch{
						Branch:     domain.NewLocalBranchName("branch"),
						NoPushHook: true,
						SHA:        domain.NewSHA("123456"),
					},
					&opcode.CreateTrackingBranch{
						Branch:     domain.NewLocalBranchName("branch"),
						NoPushHook: true,
					},
					&opcode.DeleteLocalBranch{
						Branch: domain.NewLocalBranchName("branch"),
						Force:  false,
					},
					&opcode.DeleteRemoteBranch{
						Branch: domain.NewRemoteBranchName("origin/branch"),
					},
					&opcode.DeleteParentBranch{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&opcode.DeleteTrackingBranch{
						Branch: domain.NewRemoteBranchName("origin/branch"),
					},
					&opcode.DiscardOpenChanges{},
					&opcode.EnsureHasShippableChanges{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&opcode.FetchUpstream{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&opcode.ForcePushCurrentBranch{
						NoPushHook: true,
					},
					&opcode.Merge{Branch: domain.NewBranchName("branch")},
					&opcode.MergeParent{CurrentBranch: domain.NewLocalBranchName("branch")},
					&opcode.PreserveCheckoutHistory{
						InitialBranch:                     domain.NewLocalBranchName("initial-branch"),
						InitialPreviouslyCheckedOutBranch: domain.NewLocalBranchName("initial-previous-branch"),
						MainBranch:                        domain.NewLocalBranchName("main"),
					},
					&opcode.PullCurrentBranch{},
					&opcode.PushCurrentBranch{
						CurrentBranch: domain.NewLocalBranchName("branch"),
						NoPushHook:    true,
					},
					&opcode.PushTags{},
					&opcode.RebaseBranch{Branch: domain.NewBranchName("branch")},
					&opcode.RebaseParent{CurrentBranch: domain.NewLocalBranchName("branch")},
					&opcode.RemoveFromPerennialBranches{
						Branch: domain.NewLocalBranchName("branch"),
					},
					&opcode.RemoveGlobalConfig{
						Key: config.KeyOffline,
					},
					&opcode.RemoveLocalConfig{
						Key: config.KeyOffline,
					},
					&opcode.ResetCurrentBranchToSHA{
						Hard:        true,
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
					},
					&opcode.RestoreOpenChanges{},
					&opcode.RevertCommit{
						SHA: domain.NewSHA("123456"),
					},
					&opcode.SetGlobalConfig{
						Key:   config.KeyOffline,
						Value: "1",
					},
					&opcode.SetLocalConfig{
						Key:   config.KeyOffline,
						Value: "1",
					},
					&opcode.SetParent{
						Branch: domain.NewLocalBranchName("branch"),
						Parent: domain.NewLocalBranchName("parent"),
					},
					&opcode.SkipCurrentBranch{},
					&opcode.SquashMerge{
						Branch:        domain.NewLocalBranchName("branch"),
						CommitMessage: "commit message",
						Parent:        domain.NewLocalBranchName("parent"),
					},
					&opcode.StashOpenChanges{},
					&opcode.UpdateProposalTarget{
						ProposalNumber: 123,
						NewTarget:      domain.NewLocalBranchName("new-target"),
					},
				},
			},
			UndoProgram: program.Program{},
			UnfinishedDetails: &runstate.UnfinishedRunStateDetails{
				CanSkip:   true,
				EndBranch: domain.NewLocalBranchName("end-branch"),
				EndTime:   time.Time{},
			},
			InitialActiveBranch:      domain.NewLocalBranchName("initial"),
			UndoablePerennialCommits: []domain.SHA{},
		}

		wantJSON := `
{
  "Command": "command",
  "IsAbort": true,
  "IsUndo": true,
  "AbortProgram": [],
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
      "type": "AddToPerennialBranches"
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
      "data": {},
      "type": "CommitOpenChanges"
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
      "data": {
        "Branch": "branch",
        "StartingPoint": "123456"
      },
      "type": "CreateBranch"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "CreateProposal"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true,
        "SHA": "123456"
      },
      "type": "CreateRemoteBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "NoPushHook": true
      },
      "type": "CreateTrackingBranch"
    },
    {
      "data": {
        "Branch": "branch",
        "Force": false
      },
      "type": "DeleteLocalBranch"
    },
    {
      "data": {
        "Branch": "origin/branch"
      },
      "type": "DeleteRemoteBranch"
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
        "NoPushHook": true
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
        "CurrentBranch": "branch"
      },
      "type": "MergeParent"
    },
    {
      "data": {
        "InitialBranch": "initial-branch",
        "InitialPreviouslyCheckedOutBranch": "initial-previous-branch",
        "MainBranch": "main"
      },
      "type": "PreserveCheckoutHistory"
    },
    {
      "data": {},
      "type": "PullCurrentBranch"
    },
    {
      "data": {
        "CurrentBranch": "branch",
        "NoPushHook": true
      },
      "type": "PushCurrentBranch"
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
        "CurrentBranch": "branch"
      },
      "type": "RebaseParent"
    },
    {
      "data": {
        "Branch": "branch"
      },
      "type": "RemoveFromPerennialBranches"
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
        "Hard": true,
        "MustHaveSHA": "222222",
        "SetToSHA": "111111"
      },
      "type": "ResetCurrentBranchToSHA"
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
      "type": "StashOpenChanges"
    },
    {
      "data": {
        "ProposalNumber": 123,
        "NewTarget": "new-target"
      },
      "type": "UpdateProposalTarget"
    }
  ],
  "UndoProgram": [],
  "InitialActiveBranch": "initial",
  "FinalUndoProgram": [],
  "UnfinishedDetails": {
    "CanSkip": true,
    "EndBranch": "end-branch",
    "EndTime": "0001-01-01T00:00:00Z"
  },
  "UndoablePerennialCommits": []
}`[1:]

		repoRoot := domain.NewRepoRootDir("/path/to/git-town-unit-tests")
		err := persistence.Save(&runState, repoRoot)
		must.NoError(t, err)
		filepath, err := persistence.FilePath(repoRoot)
		must.NoError(t, err)
		content, err := os.ReadFile(filepath)
		must.NoError(t, err)
		must.EqOp(t, wantJSON, string(content))
		var newState runstate.RunState
		err = json.Unmarshal(content, &newState)
		must.NoError(t, err)
		must.Eq(t, runState, newState)
	})
}
