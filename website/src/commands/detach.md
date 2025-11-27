# git town detach

<a type="gittown-command" />

```command-summary
git town detach [--dry-run] [-h | --help] [-v | --verbose]
```

The _detach_ command removes the current branch from the stack it is in and
makes it a stand-alone top-level branch that ships directly into your main
branch.

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

We are on the `branch-2` branch. After running `git town detach`, we end up with
with stack:

```
main
 \
  branch-1
   \
    branch-3
 \
* branch-2
```

This is useful when a branch in a stack makes changes that are independent from
the changes made by other branches in this stack. Detaching such independent
branches removes "noise" from your stack, i.e. reduces it to changes that belong
together, and allows you to get more of your changes reviewed and shipped
concurrently.

Please ensure all affected branches are in [sync](sync.md) before running this
command, and remove merge commits by [compressing](compress.md).

## Options

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [set-parent](set-parent.md) move the current branch and its descendents under
  a different parent
- [swap](swap.md) move the current branch up in the stack

<!-- keep-sorted end -->
