# git town config

<a type="git-town-command" />

```command-summary
git town config [(-d | --display-types) <type>] [-h | --help] [-v | --verbose]
```

The _config_ command displays and updates the local Git Town configuration.

## Subcommands

Running without a subcommand shows the current Git Town configuration.

- The [get-parent](config-get-parent.md) subcommand outputs the parent branch of
  the current or given branch.
- The [remove](config-remove.md) subcommand removes all Git Town related
  configuration from the current Git repository.
- The [init](init.md) subcommand launches Git Town's setup assistant.

## Options

#### `-d <branch-types>`<br>`--display-types <branch-types>`

This flag allows customizing whether Git Town also displays the branch type in
addition to the branch name when showing a list of branches. More info
[here](../preferences/display-types.md#cli-flags).

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
