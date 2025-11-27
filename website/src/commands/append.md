# git town append

<a type="gittown-command" />

```command-summary
git town append <branch-name> [--auto-resolve] [-b | --beam] [-c | --commit] [-d | --detached] [--dry-run] [-h | --help] [(-m | --message) <message>] [--propose] [-p | --prototype] [--push] [--stash] [--sync] [-v | --verbose]
```

The _append_ command creates a new feature branch with the given name as a
direct child of the current branch and brings over all uncommitted changes to
the new branch.

Consider this stack:

```
main
 \
* feature-1
```

We are on the `feature-1` branch. After running `git town append feature-2`, the
repository will have this stack:

```
main
 \
  feature-1
   \
*   feature-2
```

If your Git workspace is clean (no uncommitted changes), it also
[syncs](sync.md) the current branch to ensure your work in the new branch
happens on top of the current state of the repository. If the workspace contains
uncommitted changes, `git town append` does not perform this sync to let you
commit your open changes first and then sync manually.

## Positional argument

When given a non-existing branch name, `git town append` creates a new feature
branch with the main branch as its parent.

## Options

#### `--auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

#### `-b`<br>`--beam`

Moves ("beams") one or more commits from the current branch to the new child
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

#### `-m <msg>`<br>`--message <msg>`

Commit message to use together with `--commit`. Implies `--commit`.

#### `--propose`

Propose the created branch.

To always propose new branches, set the
[share new branches](../preferences/share-new-branches.md) setting to `propose`.

#### `-p`<br>`--prototype`

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches).

#### `--push`<br>`--no-push`

The `--push`/`--no-push` argument overrides the
[push-branches](../preferences/push-branches.md) config setting.

#### `--stash`<br>`--no-stash`

Enables or disables [stashing](../preferences/stash.md) for this invocation.

#### `--sync`

Enables or disables [automatic syncing](../preferences/auto-sync.md) of the
current branch before appending the new one.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

If [share-new-branches](../preferences/share-new-branches.md) is configured,
`git town append` also creates the tracking branch for the new feature branch.
This behavior is disabled by default to make `git town append` run fast and save
CI runs. The first run of `git town sync` will create the remote tracking
branch.

If the configuration setting
[new-branch-type](../preferences/new-branch-type.md) is set, `git town append`
creates a branch with the given [type](../branch-types.md).

## See also

<!-- keep-sorted start -->

- [hack](hack.md) creates the new branch as a child of the
  [main branch](../preferences/main-branch.md)
- [prepend](prepend.md) creates the new branch as a parent of the current branch

<!-- keep-sorted end -->
