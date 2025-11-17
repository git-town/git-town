# git town config

<a type="command-summary">

```command-summary
git town config [-d | --display-types <type>] [-h | --help] [-v | --verbose]
```

</a>

The _config_ command displays and updates the local Git Town configuration.

## Subcommands

Running without a subcommand shows the current Git Town configuration.

- The [get-parent](config-get-parent.md) subcommand outputs the parent branch of
  the current or given branch.
- The [remove](config-remove.md) subcommand removes all Git Town related
  configuration from the current Git repository.
- The [init](init.md) subcommand launches Git Town's setup assistant.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
