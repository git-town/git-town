# git town walk

```command-summary
git town walk [-a | --all] [-s | --stack] [--dry-run] [-v | --verbose] [<command and arguments>]
```

The _walk_ command ("walking the branch hierarchy") opens a shell for each local branch,
allowing you to investigate the state of the repository on that branch.

With the `--all` flag, it walks through all local branches,
with the `--stack` flag it walks through all branches in the current branch stack.

If you provide a command, it executes that command for each branch without opening a shell.


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
