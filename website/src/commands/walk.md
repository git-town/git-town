# git town walk

```command-summary
git town walk [-a | --all] [-s | --stack] [--dry-run] [-v | --verbose] [<command and arguments>]
```

The _walk_ command ("walking the branch hierarchy") executes a given command for each branch.
If you don't provide a command, it exits to the shell on each branch.
In that case, run [git town continue](continue.md) to move to the next branch.

## Examples

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

If we run `git town walk --stack make lint` it prints this output:

```bash
[branch-1] make lint
# ... output of "make lint" for branch-1

[branch-2] make lint
# ... output of "make lint" for branch-2

[branch-3] make lint
# ... output of "make lint" for branch-3
```

## Options

#### `-a`<br>`--all`

Iterates through all local branches.

#### `-s`<br>`--stack`

The `--stack` aka `-s` iterates through all branches of the stack that the current branch belongs to.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
