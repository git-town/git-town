# git town compress

<a type="gittown-command" />

```command-summary
git town compress [--dry-run] [-h | --help] [(-m | --message) <text>] [--no-verify] [-s | --stack] [-v | --verbose]
```

The _compress_ command squashes all commits on a branch into a single commit.

Assuming you have a feature branch with these commits:

```bash
$ git log --format='%s'
commit 1
commit 2
commit 3
```

Let's compress these three commits into a single commit:

```bash
git town compress
```

Now your branch has a single commit with the name of the first commit but
containing the changes of all three commits that existed on the branch before:

```bash
$ git log --format='%s'
commit 1
```

Git Town compresses feature branches and
[parked branches](https://www.git-town.com/preferences/parked-branches). It
doesn't compress
[perennial](https://www.git-town.com/preferences/perennial-branches),
[observed](https://www.git-town.com/preferences/observed-branches), and
[contribution](https://www.git-town.com/preferences/contribution-branches)
branches.

Branches must be in sync to compress them; run `git town sync` before running
this command.

## Options

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `-m <text>`<br>`--message <text>`

By default the now compressed commit uses the commit message of the first commit
in the branch. You can provide a custom commit message for the squashed commit
with the `--message <message>` aka `-m` flag, which works similar to the
[-m flag for `git commit`](https://git-scm.com/docs/git-commit#Documentation/git-commit.txt--mltmsggt).

Assuming you have a feature branch with these commits:

```bash
$ git log --format='%s'
commit 1
commit 2
commit 3
```

Let's compress these three commits into a single commit:

```bash
git town compress -m "compressed commit"
```

Now your branch has these commits:

```bash
$ git log --format='%s'
compressed commit
```

The new `compressed commit` now contains the changes from the old `commit 1`,
`commit 2`, and `commit 3`.

#### `--no-verify`

Disables the
[pre-commit hook](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

#### `-s`<br>`--stack`

To compress all branches in a [stack](../stacked-changes.md) provide the
`--stack` aka `-s` switch.

If you want to compress your commits every time you sync, choose the
[compress sync strategy](../preferences/sync-feature-strategy.md#compress) for
the respective branch type.

Assuming you have a [stacked change](../stacked-changes.md) consisting of two
feature branches. Each branch contains three commits.

```
main
 \
  branch-1
  |  * commit 1a
  |  * commit 1b
  |  * commit 1c
  branch-2
     * commit 2a
     * commit 2b
     * commit 2c
```

Let's compress the commits in all branches of this stack:

```
git town compress --stack
```

Now your stack contains these branches and commits:

```
main
 \
  branch-1
  |  * commit 1a
  branch-2
     * commit 2a
```

As usual, the new `commit 1a` contains the changes made in `branch 1`, i.e. the
changes from the old `commit 1a`, `commit 1b`, and `commit 1c`. The new
`commit 2a` contains the changes made in `branch 2`, i.e. the changes from the
old `commit 2a`, `commit 2b`, and `commit 2c`.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- the [compress sync strategy](../preferences/sync-feature-strategy.md#compress)
  automatically compresses branches when they get synced

<!-- keep-sorted end -->
