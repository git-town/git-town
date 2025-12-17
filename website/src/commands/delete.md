# git town delete

<a type="git-town-command" />

```command-summary
git town delete [<branch-name>...] [--dry-run] [-h | --help] [-v | --verbose]
```

The _delete_ command deletes the given branch from the local and if possible the
remote repository, removes commits of deleted branches from their descendents
(unless when using the
[merge sync strategy](../preferences/sync-feature-strategy.md#merge)), and
updates proposals of child branches to the parent of the deleted branch.

Consider this stack:

```
main
 \
  branch-1
   \
*   branch-2
     \
      branch-3
```

We are on the `branch-2` branch. After running `git town delete` we end up with
this stack, on the branch that was active before we switched to `branch-2`:

```
main
 \
  branch-1
   \
    branch-3
```

Git Town deletes only the parts of the branch that you own. If you delete
[feature](../branch-types.md#feature-branches),
[parked](../branch-types.md#parked-branches), or
[prototype](../branch-types.md#prototype-branches), it deletes the local and
tracking branch. When deleting
[contribution](../branch-types.md#contribution-branches),
[observed](../branch-types.md#observed-branches), or
[perennial](../branch-types.md#perennial-branches), it deletes only the local
branch because you don't own the tracking branch.

## Positional arguments

When called without arguments, the _delete_ command deletes the feature branch
you are on, including all uncommitted changes.

When called with a branch name, it deletes the given branch.

## Example

## Options

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
