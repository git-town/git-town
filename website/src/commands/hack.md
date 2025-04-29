# git town hack

```command-summary
git town hack [<branch-name>...] [-p | --prototype] [-d | --detached] [-c | --commit] [-m | --message <message>] [--propose] [--dry-run] [-v | --verbose]
```

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch.md) and brings all
uncommitted changes over to it.

Consider this stack:

```
main
 \
  branch-1
   \
*   branch-2
```

We are on the `branch-2` branch. After running `git hack branch-3`, our
workspace contains these branches:

```
main
 \
  branch-1
   \
    branch-2
 \
* branch-3
```

If your Git workspace is clean (no uncommitted changes), it also
[syncs](sync.md) the main branch to ensure you develop on top of the current
state of the repository. If the workspace is not clean (contains uncommitted
changes), `git town hack` does not perform this sync to let you commit your open
changes.

### Upstream remote

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

## Positional arguments

When given a non-existing branch name, `git town hack` creates a new feature
branch with the main branch as its parent.

When given an existing contribution, observed, parked, or prototype branch,
`git town hack` converts that branch to a feature branch.

When given no arguments, `git town hack` converts the current contribution,
observed, parked, or prototype branch into a feature branch.

## Options

#### `-p`<br>`--prototype`

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches).

#### `-d`<br>`--detached`

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch. This allows you to build out your stack and decide when to pull in
changes from other developers.

#### `-c`<br>`--commit`

When given, commits the currently staged changes into the branch to create and
remains on the current branch. This is intended to quickly commit changes
unrelated to the current branch into another branch and keep hacking on the
current branch. Committing suppresses all branch updates to allow you to get
your open changes committed.

#### `-b`<br>`--beam`

Moves ("beams") one or more commits from the current branch to the new feature
branch that gets created. Lets you select the commits to beam via a visual
dialog. Beaming suppresses all branch updates. Any merge conflicts encountered
while beaming arise from moving the beamed commits.

#### `-m`<br>`--message`

Commit message to use together with `--commit`. Implies `--commit`.

#### `--propose`

Propose the created branch.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

If [share-new-branches](../preferences/share-new-branches.md) is configured,
`git town hack` creates a remote tracking branch for the new feature branch.
This behavior is disabled by default to make `git town hack` run fast. The first
run of `git town sync` will create the remote tracking branch.

If the configuration setting
[new-branch-type](../preferences/new-branch-type.md) is set, `git town hack`
creates a branch with the given [type](../branch-types.md).
