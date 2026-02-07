# git town diff-parent

<a type="git-town-command" />

```command-summary
git town diff-parent [--diff-filter <value>] [-h | --help] [--name-only] [-v | --verbose]
```

The _diff-parent_ command displays the changes made on a feature branch, i.e.
the diff between the current branch and its parent branch.

## Options

#### `--diff-filter <value>`

When set, forwards the given value to
[git diff --diff-filter](https://git-scm.com/docs/git-diff#Documentation/git-diff.txt---diff-filterACDMRTUXB).

This allows you to restrict the diff to specific change types (for example,
added, modified, or deleted files) using the same semantics as native Git.

#### `-h`<br>`--help`

Display help for this command.

#### `--name-only`

When set, displays only the names of files changed on the current branch. It
suppresses the diff output and does not show the actual content changes within
those files.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
