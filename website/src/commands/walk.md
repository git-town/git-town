# git town walk

<a type="git-town-command" />

```command-summary
git town walk [<command and arguments>] [-a | --all] [--dry-run] [-h | --help] [--(non)-interactive] [-s | --stack] [-v | --verbose]
```

The _walk_ command ("walking the branch hierarchy") runs a given command
for all feature branches in your current stack or the entire repo.
Git Town checks out each branch, executes the command,
and then moves to the next branch.
If the command fails or leaves uncommitted changes behind,
`git town walk` stops so you can fix the problem or commit the generated changes
before continuing.

- use [git town continue](continue.md) to retry the command on the current
  branch
- use [git town skip](skip.md) to move on to the next branch
- use [git town undo](undo.md) to abort the iteration and undo all changes made
- use [git town status reset](status-reset.md) to abort the iteration and keep
  all changes made

If no shell command is provided,
drops you into an interactive shell for each branch.
You can manually run any shell commands,
then proceed to the next branch with [git town continue](continue.md)

## Examples

Consider this stack:

```sh
main
 \
  branch-1
   \
*   branch-2
     \
      branch-3
```

Running `git town walk --stack make lint` produces this output:

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

Iterate through all local branches.

#### `--dry-run`

Test-drive this command:
It prints the commands that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `--interactive`<br>`--non-interactive`

Enables or disables [interactive](../preferences/interactive.md) dialogs.

#### `-s`<br>`--stack`

Iterate through all branches of the stack that the current branch belongs to.

#### `-v`<br>`--verbose`

Print all Git commands executed under the hood to determine the repository
state.

## See also

<!-- keep-sorted start -->

- [branch](branch.md) displays the branch hierarchy and highlights the currently
  checked out branch in it
- [switch](switch.md) displays the branch hierarchy and lets you select a branch
  to switch to

<!-- keep-sorted end -->
