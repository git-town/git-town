# git town hack

> _git town hack [--prototype] [branch-name...]_

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch.md) and brings all
uncommitted changes over to it.

When running without uncommitted changes in your workspace, it also
[syncs](sync.md) the main branch to ensure you develop on top of the current
state of the repository. If the workspace contains uncommitted changes,
`git hack` does not perform this sync to let you commit your open changes first
and then sync later.

### Positional arguments

When given a non-existing branch name, `git hack` creates a new feature branch
with the main branch as its parent.

When given an existing contribution, observed, parked, or prototype branch,
`git hack` converts that branch to a feature branch.

When given no arguments, `git hack` converts the current contribution, observed,
parked, or prototype branch into a feature branch.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --prototype / -p

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches)).

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### upstream remote

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

### configuration

If [push-new-branches](../preferences/push-new-branches.md) is set, `git hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git hack` run fast. The first run of `git sync`
will create the remote tracking branch.

If the configuration setting
[create-prototype-branches](../preferences/create-prototype-branches.md) is set,
`git hack` always creates a
[prototype branch](../branch-types.md#prototype-branches).
