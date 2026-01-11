package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/fixture"
	"github.com/git-town/git-town/v22/internal/test/testgit"
	"github.com/git-town/git-town/v22/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestFixture(t *testing.T) {
	t.Parallel()

	t.Run("Branches", func(t *testing.T) {
		t.Run("different branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			fixture := fixture.NewMemoized(filepath.Join(dir, "")).AsFixture()
			// create the branches
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("d1", "main")
			devRepo.CreateBranch("d2", "main")
			originRepo := fixture.OriginRepo.GetOrPanic()
			originRepo.CreateBranch("o1", "initial")
			originRepo.CreateBranch("o2", "initial")
			// get branches
			table := fixture.Branches()
			// verify
			expected := "| REPOSITORY | BRANCHES     |\n| local      | main, d1, d2 |\n| origin     | main, o1, o2 |\n"
			must.EqOp(t, expected, table.String())
		})

		t.Run("same branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			fixture := fixture.NewMemoized(filepath.Join(dir, "")).AsFixture()
			// create the branches
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("b1", "main")
			devRepo.CreateBranch("b2", "main")
			originRepo := fixture.OriginRepo.GetOrPanic()
			originRepo.CreateBranch("b1", "main")
			originRepo.CreateBranch("b2", "main")
			// get branches
			table := fixture.Branches()
			// verify
			expected := "| REPOSITORY    | BRANCHES     |\n| local, origin | main, b1, b2 |\n"
			must.EqOp(t, expected, table.String())
		})
	})

	t.Run("CloneFixture", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		memoized := fixture.NewMemoized(filepath.Join(dir, "memoized"))
		cloned := memoized.CloneInto(filepath.Join(dir, "cloned"))
		asserts.IsGitRepo(t, filepath.Join(dir, "cloned", "origin"))
		asserts.IsGitRepo(t, filepath.Join(dir, "cloned", "developer"))
		asserts.BranchExists(t, filepath.Join(dir, "cloned", "developer"), "main")
		// check pushing
		devRepo := cloned.DevRepo.GetOrPanic()
		devRepo.PushBranchToRemote("main", gitdomain.RemoteOrigin)
	})

	t.Run("CommitTable", func(t *testing.T) {
		t.Run("without upstream repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			memoized := fixture.NewMemoized(filepath.Join(dir, "memoized"))
			cloned := memoized.CloneInto(filepath.Join(dir, "cloned"))
			clonedDevRepo := cloned.DevRepo.GetOrPanic()
			// create a few commits
			clonedDevRepo.CreateCommit(testgit.Commit{
				Branch:      "main",
				FileContent: "one",
				FileName:    "local-origin.md",
				Message:     "local-origin",
			})
			clonedDevRepo.PushBranchToRemote("main", gitdomain.RemoteOrigin)
			cloned.OriginRepo.GetOrPanic().CreateCommit(testgit.Commit{
				Branch:      "main",
				FileContent: "two",
				FileName:    "origin.md",
				Message:     "2",
			})
			// get the CommitTable
			table := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			must.Eq(t, table.Cells, [][]string{
				{"LOCATION", "FILE NAME", "FILE CONTENT"},
				{"local, origin", "local-origin.md", "one"},
				{"origin", "origin.md", "two"},
			})
		})

		t.Run("with upstream repo", func(t *testing.T) {
			t.Parallel()
			// create Fixture instance
			dir := t.TempDir()
			memoized := fixture.NewMemoized(filepath.Join(dir, "memoized"))
			cloned := memoized.CloneInto(filepath.Join(dir, "cloned"))
			cloned.AddUpstream()
			// create a few commits
			cloned.DevRepo.GetOrPanic().CreateCommit(testgit.Commit{
				Branch:      "main",
				FileContent: "one",
				FileName:    "local.md",
				Message:     "local",
			})
			cloned.UpstreamRepo.GetOrPanic().CreateCommit(testgit.Commit{
				Branch:      "main",
				FileContent: "two",
				FileName:    "upstream.md",
				Message:     "2",
			})
			// get the CommitTable
			table := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			must.Eq(t, table.Cells, [][]string{
				{"LOCATION", "FILE NAME", "FILE CONTENT"},
				{"local", "local.md", "one"},
				{"upstream", "upstream.md", "two"},
			})
		})
	})

	t.Run("CreateCommits", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoized := fixture.NewMemoized(filepath.Join(dir, "memoized"))
		cloned := memoized.CloneInto(filepath.Join(dir, "cloned"))
		// create the commits
		mainBranch := gitdomain.NewLocalBranchName("main")
		cloned.CreateCommits([]testgit.Commit{
			{
				Branch:      mainBranch,
				FileContent: "local and origin content",
				FileName:    "loc-rem-file",
				Locations:   testgit.Locations{testgit.LocationLocal, testgit.LocationOrigin},
				Message:     "local and origin commit",
			},
			{
				Branch:      mainBranch,
				FileContent: "local content",
				FileName:    "local-file",
				Locations:   testgit.Locations{testgit.LocationLocal},
				Message:     "local commit",
			},
			{
				Branch:      mainBranch,
				FileContent: "origin content",
				FileName:    "origin-file",
				Locations:   testgit.Locations{testgit.LocationOrigin},
				Message:     "origin commit",
			},
		})
		// verify local commits
		commits := cloned.DevRepo.GetOrPanic().Commits([]string{"FILE NAME", "FILE CONTENT"}, cloned.DevRepo.Value.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 2, commits)
		must.EqOp(t, "local and origin commit", commits[0].Message)
		must.EqOp(t, "loc-rem-file", commits[0].FileName)
		must.EqOp(t, "local and origin content", commits[0].FileContent)
		must.EqOp(t, "local commit", commits[1].Message)
		must.EqOp(t, "local-file", commits[1].FileName)
		must.EqOp(t, "local content", commits[1].FileContent)
		// verify origin commits
		commits = cloned.OriginRepo.GetOrPanic().Commits([]string{"FILE NAME", "FILE CONTENT"}, cloned.DevRepo.Value.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 2, commits)
		must.EqOp(t, "local and origin commit", commits[0].Message)
		must.EqOp(t, "loc-rem-file", commits[0].FileName)
		must.EqOp(t, "local and origin content", commits[0].FileContent)
		must.EqOp(t, "origin commit", commits[1].Message)
		must.EqOp(t, "origin-file", commits[1].FileName)
		must.EqOp(t, "origin content", commits[1].FileContent)
		// verify origin is at "initial" branch
		branch, err := cloned.OriginRepo.GetOrPanic().Git.CurrentBranch(cloned.DevRepo.GetOrPanic().TestRunner)
		must.NoError(t, err)
		must.EqOp(t, "initial", branch.GetOrPanic())
	})

	t.Run("CreateOriginBranch", func(t *testing.T) {
		t.Parallel()
		// create Fixture instance
		dir := t.TempDir()
		memoized := fixture.NewMemoized(filepath.Join(dir, "memoized"))
		cloned := memoized.CloneInto(filepath.Join(dir, "cloned"))
		// create the origin branch
		cloned.OriginRepo.GetOrPanic().CreateBranch("b1", "main")
		// verify it is in the origin branches
		branches, err := cloned.OriginRepo.GetOrPanic().LocalBranchesMainFirst("main")
		must.NoError(t, err)
		must.SliceContains(t, branches.AllBranches.Strings(), "b1")
		// verify it isn't in the local branches
		branches, err = cloned.DevRepo.GetOrPanic().LocalBranchesMainFirst("main")
		must.NoError(t, err)
		must.SliceNotContains(t, branches.AllBranches.Strings(), "b1")
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
