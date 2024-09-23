# git town hack

> _git town hack [--prototype] [branch-name...]_

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch.md) and brings all
uncommitted changes over to it.

If your Git workspace is clean (no uncommitted changes), it also
[syncs](sync.md) the main branch to ensure you develop on top of the current
state of the repository. If the workspace is not clean (contains uncommitted
changes), `git town hack` does not perform this sync to let you commit your open
changes.

### Positional arguments

When given a non-existing branch name, `git town hack` creates a new feature branch
with the main branch as its parent.

When given an existing contribution, observed, parked, or prototype branch,
`git town hack` converts that branch to a feature branch.

When given no arguments, `git town hack` converts the current contribution, observed,
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

If [push-new-branches](../preferences/push-new-branches.md) is set, `git town hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git town hack` run fast. The first run of `git town sync`
will create the remote tracking branch.

If the configuration setting
[create-prototype-branches](../preferences/create-prototype-branches.md) is set,
`git town hack` always creates a
[prototype branch](../branch-types.md#prototype-branches).
