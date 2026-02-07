# git town diff-parent

<a type="git-town-command" />

```command-summary
git town diff-parent [--diff-filter=<value>] [-h | --help] [--name-only] [-v | --verbose]
```

The _diff-parent_ command displays the changes made on a feature branch, i.e.
the diff between the current branch and its parent branch.

## Options

#### `--diff-filter=<value>`

#### `-h`<br>`--help`

Display help for this command.

#### `--name-only`

When set, displays only the names of files changed on the current branch. It
suppresses the diff output and does not show the actual content changes within
those files.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
