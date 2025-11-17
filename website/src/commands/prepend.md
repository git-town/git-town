# git town prepend

<a type="command-summary">

```command-summary
git town prepend [<branch-name>...] [-b | --beam] [--body <string>] [--propose] [-p | --prototype] [-t <text> | --title <text>] [-d | --detached] [-c | --commit] [-m | --message <message>] [--dry-run] [-v | --verbose] [-h | --help]
```

</a>

The _prepend_ command creates a new feature branch as the parent of the current
branch. It does that by inserting the new feature branch between the current
feature branch and it's existing parent.

Consider this stack:

```
main
 \
* feature-2
```

We are on the `feature-2` branch. After running `git town prepend feature-1`,
our repository has this stack:

```
main
 \
* feature-1
   \
    feature-2
```

If your Git workspace is clean (no uncommitted changes), it also
[syncs](sync.md) the current feature branch to ensure you work on top of the
current state of the repository. If the workspace is not clean (contains
uncommitted changes), `git town prepend` does not perform this sync to let you
commit your open changes.

If the branch you call this command from has a proposal, this command updates
it. To do so, it pushes the new branch.

## Options

#### `--body <string>`

Pre-populate the body of the pull request to create with the given text.
Requires `--propose`.

#### `--propose`

Propose the created branch.

To always propose new branches, set the
[share new branches](../preferences/share-new-branches.md) setting to `propose`.

#### `-p`<br>`--prototype`

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches).

#### `-t <text>`<br>`--title <text>`

Pre-populate the title of the pull request to create with the given text.
Requires `--propose`.

#### `-d`<br>`--detached`<br>`--no-detached`

The `--detached` aka `-d` flag enables
[detached mode](../preferences/detached.md) for the current command. If detached
mode is enabled through [configuration data](../preferences/detached.md), the
`--no-detached` flag disables detached mode for the current command.

In detached mode, feature branches don't receive updates from the perennial
branch at the root of your branch hierarchy. This can be useful in busy
monorepos.

#### `-c`<br>`--commit`

When given, commits the currently staged changes into the branch to create and
remains on the current branch. This is intended to quickly commit changes
unrelated to the current branch into another branch and keep hacking on the
current branch. Committing suppresses all branch updates to allow you to get
your open changes committed.

#### `-b`<br>`--beam`

Moves ("beams") one or more commits from the current branch to the new parent
branch that gets created. Lets you select the commits to beam via a visual
dialog. Beaming suppresses all branch updates. Any merge conflicts encountered
while beaming arise from moving the beamed commits.

#### `-m`<br>`--message`

Commit message to use together with `--commit`. Implies `--commit`.

#### `--push`<br>`--no-push`

The `--push`/`--no-push` argument overrides the
[push-branches](../preferences/push-branches.md) config setting.

#### `--stash`<br>`--no-stash`

Enables or disables [stashing](../preferences/stash.md) for this invocation.

#### `--sync`

Enables or disables [automatic syncing](../preferences/auto-sync.md) of the
current branch before prepending the new one.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

#### `--auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

## Configuration

If [share-new-branches](../preferences/share-new-branches.md) is configured,
`git town hack` creates a remote tracking branch for the new feature branch.
This behavior is disabled by default to make `git town hack` run fast. The first
run of `git town sync` will create the remote tracking branch.

If the configuration setting
[new-branch-type](../preferences/new-branch-type.md) is set, `git town prepend`
creates a branch with the given [type](../branch-types.md).

## See also

- [append](append.md) creates the new branch as a child of the current branch
- [hack](hack.md) creates the new branch as a child of the
  [main branch](../preferences/main-branch.md)
