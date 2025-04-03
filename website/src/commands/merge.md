# git town merge

```command-summary
git town merge [--dry-run] [-v | --verbose]
```

The _merge_ command merges the current branch with the branch ahead of it in the
current stack.

Consider this branch stack:

```
main
 \
  branch-1
   \
    branch-2
     \
*     branch-3
       \
        branch-4
```

We are on the `branch-3` branch. After running `git town merge`, the stack looks
like below, and the new `branch-3` branch contains the changes from the old
`branch-2` and `branch-3` branches.

```
main
 \
  branch-1
   \
*   branch-3
     \
      branch-4
```

Both branches must be [feature branches](../branch-types.md#feature-branches)
and in sync, i.e. run [git town sync](sync.md) before running`git town merge`.
All affected branches must be owned by you, i.e. no be
[contribution](../branch-types.md#contribution-branches),
[observed](../branch-types.md#observed-branches), or
[perennial](../branch-types.md#perennial-branches) branches.

When using the
[compress sync strategy](../preferences/sync-feature-strategy.md#compress), the
merged branch will contain two separate commits: one per merged branch. This
makes it easy to verify that both branches were merged as expected. To
consolidate these commits, run [git town sync](sync.md).

## Options

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
