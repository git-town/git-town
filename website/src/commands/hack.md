# git town hack

<a type="git-town-command" />

```command-summary
git town hack [<branch-name>...] [--(no)-auto-resolve] [-b | --beam] [-c | --commit] [-d | --(no)-detached] [--dry-run] [-h | --help] [(-m | --message) <message>] [--propose] [-p | --prototype] [--(no)-stash] [--(no)-sync] [-v | --verbose]
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

## Options

#### `--auto-resolve`<br>`--no-auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

#### `-b`<br>`--beam`

Moves ("beams") one or more commits from the current branch to the new feature
branch that gets created. Lets you select the commits to beam via a visual
dialog. Beaming suppresses all branch updates. Any merge conflicts encountered
while beaming arise from moving the beamed commits.

#### `-c`<br>`--commit`

When given, commits the currently staged changes into the branch to create and
remains on the current branch. This is intended to quickly commit changes
unrelated to the current branch into another branch and keep hacking on the
current branch. Committing suppresses all branch updates to allow you to get
your open changes committed.

#### `-d`<br>`--detached`<br>`--no-detached`

The `--detached` aka `-d` flag enables
[detached mode](../preferences/detached.md) for the current command. If detached
mode is enabled through [configuration data](../preferences/detached.md), the
`--no-detached` flag disables detached mode for the current command.

In detached mode, feature branches don't receive updates from the perennial
branch at the root of your branch hierarchy. This can be useful in busy
monorepos.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `-m <text>`<br>`--message <text>`

Commit message to use together with `--commit`. Implies `--commit`.

#### `--propose`

Propose the created branch.

To always propose new branches, set the
[share new branches](../preferences/share-new-branches.md) setting to `propose`.

#### `-p`<br>`--prototype`

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches).

#### `--stash`<br>`--no-stash`

Enables or disables [stashing](../preferences/stash.md) for this invocation.

#### `--sync`<br>`--no-sync`

Enables or disables [automatic syncing](../preferences/auto-sync.md) before
creating the new branch.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

If the configuration setting
[new-branch-type](../preferences/new-branch-type.md) is set, `git town hack`
creates a branch with the given [type](../branch-types.md).

If [share-new-branches](../preferences/share-new-branches.md) is configured,
`git town hack` creates a remote tracking branch and optionally a
[proposal](propose.md) for the new feature branch. This behavior is disabled by
default to make `git town hack` run fast. The first run of `git town sync` will
create the remote tracking branch.

## See also

<!-- keep-sorted start -->

- [append](append.md) creates the new branch as a child of the current branch
- [prepend](prepend.md) creates the new branch as a parent of the current branch

<!-- keep-sorted end -->
