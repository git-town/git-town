# git town walk

```command-summary
git town walk [-a | --all] [-s | --stack] [--dry-run] [-v | --verbose]
```

The _walk_ command ("walking the branch hierarchy") executes the given command for multiple branches.
With the `--all` flag, it executes the command for all local branches,
with the `--stack` flag for all branches in the current branch stack.

## Options

#### `-a`<br>`--all`

Iterates through all local branches.

#### `-s`<br>`--stack`

The `--stack` aka `-s` parameter makes Git Town sync all branches in the stack
that the current branch belongs to.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
