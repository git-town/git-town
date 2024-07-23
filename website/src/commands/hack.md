# git hack &lt;branch&gt;

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch.md) and brings all
uncommitted changes over to it.

When running without uncommitted changes in your workspace, it also
[syncs](sync.md) the main branch to ensure you develop on top of the current
state of the repository. If the workspace contains uncommitted changes,
`git hack` does not perform this sync to let you commit your open changes first
and then sync manually.

### Configuration

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

If [push-new-branches](../preferences/push-new-branches.md) is set, `git hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git hack` run fast. The first run of `git sync`
will create the remote tracking branch.

### Arguments

When given a non-existing branch name, `git hack` creates a new feature branch
with the main branch as its parent. Adding the `--prototype` or `-p` switch
makes it create a [prototype branch](../branch-types.md#prototype-branches)).

When given an existing contribution, observed, parked, or prototype branch,
`git hack` converts that branch to a feature branch.

When given no arguments, `git hack` converts the current contribution, observed,
parked, or prototype branch into a proper feature branch.
