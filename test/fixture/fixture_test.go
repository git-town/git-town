package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/test/asserts"
	"github.com/git-town/git-town/v13/test/fixture"
	"github.com/git-town/git-town/v13/test/git"
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
		cloned.DevRepo.PushBranchToRemote(gitdomain.NewLocalBranchName("main"), gitdomain.OriginRemote)
	})

	t.Run("NewStandardFixture", func(t *testing.T) {
		t.Parallel()
		gitEnvRootDir := t.TempDir()
		result := fixture.NewStandardFixture(gitEnvRootDir)
		// verify the origin repo
		asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
		branch, err := result.OriginRepo.CurrentBranch()
		must.NoError(t, err)
		must.EqOp(t, gitdomain.NewLocalBranchName("initial"), branch)
		// verify the developer repo
		asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
		assertHasGitConfiguration(t, gitEnvRootDir)
		branch, err = result.DevRepo.CurrentBranch()
		must.NoError(t, err)
		must.EqOp(t, gitdomain.NewLocalBranchName("main"), branch)
	})

	t.Run("Branches", func(t *testing.T) {
		t.Run("different branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			gitEnv := fixture.NewStandardFixture(filepath.Join(dir, ""))
			// create the branches
			gitEnv.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("d1"), gitdomain.NewLocalBranchName("main"))
			gitEnv.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("d2"), gitdomain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(gitdomain.NewLocalBranchName("o1"), gitdomain.NewLocalBranchName("initial"))
			gitEnv.OriginRepo.CreateBranch(gitdomain.NewLocalBranchName("o2"), gitdomain.NewLocalBranchName("initial"))
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
			gitEnv.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("b1"), gitdomain.NewLocalBranchName("main"))
			gitEnv.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("b2"), gitdomain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(gitdomain.NewLocalBranchName("b1"), gitdomain.NewLocalBranchName("main"))
			gitEnv.OriginRepo.CreateBranch(gitdomain.NewLocalBranchName("b2"), gitdomain.NewLocalBranchName("main"))
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
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "lc",
				FileName:    "local-file",
				Locations:   []string{"local"},
				Message:     "local commit",
			},
			{
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "rc",
				FileName:    "origin-file",
				Locations:   []string{"origin"},
				Message:     "origin commit",
			},
			{
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "lrc",
				FileName:    "loc-rem-file",
				Locations:   []string{"local", "origin"},
				Message:     "local and origin commit",
			},
		})
		// verify local commits
		commits := cloned.DevRepo.Commits([]string{"FILE NAME", "FILE CONTENT"}, gitdomain.NewLocalBranchName("main"))
		must.Len(t, 2, commits)
		must.EqOp(t, "local commit", commits[0].Message)
		must.EqOp(t, "local-file", commits[0].FileName)
		must.EqOp(t, "lc", commits[0].FileContent)
		must.EqOp(t, "local and origin commit", commits[1].Message)
		must.EqOp(t, "loc-rem-file", commits[1].FileName)
		must.EqOp(t, "lrc", commits[1].FileContent)
		// verify origin commits
		commits = cloned.OriginRepo.Commits([]string{"FILE NAME", "FILE CONTENT"}, gitdomain.NewLocalBranchName("main"))
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
		must.EqOp(t, gitdomain.NewLocalBranchName("initial"), branch)
	})

	t.Run("CreateOriginBranch", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoizedGitEnv := fixture.NewStandardFixture(filepath.Join(dir, "memoized"))
		cloned := fixture.CloneFixture(memoizedGitEnv, filepath.Join(dir, "cloned"))
		// create the origin branch
		cloned.OriginRepo.CreateBranch(gitdomain.NewLocalBranchName("b1"), gitdomain.NewLocalBranchName("main"))
		// verify it is in the origin branches
		branches, err := cloned.OriginRepo.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("main"))
		must.NoError(t, err)
		must.SliceContains(t, branches.Strings(), "b1")
		// verify it isn't in the local branches
		branches, err = cloned.DevRepo.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("main"))
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
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "one",
				FileName:    "local-origin.md",
				Message:     "local-origin",
			})
			cloned.DevRepo.PushBranchToRemote(gitdomain.NewLocalBranchName("main"), gitdomain.OriginRemote)
			cloned.OriginRepo.CreateCommit(git.Commit{
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "two",
				FileName:    "origin.md",
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
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "one",
				FileName:    "local.md",
				Message:     "local",
			})
			cloned.UpstreamRepo.CreateCommit(git.Commit{
				Branch:      gitdomain.NewLocalBranchName("main"),
				FileContent: "two",
				FileName:    "upstream.md",
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
}

func assertHasGitConfiguration(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	must.NoError(t, err)
	for e := range entries {
		if entries[e].Name() == ".gitconfig" {
			return
		}
	}
	t.Fatalf(".gitconfig not found in %q", dir)
}
