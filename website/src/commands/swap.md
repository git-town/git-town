# git town swap

```command-summary
git town swap [--dry-run] [-v | --verbose]
```

The _swap_ command moves the current branch one position forward in the stack.

Consider this branch stack:

```
main
 \
  branch-1
   \
*   branch-2
     \
      branch-3
```

After running `git town swap` you end up with this branch stack:

```
main
 \
  branch-2
   \
*   branch-1
     \
      branch-3
```

Moving branches up and down the stack helps you organize related branches
together, for example to review and ship them together, or [merge](merge.md)
them.

Please ensure all affected branches are in sync before running this command, and
optionally remove merge commits by [compressing](compress.md).

## Options

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
