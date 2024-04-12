# git compress [branches]

The _compress_ command squashes all commits on a branch into a single commit.
Git Town compresses feature branches and
[parked branches](https://www.git-town.com/advanced-syncing#parked-branches) if
they are currently checked out. It doesn't compress
[perennial](https://www.git-town.com/preferences/perennial-branches),
[observed](https://www.git-town.com/advanced-syncing#observed-branches),
[contribution](https://www.git-town.com/advanced-syncing#contribution-branches)
branches.

Branches to compress must be in sync, so run `git sync` and resolve possible
merge conflicts before compressing a branch.

### Configuration

By default the compressed commit uses the commit message of the first commit in
the branch. You can provide a custom commit message for the squashed commit with
the `-m` branch, which works similar to the
[-m switch for `git commit`](https://git-scm.com/docs/git-commit#Documentation/git-commit.txt--mltmsggt).

To compress all branches in a [branch stack](../stacked-changes.md) provide the
`--stack` switch.

### Example: compressing commits on a branch

Assuming you have a feature branch with these commits:

```fish
$ git log --pretty=format:'%s'
commit 1
commit 2
commit 3
```

Let's compress these three commits into a single commit:

```fish
git compress
```

Now your branch has these commits:

```fish
$ git log --pretty=format:'%s'
commit 1
```

The new `commit 1` now contains the changes from the old `commit 1`, `commit 2`,
and `commit 3`.

### Example: compressing using a custom commit message

Assuming you have a feature branch with these commits:

```fish
$ git log --pretty=format:'%s'
commit 1
commit 2
commit 3
```

Let's compress these three commits into a single commit:

```fish
git compress -m "compressed commit"
```

Now your branch has these commits:

```fish
$ git log --pretty=format:'%s'
compressed commit
```

The new `compressed commit` now contains the changes from the old `commit 1`,
`commit 2`, and `commit 3`.

### Example: Compressing all branches in a stacked change

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
git compress --stack
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
