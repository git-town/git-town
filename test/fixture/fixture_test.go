package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/test/asserts"
	"github.com/git-town/git-town/v11/test/fixture"
	"github.com/git-town/git-town/v11/test/git"
	"github.com/shoenig/test/must"
)

func TestFixture(t *testing.T) {
	t.Parallel()

	t.Run("CloneFixture", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
		cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
		asserts.IsGitRepo(t, filepath.Join(dir, "cloned", "origin"))
		asserts.IsGitRepo(t, filepath.Join(dir, "cloned", "developer"))
		asserts.BranchExists(t, filepath.Join(dir, "cloned", "developer"), "main")
		// check pushing
		cloned.DevRepo.PushBranchToRemote(domain.NewLocalBranchName("main"), domain.OriginRemote)
	})

	t.Run("NewStandardFixture", func(t *testing.T) {
		t.Parallel()
		gitEnvRootDir := t.TempDir()
		result := fixture.NewStandardFixture(gitEnvRootDir)
		// verify the origin repo
		asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
		branch, err := result.OriginRepo.CurrentBranch()
		must.NoError(t, err)
		must.EqOp(t, domain.NewLocalBranchName("initial"), branch)
		// verify the developer repo
		asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
		asserts.HasGitConfiguration(t, gitEnvRootDir)
		branch, err = result.DevRepo.CurrentBranch()
		must.NoError(t, err)
		must.EqOp(t, domain.NewLocalBranchName("main"), branch)
	})

	t.Run("Branches", func(t *testing.T) {
		t.Run("different branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			gitEnv := fixture.NewStandardFixture(filepath.Join(dir, ""))
			// create the branches
			gitEnv.DevRepo.CreateBranch(domain.NewLocalBranchName("d1"), domain.NewLocalBranchName("main"))
			gitEnv.DevRepo.CreateBranch(domain.NewLocalBranchName("d2"), domain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(domain.NewLocalBranchName("o1"), domain.NewLocalBranchName("initial"))
			gitEnv.OriginRepo.CreateBranch(domain.NewLocalBranchName("o2"), domain.NewLocalBranchName("initial"))
			// get branches
			table := gitEnv.Branches()
			// verify
			expected := "| REPOSITORY | BRANCHES     |\n| local      | main, d1, d2 |\n| origin     | main, o1, o2 |\n"
			must.EqOp(t, expected, table.String())
		})

		t.Run("same branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			gitEnv := fixture.NewStandardFixture(filepath.Join(dir, ""))
			// create the branches
			gitEnv.DevRepo.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("main"))
			gitEnv.DevRepo.CreateBranch(domain.NewLocalBranchName("b2"), domain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(domain.NewLocalBranchName("b2"), domain.NewLocalBranchName("main"))
			// get branches
			table := gitEnv.Branches()
			// verify
			expected := "| REPOSITORY    | BRANCHES     |\n| local, origin | main, b1, b2 |\n"
			must.EqOp(t, expected, table.String())
		})
	})

	t.Run("CreateCommits", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
		cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
		// create the commits
		cloned.CreateCommits([]git.Commit{
			{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "local-file",
				FileContent: "lc",
				Locations:   []string{"local"},
				Message:     "local commit",
			},
			{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "origin-file",
				FileContent: "rc",
				Locations:   []string{"origin"},
				Message:     "origin commit",
			},
			{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "loc-rem-file",
				FileContent: "lrc",
				Locations:   []string{"local", "origin"},
				Message:     "local and origin commit",
			},
		})
		// verify local commits
		commits := cloned.DevRepo.Commits([]string{"FILE NAME", "FILE CONTENT"}, domain.NewLocalBranchName("main"))
		must.Len(t, 2, commits)
		must.EqOp(t, "local commit", commits[0].Message)
		must.EqOp(t, "local-file", commits[0].FileName)
		must.EqOp(t, "lc", commits[0].FileContent)
		must.EqOp(t, "local and origin commit", commits[1].Message)
		must.EqOp(t, "loc-rem-file", commits[1].FileName)
		must.EqOp(t, "lrc", commits[1].FileContent)
		// verify origin commits
		commits = cloned.OriginRepo.Commits([]string{"FILE NAME", "FILE CONTENT"}, domain.NewLocalBranchName("main"))
		must.Len(t, 2, commits)
		must.EqOp(t, "origin commit", commits[0].Message)
		must.EqOp(t, "origin-file", commits[0].FileName)
		must.EqOp(t, "rc", commits[0].FileContent)
		must.EqOp(t, "local and origin commit", commits[1].Message)
		must.EqOp(t, "loc-rem-file", commits[1].FileName)
		must.EqOp(t, "lrc", commits[1].FileContent)
		// verify origin is at "initial" branch
		branch, err := cloned.OriginRepo.CurrentBranch()
		must.NoError(t, err)
		must.EqOp(t, domain.NewLocalBranchName("initial"), branch)
	})

	t.Run("CreateOriginBranch", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
		cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
		// create the origin branch
		cloned.CreateOriginBranch("b1", "main")
		// verify it is in the origin branches
		branches, err := cloned.OriginRepo.LocalBranchesMainFirst(domain.NewLocalBranchName("main"))
		must.NoError(t, err)
		must.SliceContains(t, branches.Strings(), "b1")
		// verify it isn't in the local branches
		branches, err = cloned.DevRepo.LocalBranchesMainFirst(domain.NewLocalBranchName("main"))
		must.NoError(t, err)
		must.SliceNotContains(t, branches.Strings(), "b1")
	})

	t.Run("CommitTable", func(t *testing.T) {
		t.Run("without upstream repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
			cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
			// create a few commits
			cloned.DevRepo.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "local-origin.md",
				FileContent: "one",
				Message:     "local-origin",
			})
			cloned.DevRepo.PushBranchToRemote(domain.NewLocalBranchName("main"), domain.OriginRemote)
			cloned.OriginRepo.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "origin.md",
				FileContent: "two",
				Message:     "2",
			})
			// get the CommitTable
			table := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			must.Len(t, 3, table.Cells)
			must.EqOp(t, table.Cells[1][0], "local, origin")
			must.EqOp(t, table.Cells[1][1], "local-origin.md")
			must.EqOp(t, table.Cells[1][2], "one")
			must.EqOp(t, table.Cells[2][0], "origin")
			must.EqOp(t, table.Cells[2][1], "origin.md")
			must.EqOp(t, table.Cells[2][2], "two")
		})

		t.Run("with upstream repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
			cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
			cloned.AddUpstream()
			// create a few commits
			cloned.DevRepo.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "local.md",
				FileContent: "one",
				Message:     "local",
			})
			cloned.UpstreamRepo.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("main"),
				FileName:    "upstream.md",
				FileContent: "two",
				Message:     "2",
			})
			// get the CommitTable
			table := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			must.Len(t, 3, table.Cells)
			must.EqOp(t, table.Cells[1][0], "local")
			must.EqOp(t, table.Cells[1][1], "local.md")
			must.EqOp(t, table.Cells[1][2], "one")
			must.EqOp(t, table.Cells[2][0], "upstream")
			must.EqOp(t, table.Cells[2][1], "upstream.md")
			must.EqOp(t, table.Cells[2][2], "two")
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
		cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
		// remove it
		cloned.Remove()
		// verify
		_, err := os.Stat(cloned.Dir)
		must.True(t, os.IsNotExist(err))
	})
}
