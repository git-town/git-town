package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/configfile"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("RenderPerennialBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("no perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames()
			have := configfile.RenderPerennialBranches(give)
			want := "[]"
			must.EqOp(t, want, have)
		})
		t.Run("one perennial branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one")
			have := configfile.RenderPerennialBranches(give)
			want := `["one"]`
			must.EqOp(t, want, have)
		})
		t.Run("multiple perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := configfile.RenderPerennialBranches(give)
			want := `["one", "two"]`
			must.EqOp(t, want, have)
		})
	})

	t.Run("RenderTOML", func(t *testing.T) {
		t.Parallel()
		give := configdomain.DefaultConfig()
		give.MainBranch = Some(gitdomain.NewLocalBranchName("main"))
		give.PerennialBranches = gitdomain.NewLocalBranchNames("one", "two")
		have := configfile.RenderTOML(&give)
		want := `
# Git Town configuration file
#
# Run "git town config setup" to add additional entries
# to this file after updating Git Town.
#
# The "push-hook" setting determines whether Git Town
# permits or prevents Git hooks while pushing branches.
# Hooks are enabled by default. If your Git hooks are slow,
# you can disable them to speed up branch syncing.
#
# When disabled, Git Town pushes using the "--no-verify" switch.
# More info at https://www.git-town.com/preferences/push-hook.
push-hook = true

# Should Git Town push the new branches it creates
# immediately to origin even if they are empty?
#
# When enabled, you can run "git push" right away
# but creating new branches is slower and
# it triggers an unnecessary CI run on the empty branch.
#
# When disabled, many Git Town commands execute faster
# and Git Town will create the missing tracking branch
# on the first run of "git sync".
push-new-branches = false

# Should "git ship" delete the tracking branch?
# You want to disable this if your code hosting platform
# (GitHub, GitLab, etc) deletes head branches when
# merging pull requests through its UI.
ship-delete-tracking-branch = true

# Should "git ship" sync branches before shipping them?
#
# Guidance: enable when shipping branches locally on your machine
# and disable when shipping feature branches via the code hosting
# API or web UI.
#
# When enabled, branches are always fully up to date when shipped
# and you get a chance to resolve merge conflicts
# between the feature branch to ship and the main development branch
# on the feature branch. This helps keep the main branch green.
# But this also triggers another CI run and delays shipping.
sync-before-ship = false

# Should "git sync" also fetch updates from the upstream remote?
#
# If an "upstream" remote exists, and this setting is enabled,
# "git sync" will also update the local main branch
# with commits from the main branch at the upstream remote.
#
# This is useful if the repository you work on is a fork,
# and you want to keep it in sync with the repo it was forked from.
sync-upstream = true

[branches]

# The main branch is the branch from which you cut new feature branches,
# and into which you ship feature branches when they are done.
# This branch is often called "main", "master", or "development".
main = "main"

# Perennial branches are long-lived branches.
# They are never shipped and have no ancestors.
# Typically, perennial branches have names like
# "development", "staging", "qa", "production", etc.
#
# See also the "perennial-regex" setting.
perennials = ["one", "two"]

# All branches whose names match this regular expression
# are also considered perennial branches.
#
# If you are not sure, leave this empty.
perennial-regex = ""

[hosting]

# Knowing the type of code hosting platform allows Git Town
# to open browser URLs and talk to the code hosting API.
# Most people can leave this on "auto-detect".
# Only change this if your code hosting server uses as custom URL.
# platform = ""

# When using SSH identities, define the hostname
# of your source code repository. Only change this
# if the auto-detection does not work for you.
# origin-hostname = ""

[sync-strategy]

# How should Git Town synchronize feature branches?
# Feature branches are short-lived branches cut from
# the main branch and shipped back into the main branch.
# Typically you develop features and bug fixes on them,
# hence their name.
feature-branches = "merge"

# How should Git Town synchronize perennial branches?
# Perennial branches have no parent branch.
# The only updates they receive are additional commits
# made to their tracking branch somewhere else.
perennial-branches = "rebase"
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		give := configdomain.DefaultConfig()
		give.MainBranch = Some(gitdomain.NewLocalBranchName("main"))
		err := configfile.Save(&give)
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := `
# Git Town configuration file
#
# Run "git town config setup" to add additional entries
# to this file after updating Git Town.
#
# The "push-hook" setting determines whether Git Town
# permits or prevents Git hooks while pushing branches.
# Hooks are enabled by default. If your Git hooks are slow,
# you can disable them to speed up branch syncing.
#
# When disabled, Git Town pushes using the "--no-verify" switch.
# More info at https://www.git-town.com/preferences/push-hook.
push-hook = true

# Should Git Town push the new branches it creates
# immediately to origin even if they are empty?
#
# When enabled, you can run "git push" right away
# but creating new branches is slower and
# it triggers an unnecessary CI run on the empty branch.
#
# When disabled, many Git Town commands execute faster
# and Git Town will create the missing tracking branch
# on the first run of "git sync".
push-new-branches = false

# Should "git ship" delete the tracking branch?
# You want to disable this if your code hosting platform
# (GitHub, GitLab, etc) deletes head branches when
# merging pull requests through its UI.
ship-delete-tracking-branch = true

# Should "git ship" sync branches before shipping them?
#
# Guidance: enable when shipping branches locally on your machine
# and disable when shipping feature branches via the code hosting
# API or web UI.
#
# When enabled, branches are always fully up to date when shipped
# and you get a chance to resolve merge conflicts
# between the feature branch to ship and the main development branch
# on the feature branch. This helps keep the main branch green.
# But this also triggers another CI run and delays shipping.
sync-before-ship = false

# Should "git sync" also fetch updates from the upstream remote?
#
# If an "upstream" remote exists, and this setting is enabled,
# "git sync" will also update the local main branch
# with commits from the main branch at the upstream remote.
#
# This is useful if the repository you work on is a fork,
# and you want to keep it in sync with the repo it was forked from.
sync-upstream = true

[branches]

# The main branch is the branch from which you cut new feature branches,
# and into which you ship feature branches when they are done.
# This branch is often called "main", "master", or "development".
main = "main"

# Perennial branches are long-lived branches.
# They are never shipped and have no ancestors.
# Typically, perennial branches have names like
# "development", "staging", "qa", "production", etc.
#
# See also the "perennial-regex" setting.
perennials = []

# All branches whose names match this regular expression
# are also considered perennial branches.
#
# If you are not sure, leave this empty.
perennial-regex = ""

[hosting]

# Knowing the type of code hosting platform allows Git Town
# to open browser URLs and talk to the code hosting API.
# Most people can leave this on "auto-detect".
# Only change this if your code hosting server uses as custom URL.
# platform = ""

# When using SSH identities, define the hostname
# of your source code repository. Only change this
# if the auto-detection does not work for you.
# origin-hostname = ""

[sync-strategy]

# How should Git Town synchronize feature branches?
# Feature branches are short-lived branches cut from
# the main branch and shipped back into the main branch.
# Typically you develop features and bug fixes on them,
# hence their name.
feature-branches = "merge"

# How should Git Town synchronize perennial branches?
# Perennial branches have no parent branch.
# The only updates they receive are additional commits
# made to their tracking branch somewhere else.
perennial-branches = "rebase"
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("TOMLComment", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("")
			want := ""
			must.Eq(t, want, have)
		})
		t.Run("single line", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1")
			want := "# line 1"
			must.Eq(t, want, have)
		})
		t.Run("multiple lines", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1\nline 2\nline 3")
			want := "# line 1\n# line 2\n# line 3"
			must.Eq(t, want, have)
		})
		t.Run("multiple lines with terminating newline", func(t *testing.T) {
			t.Parallel()
			have := configfile.TOMLComment("line 1\nline 2\nline 3\n")
			want := "# line 1\n# line 2\n# line 3\n#"
			must.Eq(t, want, have)
		})
	})
}
